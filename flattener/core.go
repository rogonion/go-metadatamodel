package flattener

import (
	"errors"

	"github.com/rogonion/go-metadatamodel/core"
)

var (
	// ErrFlattenError default error for field columns module.
	ErrFlattenError = errors.New("flattening encountered an error")

	// ErrNoGroupFields for when FieldGroupConversion.GroupFields is empty if field is a group.
	ErrNoGroupFields = errors.New("no group fields to extract found")
)

func NewError() *core.Error {
	n := core.NewError().WithDefaultBaseError(ErrFlattenError)
	return n
}
