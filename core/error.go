package core

import (
	"errors"
	"fmt"

	gojsoncore "github.com/rogonion/go-json/core"
)

/*
Error is the default base error for the metadata model package.
*/
type Error struct {
	Err          error
	FunctionName string
	Message      string
	Data         gojsoncore.JsonObject

	defaultBaseError error
}

func (e *Error) SetDefaultBaseError(value error) {
	e.defaultBaseError = value
}

func (e *Error) WithDefaultBaseError(value error) *Error {
	e.SetDefaultBaseError(value)
	return e
}

func (e *Error) Error() string {
	var err error
	if e.Message != "" {
		err = errors.New(e.Message)
	}
	return fmt.Errorf("%w: %w", err, e.Err).Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) String() string {
	str := e.Error()
	if e.FunctionName != "" {
		str = str + " \nFunctionName: " + e.FunctionName
	}
	str = str + "\nMessage: " + e.Message
	if e.Data != nil {
		str = str + " \nData: " + e.Data.String()
	}
	return str
}

func (e *Error) WithData(value gojsoncore.JsonObject) *Error {
	e.Data = value
	return e
}

func (e *Error) WithNestedError(value error) *Error {
	e.Err = fmt.Errorf("%w: %w", e.defaultBaseError, value)
	return e
}

func (e *Error) WithFunctionName(value string) *Error {
	e.FunctionName = value
	return e
}

func (e *Error) WithMessage(value string) *Error {
	e.Message = value
	return e
}

func NewError() *Error {
	n := new(Error)
	n.defaultBaseError = errors.New("metadata model error")
	return n
}
