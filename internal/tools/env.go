package tools

import (
	"errors"
	"os"
	"path/filepath"
)

func FindRepoRoot(start string) (string, error) {
	dir := start
	for {
		if dir == "" || dir == string(filepath.Separator) {
			return "", errors.New("could not locate repository root")
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not locate repository root")
		}
		dir = parent
	}
}
