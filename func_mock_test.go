package funcmock_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ktsivkov/funcmock"
)

var ExampleError = errors.New("test")

type TestFuncTypeNoArgumentsNoOuts func()

type TestFuncTypeArgumentsNoOuts func(a, b string)

type TestFuncTypeNoArgumentsOuts func() (string, error)

type TestFuncTypeArgumentsOuts func(a, b string) (string, error)

func TestFor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		assert.NotPanics(t, func() {
			assert.IsType(t, &funcmock.Builder[TestFuncTypeNoArgumentsNoOuts]{}, funcmock.For[TestFuncTypeNoArgumentsNoOuts]())
		})
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("not a func", func(t *testing.T) {
			assert.Panics(t, func() {
				funcmock.For[int]()
			})
		})
	})
}

func TestBuilder_Build(t *testing.T) {
	t.Run("with no arguments and with no outputs", func(t *testing.T) {
		funcMock := funcmock.For[TestFuncTypeNoArgumentsNoOuts]()
		funcMock.On().Return()

		fn := funcMock.Build()
		assert.NotPanics(t, func() {
			fn()
		})
	})
	t.Run("with arguments and with no outputs", func(t *testing.T) {
		funcMock := funcmock.For[TestFuncTypeArgumentsNoOuts]()
		funcMock.On("a", "b").Return()

		fn := funcMock.Build()
		assert.NotPanics(t, func() {
			fn("a", "b")
		})
	})
	t.Run("with no arguments and with outputs", func(t *testing.T) {
		t.Run("with no error", func(t *testing.T) {
			funcMock := funcmock.For[TestFuncTypeNoArgumentsOuts]()
			funcMock.On().Return("test", nil)

			fn := funcMock.Build()
			assert.NotPanics(t, func() {
				res, err := fn()
				assert.NoError(t, err)
				assert.Equal(t, "test", res)
			})
		})

		t.Run("with error", func(t *testing.T) {
			funcMock := funcmock.For[TestFuncTypeNoArgumentsOuts]()
			funcMock.On().Return("test", ExampleError)

			fn := funcMock.Build()
			assert.NotPanics(t, func() {
				res, err := fn()
				assert.ErrorIs(t, err, ExampleError)
				assert.Equal(t, "test", res)
			})
		})
	})
	t.Run("with arguments and with outputs", func(t *testing.T) {
		t.Run("with no error", func(t *testing.T) {
			funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
			funcMock.On("1", "2").Return("1 2", nil)

			fn := funcMock.Build()
			assert.NotPanics(t, func() {
				res, err := fn("1", "2")
				assert.NoError(t, err)
				assert.Equal(t, "1 2", res)
			})
		})
		t.Run("with error", func(t *testing.T) {
			funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
			defer funcMock.AssertNumberOfCalls(t, 1)
			funcMock.On("1", "2").Return("1 2", ExampleError)

			fn := funcMock.Build()
			assert.NotPanics(t, func() {
				res, err := fn("1", "2")
				assert.ErrorIs(t, err, ExampleError)
				assert.Equal(t, "1 2", res)
			})
		})
	})
}

func TestBuilder_On(t *testing.T) {
	var (
		givenArg1   = "a"
		givenArg2   = "b"
		expectedOut = "a"
		expectedErr = errors.New("test")
	)
	funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
	funcMock.On(givenArg1, givenArg2).Return(expectedOut, expectedErr)

	fn := funcMock.Build()
	res, err := fn(givenArg1, givenArg2)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedOut, res)
}

func TestBuilder_AssertCalled(t *testing.T) {
	var (
		givenArg1   = "a"
		givenArg2   = "b"
		expectedOut = "a"
		expectedErr = errors.New("test")
	)
	funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
	funcMock.On(givenArg1, givenArg2).Return(expectedOut, expectedErr)
	defer funcMock.AssertCalled(t, givenArg1, givenArg2)

	fn := funcMock.Build()
	res, err := fn(givenArg1, givenArg2)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedOut, res)
}

func TestBuilder_AssertNotCalled(t *testing.T) {
	var (
		givenArg1   = "a"
		givenArg2   = "b"
		expectedOut = "a"
		expectedErr = errors.New("test")
	)
	funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
	funcMock.On(givenArg1, givenArg2).Return(expectedOut, expectedErr)
	defer funcMock.AssertNotCalled(t, givenArg2, givenArg1)

	fn := funcMock.Build()
	res, err := fn(givenArg1, givenArg2)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedOut, res)
}

func TestBuilder_AssertNumberOfCalls(t *testing.T) {
	var (
		givenArg1   = "a"
		givenArg2   = "b"
		expectedOut = "a"
		expectedErr = errors.New("test")
	)
	funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
	funcMock.On(givenArg1, givenArg2).Return(expectedOut, expectedErr)
	defer funcMock.AssertNumberOfCalls(t, 1)

	fn := funcMock.Build()
	res, err := fn(givenArg1, givenArg2)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedOut, res)
}

func TestBuilder_AssertExpectations(t *testing.T) {
	var (
		givenArg1   = "a"
		givenArg2   = "b"
		expectedOut = "a"
		expectedErr = errors.New("test")
	)
	funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
	funcMock.On(givenArg1, givenArg2).Return(expectedOut, expectedErr)
	defer funcMock.AssertExpectations(t)

	fn := funcMock.Build()
	res, err := fn(givenArg1, givenArg2)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedOut, res)
}

func TestBuilder_String(t *testing.T) {
	funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
	assert.Equal(t, funcMock.String(), reflect.TypeFor[TestFuncTypeArgumentsOuts]().String())
}

func TestBuilder_TestData(t *testing.T) {
	funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
	assert.IsType(t, objx.Map{}, funcMock.TestData())
}

func TestBuilder_IsCallable(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
		funcMock.On("a", "b").Return("a", ExampleError)
		assert.True(t, funcMock.IsCallable(t, "a", "b"))
	})
	t.Run("false", func(t *testing.T) {
		funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
		assert.False(t, funcMock.IsCallable(t, "a", "b"))
	})
}

func TestBuilder_Called(t *testing.T) {
	out1, out2 := "a", ExampleError
	funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
	funcMock.On("a", "b").Return(out1, out2)
	assert.Equal(t, funcMock.Called("a", "b"), mock.Arguments{out1, out2})
}

func TestBuilder_Test(t *testing.T) {
	tt := &mockT{}
	tt.On("Errorf", mock.Anything, mock.Anything)
	tt.On("FailNow")
	defer tt.AssertExpectations(t)

	funcMock := funcmock.For[TestFuncTypeArgumentsOuts]()
	funcMock.Test(tt)
	assert.PanicsWithValue(t, examplePanicMessage, func() {
		funcMock.Called("a", "b")
	})
}
