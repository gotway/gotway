package cert

import (
	"path/filepath"
	"runtime"
)

// Path returns the absolute path the given relative file or directory path,
// relative to the github.com/gosmo-devs/microgateway/cert directory in the user's GOPATH.
// If rel is already absolute, it is returned unmodified.
func Path(rel string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(currentFile)

	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}
