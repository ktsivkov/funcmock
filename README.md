MockFunc - Overview and Disclaimer
========
[![Build Status](https://github.com/ktsivkov/funcmock/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/ktsivkov/funcmock/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/ktsivkov/funcmock)](https://goreportcard.com/report/github.com/ktsivkov/funcmock) [![PkgGoDev](https://pkg.go.dev/badge/github.com/ktsivkov/funcmock)](https://pkg.go.dev/github.com/ktsivkov/funcmock)

> This is a wrapper around `github.com/stretchr/testify` package.
> 
> We keep the API as close as possible to the original package.
> 
> We only aim to cover the most common use cases of function mocking.
> 
> At the moment we support GoLang versions >= 1.22.

Installation
============
```bash
go get github.com/ktsivkov/funcmock
```

Usage
=====

## Using a concrete type
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

## Using an assumed type
```go
package example

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ktsivkov/funcmock"
)

type SimpleFunction func(string, string) (string, error)

func MyFunction(a string, b string) (string, error) {
	return fmt.Sprintf("%s_%s", a, b), nil
}

type Service struct {
	simpleFunction SimpleFunction
}

func (s Service) Do() (string, error) {
	return s.simpleFunction("a", "b")
}

func TestService(t *testing.T) {
	expectedErr := errors.New("test")
	expectedOut := "a"

	fnBuilder := funcmock.As(MyFunction)
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
