package tls_test

import (
	"path/filepath"
	"runtime"
)

// path returns the absolute path the given relative file or directory path,
// relative to the cert directory in the user's GOPATH.
// If rel is already absolute, it is returned unmodified.
func path(rel string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(currentFile)

	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}

func CA() string {
	return path("ca.pem")
}

func Cert() string {
	return path("server.pem")
}

func Key() string {
	return path("server.key")
}
