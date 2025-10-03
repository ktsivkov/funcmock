package funcmock_test

import "github.com/stretchr/testify/mock"

const examplePanicMessage = "example panic message"

type mockT struct {
	mock.Mock
}

func (m *mockT) Logf(format string, args ...interface{}) {
	_ = m.Called(format, args)
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	_ = m.Called(format, args)
}

func (m *mockT) FailNow() {
	_ = m.Called()
	panic(examplePanicMessage)
}
