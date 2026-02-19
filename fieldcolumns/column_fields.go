package fieldcolumns

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/core"
)

// GetColumnFieldByIndexInUnskippedReadOrder retrieves a ColumnField by its index in the unskipped read order.
func (n *ColumnFields) GetColumnFieldByIndexInUnskippedReadOrder(index int) (*ColumnField, bool) {
	if index >= 0 && index < len(n.UnskippedReadOrderOfColumnFields) {
		return n.GetColumnFieldByIndexInOriginalReadOrder(n.UnskippedReadOrderOfColumnFields[index])
	}
	return nil, false
}

// GetColumnFieldByIndexInRepositionedReadOrder retrieves a ColumnField by its index in the repositioned read order.
func (n *ColumnFields) GetColumnFieldByIndexInRepositionedReadOrder(index int) (*ColumnField, bool) {
	if index >= 0 && index < len(n.RepositionedReadOrderOfColumnFields) {
		return n.GetColumnFieldByIndexInOriginalReadOrder(n.RepositionedReadOrderOfColumnFields[index])
	}
	return nil, false
}

// GetColumnFieldByIndexInOriginalReadOrder retrieves a ColumnField by its index in the original extraction order.
func (n *ColumnFields) GetColumnFieldByIndexInOriginalReadOrder(index int) (*ColumnField, bool) {
	if index >= 0 && index < len(n.OriginalReadOrderOfColumnFields) {
		return n.GetColumnFieldByFieldGroupJsonPathKey(n.OriginalReadOrderOfColumnFields[index].JSONPath())
	}
	return nil, false
}

// GetColumnFieldByFieldGroupJsonPathKey retrieves a ColumnField by its JSON path key.
func (n *ColumnFields) GetColumnFieldByFieldGroupJsonPathKey(jsonPathKey path.JSONPath) (*ColumnField, bool) {
	field, ok := n.Fields[jsonPathKey]
	return field, ok
}

/*
Skip a field if skip.FirstMatch returns `true` or add.FirstMatch returns `false`.

Call after ColumnFields.Reposition.

Populates ColumnFields.UnskippedReadOrderOfColumnFields.

Updates ColumnField.IndexInUnskippedColumnFields.
*/
func (n *ColumnFields) Skip(skip core.FieldGroupPropertiesMatch, add core.FieldGroupPropertiesMatch) {
	n.FieldsToSkip = make(FieldsToSkip)

	n.UnskippedReadOrderOfColumnFields = make([]int, 0)
	for _, originalIndex := range n.RepositionedReadOrderOfColumnFields {
		if field, ok := n.GetColumnFieldByIndexInOriginalReadOrder(originalIndex); ok {
			field.Skip = false
			if skip.IsValid() {
				if skip.FirstMatch(field.Property) {
					field.Skip = true
				}
			}
			if add.IsValid() {
				if !add.FirstMatch(field.Property) {
					field.Skip = true
				}
			}
			if field.Skip {
				n.FieldsToSkip[field.FieldColumnPosition.JSONPath()] = FieldToSkip()
				field.IndexInUnskippedColumnFields = -1
			} else {
				n.UnskippedReadOrderOfColumnFields = append(n.UnskippedReadOrderOfColumnFields, originalIndex)
				field.IndexInUnskippedColumnFields = len(n.UnskippedReadOrderOfColumnFields) - 1
			}
		}
	}

}

