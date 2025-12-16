package fieldcolumns

import (
	"errors"
	"fmt"

	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
)

type FieldColumnPosition struct {
	SourceIndex           int
	FieldGroupJsonPathKey path.JSONPath

	FieldViewInSeparateColumns                  bool
	FieldViewValuesInSeparateColumnsHeaderIndex int

	GroupViewInSeparateColumns                  bool
	GroupViewValuesInSeparateColumnsHeaderIndex int
	GroupViewParentJsonPathKey                  path.JSONPath
	FieldJsonPathKeySuffix                      string

	FieldGroupPositionBefore bool
}

func (n *FieldColumnPosition) String() string {
	if n.GroupViewInSeparateColumns {
		return string(n.GroupViewParentJsonPathKey) + path.JsonpathLeftBracket + fmt.Sprintf("%d", n.GroupViewValuesInSeparateColumnsHeaderIndex) + path.JsonpathRightBracket + path.JsonpathDotNotation + n.FieldJsonPathKeySuffix
	}

	if n.FieldViewInSeparateColumns {
		return string(n.FieldGroupJsonPathKey) + path.JsonpathLeftBracket + fmt.Sprintf("%d", n.FieldViewValuesInSeparateColumnsHeaderIndex) + path.JsonpathRightBracket
	}

	return string(n.FieldGroupJsonPathKey)
}

func (n *FieldColumnPosition) JSONPath() path.JSONPath {
	return path.JSONPath(n.String())
}

type RepositionFieldColumns []FieldColumnPosition

var (
	// ErrFieldColumnsError default error for field columns module.
	ErrFieldColumnsError = errors.New("field columns error")
)

func NewError() *core.Error {
	n := core.NewError().WithDefaultBaseError(ErrFieldColumnsError)
	return n
}
