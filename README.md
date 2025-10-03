MockFunc - Overview and Disclaimer
========

> This is a wrapper around `github.com/stretchr/testify` package.
> 
> We keep the API as close as possible to the original package.
> 
> We only aim to cover the most common use cases of function mocking.

Installation
============
```bash
go get github.com/ktsivkov/funcmock
```

Usage
=====

```go
package example

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ktsivkov/funcmock"
)

type SimpleFunction func(string, string) (string, error)

type Service struct {
	simpleFunction SimpleFunction
}

func (s Service) Do() (string, error) {
	return s.simpleFunction("a", "b")
}

func TestService(t *testing.T) {
	expectedErr := errors.New("test")
	expectedOut := "a"

	fnBuilder := funcmock.For[SimpleFunction]()
	fnBuilder.On("a", "b").Return(expectedOut, expectedErr)
	defer fnBuilder.AssertNumberOfCalls(t, 1)
	defer fnBuilder.AssertExpectations(t)

	fn := fnBuilder.Build()
	service := Service{fn}

	res, err := service.Do()
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, expectedOut, res)
}

```
