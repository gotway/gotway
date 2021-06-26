package tls_test

import (
	"path/filepath"
	"runtime"
)

// Path returns the absolute path the given relative file or directory path,
// relative to the cert directory in the user's GOPATH.
// If rel is already absolute, it is returned unmodified.
func Path(rel string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(currentFile)

	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}

func Cert() string {
	return Path("server.pem")
}

func Key() string {
	return Path("server.key")
}
