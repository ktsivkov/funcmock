package funcmock

import (
	"reflect"

	"github.com/stretchr/objx"
	"github.com/stretchr/testify/mock"
)

func For[T any]() *Builder[T] {
	typ := reflect.TypeFor[T]()
	if typ == nil || typ.Kind() != reflect.Func {
		panic("type must be a function")
	}

	return &Builder[T]{
		mock: mock.Mock{},
	}
}

func As[T any](_ T) *Builder[T] {
	typ := reflect.TypeFor[T]()
	if typ == nil || typ.Kind() != reflect.Func {
		panic("type must be a function")
	}

	return &Builder[T]{
		mock: mock.Mock{},
	}
}

type Builder[T any] struct {
	mock mock.Mock
}

func (m *Builder[T]) Build() T {
	typ := reflect.TypeFor[T]()
	return reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
		argsAsInterface := make([]interface{}, len(args))
		for i, arg := range args {
			argsAsInterface[i] = arg.Interface()
		}
		outs := m.Called(argsAsInterface...)
		res := make([]reflect.Value, typ.NumOut())
		for i := 0; i < typ.NumOut(); i++ {
			val := outs.Get(i)
			if val == nil {
				res[i] = reflect.Zero(typ.Out(i))
				continue
			}
			res[i] = reflect.ValueOf(val)
		}
		return res
	}).Interface().(T)
}

func (m *Builder[T]) On(args ...interface{}) *mock.Call {
	return m.mock.On("func", args...)
}

func (m *Builder[T]) AssertExpectations(t mock.TestingT) {
	m.mock.AssertExpectations(t)
}

func (m *Builder[T]) AssertNotCalled(t mock.TestingT, arguments ...interface{}) {
	m.mock.AssertNotCalled(t, "func", arguments...)
}

func (m *Builder[T]) AssertCalled(t mock.TestingT, arguments ...interface{}) {
	m.mock.AssertCalled(t, "func", arguments...)
}

func (m *Builder[T]) AssertNumberOfCalls(t mock.TestingT, expectedCalls int) {
	m.mock.AssertNumberOfCalls(t, "func", expectedCalls)
}

func (m *Builder[T]) IsCallable(t mock.TestingT, arguments ...interface{}) bool {
	return m.mock.IsMethodCallable(t, "func", arguments...)
}

func (m *Builder[T]) Test(t mock.TestingT) {
	m.mock.Test(t)
}

func (m *Builder[T]) Called(arguments ...interface{}) mock.Arguments {
	return m.mock.MethodCalled("func", arguments...)
}

func (m *Builder[T]) TestData() objx.Map {
	return m.mock.TestData()
}

func (m *Builder[T]) String() string {
	return reflect.TypeFor[T]().String()
}
