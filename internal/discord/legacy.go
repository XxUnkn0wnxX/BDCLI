package discord

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/betterdiscord/cli/internal/models"
)

//go:embed assets/legacy_injection.js
var legacyInjectionScript []byte

const canonicalLegacyLoader = "module.exports = require(\"./core.asar\");"

var legacyVersionRegex = regexp.MustCompile(`^(?:app-)?[0-9]+\.[0-9]+\.[0-9]+$`)

// legacyInjectionState is the complete information needed to restore one
// legacy loader after a later operation fails.
type legacyInjectionState struct {
	indexPath    string
	original     []byte
	originalMode os.FileMode
}

// findLegacyInjectionStates locates historical core loaders below the supplied
// config root. It deliberately does not use platform globals, so discovery can
// be tested against a temporary directory on every platform.
func findLegacyInjectionStates(configRoot string, channel models.DiscordChannel) ([]legacyInjectionState, error) {
	if channel.Name() == "" {
		return nil, nil
	}

	channelRoot := filepath.Join(configRoot, legacyChannelName(channel))
	entries, err := os.ReadDir(channelRoot)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("inspect legacy %s directory: %w", channel.Name(), err)
	}

	var candidates []string
	for _, entry := range entries {
		if !legacyVersionRegex.MatchString(entry.Name()) || !isRealDirectory(filepath.Join(channelRoot, entry.Name())) {
			continue
		}

		versionRoot := filepath.Join(channelRoot, entry.Name(), "modules")
		if !isRealDirectory(versionRoot) {
			continue
		}
		moduleEntries, err := os.ReadDir(versionRoot)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("inspect legacy modules in %s: %w", entry.Name(), err)
		}

		for _, module := range moduleEntries {
			if !module.IsDir() || !isRealDirectory(filepath.Join(versionRoot, module.Name())) {
				continue
			}

			moduleRoot := filepath.Join(versionRoot, module.Name())
			switch {
			case module.Name() == "discord_desktop_core":
				candidates = append(candidates, moduleRoot)
			case strings.HasPrefix(module.Name(), "discord_desktop_core-") && len(module.Name()) > len("discord_desktop_core-"):
				wrappedRoot := filepath.Join(moduleRoot, "discord_desktop_core")
				if isRealDirectory(wrappedRoot) {
					candidates = append(candidates, wrappedRoot)
				}
			}
		}
	}

	sort.Strings(candidates)
	states := make([]legacyInjectionState, 0, len(candidates))
	for _, candidate := range candidates {
		state, ok, err := inspectLegacyCandidate(candidate)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		states = append(states, state)
	}

	return states, nil
}

func legacyChannelName(channel models.DiscordChannel) string {
	return strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "")
}

func inspectLegacyCandidate(moduleRoot string) (legacyInjectionState, bool, error) {
	corePath := filepath.Join(moduleRoot, "core.asar")
	indexPath := filepath.Join(moduleRoot, "index.js")

	if !isRegularFile(corePath) || !isRegularNonSymlink(indexPath) {
		return legacyInjectionState{}, false, nil
	}

	indexInfo, err := os.Stat(indexPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return legacyInjectionState{}, false, nil
		}
		return legacyInjectionState{}, false, fmt.Errorf("inspect legacy loader %s: %w", indexPath, err)
	}
	original, err := os.ReadFile(indexPath)
	if err != nil {
		return legacyInjectionState{}, false, fmt.Errorf("read legacy loader %s: %w", indexPath, err)
	}

	if bytes.Equal(original, legacyInjectionScript) {
		return legacyInjectionState{
			indexPath:    indexPath,
			original:     original,
			originalMode: indexInfo.Mode(),
		}, true, nil
	}

	if containsBetterDiscordMarker(original) {
		return legacyInjectionState{}, false, fmt.Errorf(
			"unknown BetterDiscord legacy loader at %s; refusing to modify it; inspect or remove it manually before retrying",
			indexPath,
		)
	}

	return legacyInjectionState{}, false, nil
}

func isRegularFile(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0
}

func isRealDirectory(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.IsDir() && info.Mode()&os.ModeSymlink == 0
}

func isRegularNonSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0
}

func containsBetterDiscordMarker(data []byte) bool {
	return bytes.Contains(bytes.ToLower(data), []byte("betterdiscord"))
}

// cleanupLegacyInjectionStates replaces each exact historical loader with the
// canonical passthrough. A failure restores every target touched by the
// operation, including a target whose write was only partially completed.
func cleanupLegacyInjectionStates(states []legacyInjectionState) error {
	if len(states) == 0 {
		return nil
	}

	attempted := make([]legacyInjectionState, 0, len(states))
	fail := func(err error) error {
		if len(attempted) == 0 {
			return err
		}
		if rollbackErr := restoreLegacyInjectionStates(attempted); rollbackErr != nil {
			return errors.Join(err, fmt.Errorf("rollback legacy loader cleanup: %w", rollbackErr))
		}
		return err
	}
	for _, state := range states {
		if !bytes.Equal(state.original, legacyInjectionScript) {
			return fail(fmt.Errorf("refusing to clean non-legacy loader %s", state.indexPath))
		}
		current, err := os.ReadFile(state.indexPath)
		if err != nil {
			return fail(fmt.Errorf("recheck legacy loader %s: %w", state.indexPath, err))
		}
		if !bytes.Equal(current, state.original) {
			return fail(fmt.Errorf("legacy loader %s changed after preflight; refusing to modify it", state.indexPath))
		}
		attempted = append(attempted, state)
		if err := writeLegacyIndex(state.indexPath, []byte(canonicalLegacyLoader), state.originalMode); err != nil {
			return fail(err)
		}
	}

	return nil
}

// runLegacyInstallAction migrates exact legacy loaders before action. If the
// install action fails, the migration is reverted so the pre-action state is
// preserved; a restoration failure is joined to the primary action error.
func runLegacyInstallAction(states []legacyInjectionState, action func() error) error {
	if err := cleanupLegacyInjectionStates(states); err != nil {
		return err
	}
	if err := action(); err != nil {
		if restoreErr := restoreLegacyInjectionStates(states); restoreErr != nil {
			return errors.Join(err, restoreErr)
		}
		return err
	}
	return nil
}

// runLegacyRemovalAction migrates exact legacy loaders before action. Removal
// deliberately does not restore them if the later action fails: the canonical
// loader remains in place and is safe for a subsequent retry.
func runLegacyRemovalAction(states []legacyInjectionState, action func() error) error {
	if err := cleanupLegacyInjectionStates(states); err != nil {
		return err
	}
	return action()
}

// restoreLegacyInjectionStates writes the exact bytes and permissions saved by
// findLegacyInjectionStates. Restoration attempts every state and joins errors
// so one failed target does not hide failures for the remaining targets.
func restoreLegacyInjectionStates(states []legacyInjectionState) error {
	var restoreErrs []error
	for _, state := range states {
		if err := writeLegacyIndex(state.indexPath, state.original, state.originalMode); err != nil {
			restoreErrs = append(restoreErrs, fmt.Errorf("restore %s: %w", state.indexPath, err))
		}
	}
	return errors.Join(restoreErrs...)
}

func writeLegacyIndex(path string, contents []byte, mode os.FileMode) error {
	if !isRegularNonSymlink(path) {
		return fmt.Errorf("legacy loader %s is no longer a regular non-symlink file", path)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, legacyPermissionBits(mode))
	if err != nil {
		return fmt.Errorf("write legacy loader %s: %w", path, err)
	}

	n, writeErr := file.Write(contents)
	if writeErr == nil && n != len(contents) {
		writeErr = io.ErrShortWrite
	}
	closeErr := file.Close()
	if writeErr != nil || closeErr != nil {
		return errors.Join(writeErr, closeErr)
	}
	if err := os.Chmod(path, legacyPermissionBits(mode)); err != nil {
		return fmt.Errorf("preserve permissions on legacy loader %s: %w", path, err)
	}
	return nil
}

func legacyPermissionBits(mode os.FileMode) os.FileMode {
	return mode.Perm() | mode&(os.ModeSetuid|os.ModeSetgid|os.ModeSticky)
}
