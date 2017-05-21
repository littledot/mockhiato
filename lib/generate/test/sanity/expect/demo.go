package example

import (
	"bytes"
	"encoding/json"
	"image/png"
	"os"
)

// Target contains 2 interfaces that needs to be mocked.
type Target interface {
	Hello(a string) (b int, err error)
	World(c int, a ...string) (b map[int]interface{})
	Yes(fi os.FileInfo) (pnge *png.Encoder, err error)

	Bool(bool) bool
	Error(error) error
	Interface(interface{}) interface{}
	Ptr(uintptr) uintptr
	Func(func(int) error) func(int) error
	Slice([]int) []error
	Chan(chan int) chan error
	Map(map[int]error) map[int]error
	Float(float32, float64) (float32, float64)
	Complex(complex64, complex128) (complex64, complex128)
	Byte(byte, rune, string) (byte, rune, string)
	Int(int, int8, int16, int32, int64) (int, int8, int16, int32, int64)
	Uint(uint, uint8, uint16, uint32, uint64) (uint, uint8, uint16, uint32, uint64)
}

// B B
type B interface {
	No(*json.Decoder) *bytes.Buffer
}
