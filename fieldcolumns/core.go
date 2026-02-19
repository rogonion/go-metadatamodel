package fieldcolumns

import (
	"errors"
	"fmt"

	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
)

// FieldColumnPosition represents the position and context of a field within the metadata model hierarchy.
// It tracks the JSON path, view settings (e.g., pivoting), and relative positioning.
type FieldColumnPosition struct {
	SourceIndex           int
	FieldGroupJsonPathKey path.JSONPath

	// FieldViewInSeparateColumns indicates if individual entries in a single field should be viewed as separate columns.
	FieldViewInSeparateColumns                  bool
	FieldViewValuesInSeparateColumnsHeaderIndex int

	// For fields in 1D groups that should be viewed in separate columns.
	GroupViewInSeparateColumns                  bool
	GroupViewValuesInSeparateColumnsHeaderIndex int
	GroupViewParentJsonPathKey                  path.JSONPath
	FieldJsonPathKeySuffix                      string

	FieldGroupPositionBefore bool
}

// String returns the string representation of the field's path, including pivot indices if applicable.
func (n *FieldColumnPosition) String() string {
	if n.GroupViewInSeparateColumns {
		return string(n.GroupViewParentJsonPathKey) + path.JsonpathDotNotation + core.GroupFields + path.JsonpathLeftBracket + fmt.Sprintf("%d", n.GroupViewValuesInSeparateColumnsHeaderIndex) + path.JsonpathRightBracket + path.JsonpathDotNotation + n.FieldJsonPathKeySuffix
	}

	if n.FieldViewInSeparateColumns {
		return string(n.FieldGroupJsonPathKey) + path.JsonpathLeftBracket + fmt.Sprintf("%d", n.FieldViewValuesInSeparateColumnsHeaderIndex) + path.JsonpathRightBracket
	}

	return string(n.FieldGroupJsonPathKey)
}

// JSONPath returns the FieldColumnPosition as a path.JSONPath.
func (n *FieldColumnPosition) JSONPath() path.JSONPath {
	return path.JSONPath(n.String())
}

// FieldsColumnsPositions is a slice of FieldColumnPosition pointers.
type FieldsColumnsPositions []*FieldColumnPosition

// RepositionFieldColumns is a slice of FieldColumnPosition used for reordering.
type RepositionFieldColumns []FieldColumnPosition

var (
	// ErrFieldColumnsError default error for field columns module.
	ErrFieldColumnsError = errors.New("field columns error")
)

// NewError creates a new core.Error with the default field columns error base.
func NewError() *core.Error {
	n := core.NewError().WithDefaultBaseError(ErrFieldColumnsError)
	return n
}
