package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	var _, err = os.Stat(path)
	return err == nil
}

func Filter[T any](source []T, filterFunc func(T) bool) (ret []T) {
	var returnArray = []T{}
	for _, s := range source {
		if filterFunc(s) {
			returnArray = append(returnArray, s)
		}
	}
	return returnArray
}

func FindSegment(p, segment string) (string, error) {
	p = filepath.ToSlash(p)

	// Ensure segment is slash-wrapped for exact matching
	search := "/" + strings.Trim(segment, "/") + "/"
	idx := strings.Index(p, search)
	if idx == -1 {
		return "", fmt.Errorf("segment %q not found in path: %s", segment, p)
	}

	return filepath.FromSlash(p[:idx+len(search)-1]), nil
}
