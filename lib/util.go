package lib

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
)

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