/*
Reposition reorganizes ColumnFields.RepositionedReadOrderOfColumnFields based on ColumnFields.RepositionFieldColumns information.

Only call this method after Extraction.Extract has been run successfully.

Populates ColumnFields.RepositionedReadOrderOfColumnFields.

Updates ColumnField.IndexInRepositionedColumnFields.
*/
func (n *ColumnFields) Reposition() {
	totalNoOfFields := len(n.RepositionedReadOrderOfColumnFields)

	for _, newPosition := range n.RepositionFieldColumns {
		if destinationField, ok := n.Fields[newPosition.JSONPath()]; ok {
			sourceIndex := -1
			destinationIndex := -1
			for cIndex, cValue := range n.RepositionedReadOrderOfColumnFields {
				if cValue == destinationField.IndexInOriginalReadOrderOfColumnFields {
					if newPosition.FieldGroupPositionBefore || cIndex >= totalNoOfFields-1 {
						destinationIndex = cIndex
					} else {
						destinationIndex = cIndex + 1
					}
				} else {
					if cValue == newPosition.SourceIndex {
						sourceIndex = cIndex
					}
				}
				if destinationIndex >= 0 && sourceIndex >= 0 {
					break
				}
			}

			if destinationIndex >= 0 && sourceIndex >= 0 && destinationIndex != sourceIndex {
				n.RepositionedReadOrderOfColumnFields = append(n.RepositionedReadOrderOfColumnFields[:sourceIndex], n.RepositionedReadOrderOfColumnFields[sourceIndex+1:]...)
				n.RepositionedReadOrderOfColumnFields = append(n.RepositionedReadOrderOfColumnFields[:destinationIndex], append([]int{newPosition.SourceIndex}, n.RepositionedReadOrderOfColumnFields[destinationIndex:]...)...)
			}
		}
	}

	for indexOfField, jsonPathKey := range n.OriginalReadOrderOfColumnFields {
		if columnField, ok := n.GetColumnFieldByFieldGroupJsonPathKey(jsonPathKey.JSONPath()); ok {
			for readIndex, indexOfFieldInReadOrder := range n.RepositionedReadOrderOfColumnFields {
				if indexOfFieldInReadOrder == indexOfField {
					columnField.IndexInRepositionedColumnFields = readIndex
				}
			}
		}
	}
}

// NewColumnFields creates a new ColumnFields instance.
func NewColumnFields() *ColumnFields {
	return &ColumnFields{
		Fields:                              make(Fields),
		OriginalReadOrderOfColumnFields:     make(FieldsColumnsPositions, 0),
		RepositionedReadOrderOfColumnFields: make([]int, 0),
		FieldsToSkip:                        make(FieldsToSkip),
		RepositionFieldColumns:              make(RepositionFieldColumns, 0),
	}
}

/*
ColumnFields represents the metadata model fields as columns in a table.
*/
type ColumnFields struct {
	// Fields store field information.
	Fields Fields

	// OriginalReadOrderOfColumnFields store order of Fields as per read order of metadata model tree.
	OriginalReadOrderOfColumnFields FieldsColumnsPositions

	// Derived from OriginalReadOrderOfColumnFields
	//
	// The value stored at each index is the actual index in the OriginalReadOrderOfColumnFields.
	//
	// Its size MUST be equal to OriginalReadOrderOfColumnFields.
	RepositionedReadOrderOfColumnFields []int

	// Derived from RepositionedReadOrderOfColumnFields
	//
	// The value stored at each index is the actual index in the OriginalReadOrderOfColumnFields.
	//
	// Its size may not be equal to OriginalReadOrderOfColumnFields.
	UnskippedReadOrderOfColumnFields []int

	// FieldsToSkip
	FieldsToSkip FieldsToSkip

	// RepositionFieldColumns for repositioning OriginalReadOrderOfColumnFields to RepositionedReadOrderOfColumnFields.
	RepositionFieldColumns RepositionFieldColumns
}

// FieldsToSkip is a set of JSON paths representing fields that should be skipped.
type FieldsToSkip map[path.JSONPath]struct{}

// FieldToSkip returns an empty struct, used as a value in the FieldsToSkip set.
func FieldToSkip() struct{} {
	return struct{}{}
}

// Fields is a map of JSON paths to ColumnField pointers.
type Fields map[path.JSONPath]*ColumnField

// ColumnField represents a single field (column) extracted from the metadata model.
type ColumnField struct {
	FieldColumnPosition FieldColumnPosition
	// Property field metadata model property.
	Property gojsoncore.JsonObject
	// Schema Only set if Extraction.schema is set.
	Schema schema.Schema
	// IndexInOriginalReadOrderOfColumnFields original index of field as it was being Extraction.Extract.
	IndexInOriginalReadOrderOfColumnFields int
	// IndexInRepositionedColumnFields new index of field after ColumnFields.Reposition.
	IndexInRepositionedColumnFields int
	// Skip field. Set using ColumnFields.Skip.
	Skip bool
	// IndexInUnskippedColumnFields new index of field after ColumnFields.Skip.
	IndexInUnskippedColumnFields int
}
