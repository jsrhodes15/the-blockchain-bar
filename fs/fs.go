package fs

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// Expands a file path
// 1. replace tilde with users home dir
// 2. expands embedded environment variables
// 3. cleans the path, e.g. /a/b/../c -> /a/c
// Note, it has limitations, e.g. ~someuser/tmp will not be expanded
func ExpandPath(p string) string {
	if i := strings.Index(p, ":"); i > 0 {
		return p
	}
	if i := strings.Index(p, "@"); i > 0 {
		return p
	}
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		if home := GetHomeDir(); home != "" {
			p = home + p[1:]
		}
	}
	return path.Clean(os.ExpandEnv(p))
}

func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

func GetHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return homeDir
}
