package fieldcolumns

import (
	"errors"

	"github.com/rogonion/go-metadatamodel/core"
)

var (
	// ErrFieldColumnsError default error for field columns module.
	ErrFieldColumnsError = errors.New("field columns error")
)

func NewError() *core.Error {
	n := core.NewError().WithDefaultBaseError(ErrFieldColumnsError)
	return n
}
