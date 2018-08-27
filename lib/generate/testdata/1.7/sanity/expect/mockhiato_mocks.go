package example

// Code generated by Mockhiato. DO NOT EDIT.
import (
	bytes "bytes"
	json "encoding/json"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

// NewDependenciesMock creates a new DependenciesMock
func NewDependenciesMock() *DependenciesMock { return &DependenciesMock{} }

// DependenciesMock implements example.Dependencies
type DependenciesMock struct{ mock.Mock }

// GoDependency implements (example.Dependencies).GoDependency
func (r *DependenciesMock) GoDependency(p0 *json.Decoder) *bytes.Buffer {
	ret := r.Called(p0)
	var ret0 *bytes.Buffer
	if a := ret.Get(0); a != nil {
		ret0 = a.(*bytes.Buffer)
	}
	return ret0
}

// VendorDependency implements (example.Dependencies).VendorDependency
func (r *DependenciesMock) VendorDependency(p0 assert.TestingT) {
	r.Called(p0)
	return
}

// NewPrimitivesMock creates a new PrimitivesMock
func NewPrimitivesMock() *PrimitivesMock { return &PrimitivesMock{} }

// PrimitivesMock implements example.Primitives
type PrimitivesMock struct{ mock.Mock }

// Bool implements (example.Primitives).Bool
func (r *PrimitivesMock) Bool(p0 bool) bool {
	ret := r.Called(p0)
	var ret0 bool
	if a := ret.Get(0); a != nil {
		ret0 = a.(bool)
	}
	return ret0
}

// Byte implements (example.Primitives).Byte
func (r *PrimitivesMock) Byte(p0 byte, p1 rune, p2 string) (byte, rune, string) {
	ret := r.Called(p0, p1, p2)
	var ret0 byte
	if a := ret.Get(0); a != nil {
		ret0 = a.(byte)
	}
	var ret1 rune
	if a := ret.Get(1); a != nil {
		ret1 = a.(rune)
	}
	var ret2 string
	if a := ret.Get(2); a != nil {
		ret2 = a.(string)
	}
	return ret0, ret1, ret2
}

// Chan implements (example.Primitives).Chan
func (r *PrimitivesMock) Chan(p0 chan int) chan bool {
	ret := r.Called(p0)
	var ret0 chan bool
	if a := ret.Get(0); a != nil {
		ret0 = a.(chan bool)
	}
	return ret0
}

// Complex implements (example.Primitives).Complex
func (r *PrimitivesMock) Complex(p0 complex64, p1 complex128) (complex64, complex128) {
	ret := r.Called(p0, p1)
	var ret0 complex64
	if a := ret.Get(0); a != nil {
		ret0 = a.(complex64)
	}
	var ret1 complex128
	if a := ret.Get(1); a != nil {
		ret1 = a.(complex128)
	}
	return ret0, ret1
}

// Error implements (example.Primitives).Error
func (r *PrimitivesMock) Error(p0 error) error {
	ret := r.Called(p0)
	var ret0 error
	if a := ret.Get(0); a != nil {
		ret0 = a.(error)
	}
	return ret0
}

// Float implements (example.Primitives).Float
func (r *PrimitivesMock) Float(p0 float32, p1 float64) (float32, float64) {
	ret := r.Called(p0, p1)
	var ret0 float32
	if a := ret.Get(0); a != nil {
		ret0 = a.(float32)
	}
	var ret1 float64
	if a := ret.Get(1); a != nil {
		ret1 = a.(float64)
	}
	return ret0, ret1
}

// Func implements (example.Primitives).Func
func (r *PrimitivesMock) Func(p0 func(int) bool) func(int) bool {
	ret := r.Called(p0)
	var ret0 func(int) bool
	if a := ret.Get(0); a != nil {
		ret0 = a.(func(int) bool)
	}
	return ret0
}

// Hello implements (example.Primitives).Hello
func (r *PrimitivesMock) Hello(p0 string) (int, bool) {
	ret := r.Called(p0)
	var ret0 int
	if a := ret.Get(0); a != nil {
		ret0 = a.(int)
	}
	var ret1 bool
	if a := ret.Get(1); a != nil {
		ret1 = a.(bool)
	}
	return ret0, ret1
}

// Int implements (example.Primitives).Int
func (r *PrimitivesMock) Int(p0 int, p1 int8, p2 int16, p3 int32, p4 int64) (int, int8, int16, int32, int64) {
	ret := r.Called(p0, p1, p2, p3, p4)
	var ret0 int
	if a := ret.Get(0); a != nil {
		ret0 = a.(int)
	}
	var ret1 int8
	if a := ret.Get(1); a != nil {
		ret1 = a.(int8)
	}
	var ret2 int16
	if a := ret.Get(2); a != nil {
		ret2 = a.(int16)
	}
	var ret3 int32
	if a := ret.Get(3); a != nil {
		ret3 = a.(int32)
	}
	var ret4 int64
	if a := ret.Get(4); a != nil {
		ret4 = a.(int64)
	}
	return ret0, ret1, ret2, ret3, ret4
}

// Interface implements (example.Primitives).Interface
func (r *PrimitivesMock) Interface(p0 interface{}) interface{} {
	ret := r.Called(p0)
	var ret0 interface{}
	if a := ret.Get(0); a != nil {
		ret0 = a.(interface{})
	}
	return ret0
}

// Map implements (example.Primitives).Map
func (r *PrimitivesMock) Map(p0 map[int]bool) map[int]bool {
	ret := r.Called(p0)
	var ret0 map[int]bool
	if a := ret.Get(0); a != nil {
		ret0 = a.(map[int]bool)
	}
	return ret0
}

// Ptr implements (example.Primitives).Ptr
func (r *PrimitivesMock) Ptr(p0 uintptr) uintptr {
	ret := r.Called(p0)
	var ret0 uintptr
	if a := ret.Get(0); a != nil {
		ret0 = a.(uintptr)
	}
	return ret0
}

// Slice implements (example.Primitives).Slice
func (r *PrimitivesMock) Slice(p0 []int) []bool {
	ret := r.Called(p0)
	var ret0 []bool
	if a := ret.Get(0); a != nil {
		ret0 = a.([]bool)
	}
	return ret0
}

// Uint implements (example.Primitives).Uint
func (r *PrimitivesMock) Uint(p0 uint, p1 uint8, p2 uint16, p3 uint32, p4 uint64) (uint, uint8, uint16, uint32, uint64) {
	ret := r.Called(p0, p1, p2, p3, p4)
	var ret0 uint
	if a := ret.Get(0); a != nil {
		ret0 = a.(uint)
	}
	var ret1 uint8
	if a := ret.Get(1); a != nil {
		ret1 = a.(uint8)
	}
	var ret2 uint16
	if a := ret.Get(2); a != nil {
		ret2 = a.(uint16)
	}
	var ret3 uint32
	if a := ret.Get(3); a != nil {
		ret3 = a.(uint32)
	}
	var ret4 uint64
	if a := ret.Get(4); a != nil {
		ret4 = a.(uint64)
	}
	return ret0, ret1, ret2, ret3, ret4
}

// VariadicSlice implements (example.Primitives).VariadicSlice
func (r *PrimitivesMock) VariadicSlice(p0 ...[][][]string) {
	r.Called(p0)
	return
}

// World implements (example.Primitives).World
func (r *PrimitivesMock) World(p0 int, p1 ...string) map[int]interface{} {
	ret := r.Called(p0, p1)
	var ret0 map[int]interface{}
	if a := ret.Get(0); a != nil {
		ret0 = a.(map[int]interface{})
	}
	return ret0
}

// NewFunctionsMock creates a new FunctionsMock
func NewFunctionsMock() *FunctionsMock { return &FunctionsMock{} }

// FunctionsMock implements example.Functions
type FunctionsMock struct{ mock.Mock }

// FunctionsRun implements (example.Functions).FunctionsRun
func (r *FunctionsMock) FunctionsRun(p0 int, p1 []int, p2 ...[][]int) (byte, []byte, [][]byte) {
	ret := r.Called(p0, p1, p2)
	var ret0 byte
	if a := ret.Get(0); a != nil {
		ret0 = a.(byte)
	}
	var ret1 []byte
	if a := ret.Get(1); a != nil {
		ret1 = a.([]byte)
	}
	var ret2 [][]byte
	if a := ret.Get(2); a != nil {
		ret2 = a.([][]byte)
	}
	return ret0, ret1, ret2
}
