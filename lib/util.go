package lib

import (
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"

	log "github.com/sirupsen/logrus"
)

// GetPackagePath returns the package path relative to $GOPATH/src
func GetPackagePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	const src = "src/"
	srcPos := strings.Index(absPath, src)
	return absPath[srcPos+len(src):]
}

// Unpanic recovers from a panic.
func Unpanic() func() {
	return func() {
		if err := Err(recover()); err != nil {
			log.Error(err.ErrorStack())
			log.Error("Mockhiato encountered an error! If this seems like a bug, please report it to https://github.com/littledot/mockhiato/issues/new")
		}
	}
}

// Err wraps value in a stack trace
func Err(value interface{}) *errors.Error {
	switch rec := value.(type) {
	case nil:
		return nil
	default:
		return errors.Wrap(rec, 1)
	}
}
