package lib

import (
	"path/filepath"
	"strings"
)

// GetPackagePath returns the package path relative to $GOPATH/src
func GetPackagePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	const src = "src/"
	srcPos := strings.Index(absPath, src)
	return absPath[srcPos+len(src) : len(absPath)]
}
