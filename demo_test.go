package main_test

import (
	"errors"
	"testing"

	"gitlab.com/littledot/mockhiato/mocks"
)

func TestMock(t *testing.T) {
	m := &mocks.AAA{}
	m.On("Hello", "world").Return(1, errors.New("boo!"))
	a0, a1 := m.Hello("world")
	if a0 != 1 {
		t.Errorf("a0: %d", a0)
	}
	if a1 == nil {
		t.Errorf("a1: %s", a1)
	}
}
