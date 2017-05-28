package example

import (
	"bytes"
	"encoding/json"
)

// Target is an interface that should be mocked.
type Target interface {
	// All Go types should be supported
	Bool(bool) bool
	Interface(interface{}) interface{}
	Ptr(uintptr) uintptr
	Func(func(int) bool) func(int) bool
	Slice([]int) []bool
	Chan(chan int) chan bool
	Map(map[int]bool) map[int]bool
	Float(float32, float64) (float32, float64)
	Complex(complex64, complex128) (complex64, complex128)
	Byte(byte, rune, string) (byte, rune, string)
	Int(int, int8, int16, int32, int64) (int, int8, int16, int32, int64)
	Uint(uint, uint8, uint16, uint32, uint64) (uint, uint8, uint16, uint32, uint64)

	// Param names should be ignored
	Hello(a string) (b int, err bool)
	World(c int, a ...string) (b map[int]interface{})
}

// B is an interface that should be mocked
type B interface {
	No(*json.Decoder) *bytes.Buffer
}
