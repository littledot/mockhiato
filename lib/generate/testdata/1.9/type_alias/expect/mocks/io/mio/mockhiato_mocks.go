package mio

// Code generated by Mockhiato. DO NOT EDIT.
import (
	mock "github.com/stretchr/testify/mock"
)

// NewReaderMock creates a new ReaderMock
func NewReaderMock() *ReaderMock { return &ReaderMock{} }

// ReaderMock implements mio.Reader
type ReaderMock struct{ mock.Mock }

// Read implements (mio.Reader).Read
func (r *ReaderMock) Read(p0 []byte) (int, error) {
	ret := r.Called(p0)
	var ret0 int
	if a := ret.Get(0); a != nil {
		ret0 = a.(int)
	}
	var ret1 error
	if a := ret.Get(1); a != nil {
		ret1 = a.(error)
	}
	return ret0, ret1
}
