package example

// Code generated by Mockhiato. DO NOT EDIT.
import (
	mock "github.com/stretchr/testify/mock"
)

// NewTargetMock creates a new TargetMock
func NewTargetMock() *TargetMock { return &TargetMock{} }

// TargetMock implements example.Target
type TargetMock struct{ mock.Mock }

// Bool implements (example.Target).Bool
func (r *TargetMock) Bool(p0 bool) bool {
	ret := r.Called(p0)
	var ret0 bool
	if a := ret.Get(0); a != nil {
		ret0 = a.(bool)
	}
	return ret0
}

// Chan implements (example.Target).Chan
func (r *TargetMock) Chan(p0 chan int) chan int {
	ret := r.Called(p0)
	var ret0 chan int
	if a := ret.Get(0); a != nil {
		ret0 = a.(chan int)
	}
	return ret0
}

// Complex implements (example.Target).Complex
func (r *TargetMock) Complex(p0 complex64, p1 complex128) (complex64, complex128) {
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

// Error implements (example.Target).Error
func (r *TargetMock) Error(p0 error) error {
	ret := r.Called(p0)
	var ret0 error
	if a := ret.Get(0); a != nil {
		ret0 = a.(error)
	}
	return ret0
}

// Float implements (example.Target).Float
func (r *TargetMock) Float(p0 float32, p1 float64) (float32, float64) {
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

// Func implements (example.Target).Func
func (r *TargetMock) Func(p0 func(int) float32) func(int) float32 {
	ret := r.Called(p0)
	var ret0 func(int) float32
	if a := ret.Get(0); a != nil {
		ret0 = a.(func(int) float32)
	}
	return ret0
}

// Int implements (example.Target).Int
func (r *TargetMock) Int(p0 int, p1 int8, p2 int16, p3 int32, p4 int64) (int, int8, int16, int32, int64) {
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

// Interface implements (example.Target).Interface
func (r *TargetMock) Interface(p0 interface{}) interface{} {
	ret := r.Called(p0)
	var ret0 interface{}
	if a := ret.Get(0); a != nil {
		ret0 = a.(interface{})
	}
	return ret0
}

// Map implements (example.Target).Map
func (r *TargetMock) Map(p0 map[int]error) map[int]error {
	ret := r.Called(p0)
	var ret0 map[int]error
	if a := ret.Get(0); a != nil {
		ret0 = a.(map[int]error)
	}
	return ret0
}

// Ptr implements (example.Target).Ptr
func (r *TargetMock) Ptr(p0 uintptr) uintptr {
	ret := r.Called(p0)
	var ret0 uintptr
	if a := ret.Get(0); a != nil {
		ret0 = a.(uintptr)
	}
	return ret0
}

// Slice implements (example.Target).Slice
func (r *TargetMock) Slice(p0 []int) []int {
	ret := r.Called(p0)
	var ret0 []int
	if a := ret.Get(0); a != nil {
		ret0 = a.([]int)
	}
	return ret0
}

// Text implements (example.Target).Text
func (r *TargetMock) Text(p0 byte, p1 rune, p2 string) (byte, rune, string) {
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

// Uint implements (example.Target).Uint
func (r *TargetMock) Uint(p0 uint, p1 uint8, p2 uint16, p3 uint32, p4 uint64) (uint, uint8, uint16, uint32, uint64) {
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
