package cmd

import (
	"os"
	"strings"
)

// isDirTarget reports whether path is a directory output target. It returns
// true for explicit directory paths (trailing separator or slash) and for
// existing directories.
func isDirTarget(path string) bool {
	if path == "" {
		return false
	}
	if strings.HasSuffix(path, string(os.PathSeparator)) || strings.HasSuffix(path, "/") {
		return true
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
