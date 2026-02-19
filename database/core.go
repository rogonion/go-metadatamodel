package database

import (
	"errors"

	"github.com/rogonion/go-metadatamodel/core"
)

var (
	//ErrDatabaseError default error for the database functions.
	ErrDatabaseError = errors.New("database error")

	//ErrDatabaseGetColumnFieldsError for when GetColumnFields.Get fails.
	ErrDatabaseGetColumnFieldsError = errors.New("database get column fields error")

	//ErrDatabaseFieldValueError for when FieldValue methods fails.
	ErrDatabaseFieldValueError = errors.New("database manipulate field value error")
)

// NewError creates a new core.Error with the default database error base.
func NewError() *core.Error {
	n := core.NewError().WithDefaultBaseError(ErrDatabaseError)
	return n
}
