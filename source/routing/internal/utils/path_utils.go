package utils

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	projectRoot     string
	projectRootOnce sync.Once
)

// getProjectRoot tries to find the project root by walking up
// until it finds a go.mod file.
func getProjectRoot() string {
	projectRootOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			panic("cannot determine working directory: " + err.Error())
		}

		// Walk up until we find go.mod
		for {
			if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
				projectRoot = wd
				return
			}
			parent := filepath.Dir(wd)
			if parent == wd { // reached filesystem root
				panic("go.mod not found, cannot determine project root")
			}
			wd = parent
		}
	})
	return projectRoot
}

// ResolvePath resolves a relative path to the project root.
// Example: ResolvePath("dev/requests/request.json")
func ResolvePath(relPath string) string {
	return filepath.Join(getProjectRoot(), relPath)
}