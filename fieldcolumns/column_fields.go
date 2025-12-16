package fieldcolumns

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/core"
)

/*
GetCurrentIndexOfReadOrderOfFields after ColumnFields.Reposition and ColumnFields.Skip.

Removes skipped fields.
*/
func (n *ColumnFields) GetCurrentIndexOfReadOrderOfFields() []int {
	readOrder := make([]int, 0)

	for _, fieldIndex := range n.CurrentIndexOfReadOrderOfColumnFields {
		if field, ok := n.GetColumnFieldByIndexInOriginalReadOrder(fieldIndex); ok {
			if !field.Skip {
				readOrder = append(readOrder, fieldIndex)
			}
		}
	}

	return readOrder
}

func (n *ColumnFields) GetColumnFieldByFieldGroupJsonPathKey(jsonPathKey path.JSONPath) (*ColumnField, bool) {
	field, ok := n.Fields[jsonPathKey]
	return field, ok
}

func (n *ColumnFields) GetColumnFieldByIndexInCurrentReadOrder(index int) (*ColumnField, bool) {
	field, ok := n.Fields[n.ReadOrderOfColumnFields[n.CurrentIndexOfReadOrderOfColumnFields[index]].JSONPath()]
	return field, ok
}

func (n *ColumnFields) GetColumnFieldByIndexInOriginalReadOrder(index int) (*ColumnField, bool) {
	field, ok := n.Fields[n.ReadOrderOfColumnFields[index].JSONPath()]
	return field, ok
}

/*
Skip a field if skip.FirstMatch returns `true` or add.FirstMatch returns `false`.
*/
func (n *ColumnFields) Skip(skip core.FieldGroupPropertiesMatch, add core.FieldGroupPropertiesMatch) {
	n.CurrentFieldsToSkip = make(FieldsToSkip)
	for jsonPathKey, field := range n.Fields {
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
			n.CurrentFieldsToSkip[jsonPathKey] = FieldToSkip()
		}
	}
}

/*
Reposition reorganizes Extraction.CurrentIndexOfReadOrderOfColumnFields based on Extraction.RepositionFieldColumns information.

Only call this method after FieldColumns.Extract has been run successfully.

To retrieve index of fields in read order, call Extraction.GetReadOrderOfColumnFields.
*/
func (n *ColumnFields) Reposition() {
	totalNoOfFields := len(n.CurrentIndexOfReadOrderOfColumnFields)
	for _, newPosition := range n.RepositionFieldColumns {
		if destinationField, ok := n.Fields[newPosition.JSONPath()]; ok {
			sourceIndex := -1
			destinationIndex := -1
			for cIndex, cValue := range n.CurrentIndexOfReadOrderOfColumnFields {
				if cValue == destinationField.IndexInReadOrderOfColumnFields {
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
				n.CurrentIndexOfReadOrderOfColumnFields = append(n.CurrentIndexOfReadOrderOfColumnFields[:sourceIndex], n.CurrentIndexOfReadOrderOfColumnFields[sourceIndex+1:]...)
				n.CurrentIndexOfReadOrderOfColumnFields = append(n.CurrentIndexOfReadOrderOfColumnFields[:destinationIndex], append([]int{newPosition.SourceIndex}, n.CurrentIndexOfReadOrderOfColumnFields[destinationIndex:]...)...)
			}
		}
	}
}

/*
UpdateIndexInRepositionedColumnFieldsInColumnField after Extraction.Reposition, update ColumnField.IndexInRepositionedColumnFields.
*/
func (n *ColumnFields) UpdateIndexInRepositionedColumnFieldsInColumnField() {
	for indexOfField, jsonPathKey := range n.ReadOrderOfColumnFields {
		if columnField, ok := n.GetColumnFieldByFieldGroupJsonPathKey(jsonPathKey.JSONPath()); ok {
			for readIndex, indexOfFieldInReadOrder := range n.CurrentIndexOfReadOrderOfColumnFields {
				if indexOfFieldInReadOrder == indexOfField {
					columnField.IndexInRepositionedColumnFields = readIndex
				}
			}
		}
	}
}

func NewColumnFields() *ColumnFields {
	return &ColumnFields{
		Fields:                                make(Fields),
		ReadOrderOfColumnFields:               make([]*FieldColumnPosition, 0),
		CurrentIndexOfReadOrderOfColumnFields: make([]int, 0),
		CurrentFieldsToSkip:                   make(FieldsToSkip),
		RepositionFieldColumns:                make(RepositionFieldColumns, 0),
	}
}

/*
ColumnFields represents the metadata model fields as columns in a table.
*/
type ColumnFields struct {
	// Fields store field information.
	Fields Fields

	// ReadOrderOfColumnFields store order of Fields.
	ReadOrderOfColumnFields []*FieldColumnPosition

	// CurrentIndexOfReadOrderOfColumnFields
	CurrentIndexOfReadOrderOfColumnFields []int

	// CurrentFieldsToSkip
	CurrentFieldsToSkip FieldsToSkip

	// RepositionFieldColumns for repositioning CurrentIndexOfReadOrderOfColumnFields.
	RepositionFieldColumns RepositionFieldColumns
}

type FieldsToSkip map[path.JSONPath]struct{}

func FieldToSkip() struct{} {
	return struct{}{}
}

type Fields map[path.JSONPath]*ColumnField

type ColumnField struct {
	// Property field metadata model property.
	Property gojsoncore.JsonObject
	// Schema Only set if Extraction.schema is set.
	Schema schema.Schema
	// IndexInReadOrderOfColumnFields original index of field as it was being Extraction.Extract.
	IndexInReadOrderOfColumnFields int
	// IndexInRepositionedColumnFields new index of field after Extraction.Reposition.
	IndexInRepositionedColumnFields int
	// Skip field. Set using ColumnFields.Skip.
	Skip bool
}
