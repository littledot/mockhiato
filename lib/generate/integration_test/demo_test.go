package example_test

import (
	"errors"
	"testing"

	. "github.com/littledot/mockhiato/lib/generate/integration_test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestifyIntegration_Bool(t *testing.T) {
	target := &TargetMock{}

	a := true
	target.On("Bool", a).Return(a)
	ans := target.Bool(a)
	assert.Equal(t, a, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilBool(t *testing.T) {
	target := &TargetMock{}

	var a bool
	target.On("Bool", true).Return(a)
	ans := target.Bool(true)
	assert.Zero(t, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Chan(t *testing.T) {
	target := &TargetMock{}

	a := make(chan int, 1)
	target.On("Chan", a).Return(a)
	ans := target.Chan(a)
	assert.Equal(t, a, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilChan(t *testing.T) {
	target := &TargetMock{}

	a := make(chan int, 1)
	target.On("Chan", a).Return(nil)
	ans := target.Chan(a)
	assert.Zero(t, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Complex(t *testing.T) {
	target := &TargetMock{}

	a := complex(float32(1.0), 2.0)
	b := complex(float64(3.0), 4.0)
	target.On("Complex", a, b).Return(a, b)
	ans, ans1 := target.Complex(a, b)
	assert.Equal(t, a, ans)
	assert.Equal(t, b, ans1)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilComplex(t *testing.T) {
	target := &TargetMock{}

	a := complex(float32(1.0), 2.0)
	b := complex(float64(3.0), 4.0)
	target.On("Complex", a, b).Return(complex64(0), complex128(0))
	ans, ans1 := target.Complex(a, b)
	assert.Zero(t, ans)
	assert.Zero(t, ans1)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Error(t *testing.T) {
	target := &TargetMock{}

	a := errors.New("test nonnil")
	target.On("Error", a).Return(a)
	ans := target.Error(a)
	assert.Equal(t, a, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilError(t *testing.T) {
	target := &TargetMock{}

	a := errors.New("test nil")
	target.On("Error", a).Return(nil)
	ans := target.Error(a)
	assert.Zero(t, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Float(t *testing.T) {
	target := &TargetMock{}

	a := float32(1.0)
	b := float64(2.0)
	target.On("Float", a, b).Return(a, b)
	ans, ans1 := target.Float(a, b)
	assert.Equal(t, a, ans)
	assert.Equal(t, b, ans1)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilFloat(t *testing.T) {
	target := &TargetMock{}

	a := float32(1.0)
	b := float64(2.0)
	target.On("Float", a, b).Return(float32(0), float64(0))
	ans, ans1 := target.Float(a, b)
	assert.Zero(t, ans)
	assert.Zero(t, ans1)

	target.AssertExpectations(t)
}

func TestTestifyIntegrationFunc(t *testing.T) {
	target := &TargetMock{}

	a := func(int) float32 { return 1.0 }
	target.On("Func", mock.AnythingOfType("func(int) float32")).Return(a)
	ans := target.Func(a)(0)
	assert.Equal(t, float32(1.0), ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilFunc(t *testing.T) {
	target := &TargetMock{}

	a := func(int) float32 { return 1.0 }
	target.On("Func", mock.AnythingOfType("func(int) float32")).Return(nil)
	ans := target.Func(a)
	assert.Zero(t, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Int(t *testing.T) {
	target := &TargetMock{}

	a := int(1.0)
	b := int8(2.0)
	c := int16(2.0)
	d := int32(2.0)
	e := int64(2.0)
	target.On("Int", a, b, c, d, e).Return(a, b, c, d, e)
	ans, ans1, ans2, ans3, ans4 := target.Int(a, b, c, d, e)
	assert.Equal(t, a, ans)
	assert.Equal(t, b, ans1)
	assert.Equal(t, c, ans2)
	assert.Equal(t, d, ans3)
	assert.Equal(t, e, ans4)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_ZeroInt(t *testing.T) {
	target := &TargetMock{}

	a := int(1.0)
	b := int8(2.0)
	c := int16(2.0)
	d := int32(2.0)
	e := int64(2.0)
	target.On("Int", a, b, c, d, e).Return(int(0), int8(0), int16(0), int32(0), int64(0))
	ans, ans1, ans2, ans3, ans4 := target.Int(a, b, c, d, e)
	assert.Zero(t, ans)
	assert.Zero(t, ans1)
	assert.Zero(t, ans2)
	assert.Zero(t, ans3)
	assert.Zero(t, ans4)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Interface(t *testing.T) {
	target := &TargetMock{}

	target.On("Interface", 1).Return(1)
	ans := target.Interface(1)
	assert.Equal(t, 1, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilInterface(t *testing.T) {
	target := &TargetMock{}

	target.On("Interface", nil).Return(nil)
	ans := target.Interface(nil)
	assert.Zero(t, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Map(t *testing.T) {
	target := &TargetMock{}

	a := map[int]error{}
	target.On("Map", a).Return(a)
	ans := target.Map(a)
	assert.Equal(t, a, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilMap(t *testing.T) {
	target := &TargetMock{}

	a := map[int]error{}
	target.On("Map", a).Return(nil)
	ans := target.Map(a)
	assert.Zero(t, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Ptr(t *testing.T) {
	target := &TargetMock{}

	a := uintptr(1)
	target.On("Ptr", a).Return(a)
	ans := target.Ptr(a)
	assert.Equal(t, a, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilPtr(t *testing.T) {
	target := &TargetMock{}

	a := uintptr(1)
	target.On("Ptr", a).Return(uintptr(0))
	ans := target.Ptr(a)
	assert.Zero(t, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Slice(t *testing.T) {
	target := &TargetMock{}

	a := []int{1}
	target.On("Slice", a).Return(a)
	ans := target.Slice(a)
	assert.Equal(t, a, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilSlice(t *testing.T) {
	target := &TargetMock{}

	a := []int{1}
	target.On("Slice", a).Return(nil)
	ans := target.Slice(a)
	assert.Zero(t, ans)

	target.AssertExpectations(t)
}

func TestTestifyIntegrationText(t *testing.T) {
	target := &TargetMock{}

	a := byte(0)
	b := rune(0)
	c := ""
	target.On("Text", a, b, c).Return(a, b, c)
	ans, ans1, ans2 := target.Text(a, b, c)
	assert.Zero(t, ans)
	assert.Zero(t, ans1)
	assert.Zero(t, ans2)

	a = byte('a')
	b = rune('b')
	c = "c"
	target.On("Text", a, b, c).Return(a, b, c)
	ans, ans1, ans2 = target.Text(a, b, c)
	assert.Equal(t, a, ans)
	assert.Equal(t, b, ans1)
	assert.Equal(t, c, ans2)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_Uint(t *testing.T) {
	target := &TargetMock{}

	a := uint(1.0)
	b := uint8(2.0)
	c := uint16(2.0)
	d := uint32(2.0)
	e := uint64(2.0)
	target.On("Uint", a, b, c, d, e).Return(a, b, c, d, e)
	ans, ans1, ans2, ans3, ans4 := target.Uint(a, b, c, d, e)
	assert.Equal(t, a, ans)
	assert.Equal(t, b, ans1)
	assert.Equal(t, c, ans2)
	assert.Equal(t, d, ans3)
	assert.Equal(t, e, ans4)

	target.AssertExpectations(t)
}

func TestTestifyIntegration_NilUint(t *testing.T) {
	target := &TargetMock{}

	a := uint(1.0)
	b := uint8(2.0)
	c := uint16(2.0)
	d := uint32(2.0)
	e := uint64(2.0)
	target.On("Uint", a, b, c, d, e).Return(uint(0), uint8(0), uint16(0), uint32(0), uint64(0))
	ans, ans1, ans2, ans3, ans4 := target.Uint(a, b, c, d, e)
	assert.Zero(t, ans)
	assert.Zero(t, ans1)
	assert.Zero(t, ans2)
	assert.Zero(t, ans3)
	assert.Zero(t, ans4)

	target.AssertExpectations(t)
}
