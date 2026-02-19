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

// SetDefaultBaseError sets the default base error.
func (e *Error) SetDefaultBaseError(value error) {
	e.defaultBaseError = value
}

// WithDefaultBaseError sets the default base error and returns the Error instance.
func (e *Error) WithDefaultBaseError(value error) *Error {
	e.SetDefaultBaseError(value)
	return e
}

// Error returns the error message.
func (e *Error) Error() string {
	var err error
	if e.Message != "" {
		err = errors.New(e.Message)
	}
	return fmt.Errorf("%w: %w", err, e.Err).Error()
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}

// String returns the string representation of the error.
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

// WithData sets the data for the error and returns the Error instance.
func (e *Error) WithData(value gojsoncore.JsonObject) *Error {
	e.Data = value
	return e
}

// WithNestedError sets the nested error and returns the Error instance.
func (e *Error) WithNestedError(value error) *Error {
	e.Err = fmt.Errorf("%w: %w", e.defaultBaseError, value)
	return e
}

// WithFunctionName sets the function name and returns the Error instance.
func (e *Error) WithFunctionName(value string) *Error {
	e.FunctionName = value
	return e
}

// WithMessage sets the message and returns the Error instance.
func (e *Error) WithMessage(value string) *Error {
	e.Message = value
	return e
}

// NewError creates a new Error instance.
func NewError() *Error {
	n := new(Error)
	n.defaultBaseError = errors.New("metadata model error")
	return n
}
