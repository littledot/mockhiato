package example

// Target is a sophisticated interface.
type Target interface {
	Bool(bool) bool
	Chan(chan int) chan int
	Complex(complex64, complex128) (complex64, complex128)
	Error(error) error
	Float(float32, float64) (float32, float64)
	Func(func(int) float32) func(int) float32
	Int(int, int8, int16, int32, int64) (int, int8, int16, int32, int64)
	Interface(interface{}) interface{}
	Map(map[int]error) map[int]error
	Ptr(uintptr) uintptr
	Slice([]int) []int
	Text(byte, rune, string) (byte, rune, string)
	Uint(uint, uint8, uint16, uint32, uint64) (uint, uint8, uint16, uint32, uint64)
}
