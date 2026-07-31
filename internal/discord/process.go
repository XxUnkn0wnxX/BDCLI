package discord

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/betterdiscord/cli/internal/output"
	"github.com/shirou/gopsutil/v3/process"
)

// stop terminates Discord if it is running. The new injection method modifies
// app.asar, which the running Discord process holds a lock on, so Discord must
// be stopped before inject/uninject can touch it. It returns the executable path
// of the killed process (captured before the kill, for a later start) and whether
// Discord was running. Flatpak/Snap relaunch via their own run commands and don't
// use the exe.
func (discord *DiscordInstall) stop() (exe string, wasRunning bool, err error) {
	if running, _ := discord.isRunning(); !running {
		output.Printf("✅ %s is not running.\n", discord.Channel.Name())
		return "", false, nil
	}

	// Capture the executable before killing — afterward the process is gone.
	exe = discord.getFullExe()

	if err := discord.kill(); err != nil {
		output.Printf("❌ Unable to stop %s. Please close it and try again.\n", discord.Channel.Name())
		output.Printf("   %s\n", err.Error())
		return exe, true, err
	}

	output.Printf("✅ Stopped %s\n", discord.Channel.Name())
	return exe, true, nil
}

// start launches Discord. exe is the executable path captured by stop() and is
// used for native installs; Flatpak/Snap launch via their run commands.
func (discord *DiscordInstall) start(exe string) error {
	// Determine command based on installation type
	var cmd *exec.Cmd
	if discord.IsFlatpak {
		cmd = exec.Command("flatpak", "run", "com.discordapp."+discord.Channel.Exe())
	} else if discord.IsSnap {
		cmd = exec.Command("snap", "run", discord.Channel.Exe())
	} else {
		// Use binary found while killing the process for non-Flatpak/Snap installs
		if exe == "" {
			output.Printf("❌ Unable to restart %s, please do so manually.\n", discord.Channel.Name())
			return fmt.Errorf("could not determine executable path for %s", discord.Channel.Name())
		}
		cmd = exec.Command(exe)
	}

	// Set working directory to user home
	cmd.Dir, _ = os.UserHomeDir()

	if err := cmd.Start(); err != nil {
		output.Printf("❌ Unable to restart %s, please do so manually.\n", discord.Channel.Name())
		output.Printf("   %s\n", err.Error())
		return err
	}
	output.Printf("✅ Restarted %s\n", discord.Channel.Name())
	return nil
}

func (discord *DiscordInstall) isRunning() (bool, error) {
	name := discord.Channel.Exe()
	processes, err := process.Processes()

	// If we can't even list processes, bail out
	if err != nil {
		return false, fmt.Errorf("could not list processes")
	}

	// Search for desired process(es)
	for _, p := range processes {
		n, err := p.Name()

		// Ignore processes requiring Admin/Sudo
		if err != nil {
			continue
		}

		// We found our target return
		if n == name {
			return true, nil
		}
	}

	// If we got here, process was not found
	return false, nil
}

func (discord *DiscordInstall) kill() error {
	name := discord.Channel.Exe()
	processes, err := process.Processes()

	// If we can't even list processes, bail out
	if err != nil {
		return fmt.Errorf("could not list processes")
	}

	// Search for desired process(es)
	for _, p := range processes {
		n, err := p.Name()

		// Ignore processes requiring Admin/Sudo
		if err != nil {
			continue
		}

		// We found our target, kill it
		if n == name {
			var killErr = p.Kill()

			// We found it but can't kill it, bail out
			if killErr != nil {
				return killErr
			}
		}
	}

	// If we got here, everything was killed without error
	return nil
}

func (discord *DiscordInstall) getFullExe() string {
	name := discord.Channel.Exe()

	var exe = ""
	processes, err := process.Processes()
	if err != nil {
		return exe
	}
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			continue
		}
		if n == name {
			if len(exe) == 0 {
				exe, _ = p.Exe()
			}
		}
	}
	return exe
}
