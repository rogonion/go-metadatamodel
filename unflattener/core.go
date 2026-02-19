package unflattener

import (
	"errors"

	"github.com/rogonion/go-metadatamodel/core"
)

var (
	// ErrFlattenError default error for field columns module.
	ErrFlattenError = errors.New("unflattening encountered an error")

	// ErrNoGroupFields for when RecursiveGroupIndexTree.MmPropGroupFields is empty if field is a group.
	ErrNoGroupFields = errors.New("no group fields to extract found")
)

// NewError creates a new core.Error with the default unflatten error base.
func NewError() *core.Error {
	n := core.NewError().WithDefaultBaseError(ErrFlattenError)
	return n
}
