package database

import (
	"fmt"
	"reflect"

	"github.com/brunoga/deep"
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/iter"
)

/*
Get Extraction database fields.
*/
func (n *GetColumnFields) Get(metadataModel any) (*ColumnFields, error) {
	const FunctionName = "SetColumnFields"

	if !n.isTableCollectionUIDValid && (!n.isJoinDepthValid || !n.isTableCollectionNameValid) {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("tableCollectionUID or joinDepth and tableCollectionName is required").WithData(gojsoncore.JsonObject{"MetadataModel": metadataModel}).WithNestedError(ErrDatabaseGetColumnFieldsError)
	}

	n.columnFields = NewColumnFields()

	var forEachError error
	iter.ForEach(metadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
		if (n.skip.IsValid() && n.skip.Match(fieldGroup)) || (n.add.IsValid() && !n.add.Match(fieldGroup)) {
			return false, false
		}

		if n.isTableCollectionUIDValid {
			if tableCollectionUID, ok := fieldGroup[core.DatabaseTableCollectionUid].(string); !ok || tableCollectionUID != *n.tableCollectionUID {
				return false, true
			}
		} else if n.isJoinDepthValid && n.isTableCollectionNameValid {
			joinDepth := int64(0)
			if value, ok := fieldGroup[core.DatabaseJoinDepth]; ok {
				if err := n.defaultConverter.Convert(value, &schema.DynamicSchemaNode{Type: reflect.TypeOf(int64(0)), Kind: reflect.Int64}, &joinDepth); err != nil {
					forEachError = NewError().WithFunctionName(FunctionName).WithMessage("convert joinDepth to int64 failed").WithNestedError(err)
					return true, true
				}
			} else {
				return false, true
			}

			tableCollectionName, err := gojsoncore.As[string](fieldGroup[core.DatabaseTableCollectionName])
			if err != nil {
				return false, true
			}

			if *n.joinDepth != joinDepth || *n.tableCollectionName != tableCollectionName {
				return false, true
			}
		} else {
			return false, false
		}

		if core.IsFieldAGroup(fieldGroup) {
			return false, false
		}

		if fieldColumName, ok := fieldGroup[core.DatabaseFieldColumnName].(string); ok && len(fieldColumName) > 0 {
			if _, ok := n.columnFields.Fields[fieldColumName]; ok {
				forEachError = NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("duplicate fieldColumnName '%s' found", fieldColumName)).WithNestedError(ErrDatabaseGetColumnFieldsError)
				return true, true
			}

			newField := fieldGroup
			if value, err := deep.Copy(fieldGroup); err == nil {
				newField = value
			}
			n.columnFields.ColumnFieldsReadOrder = append(n.columnFields.ColumnFieldsReadOrder, fieldColumName)
			n.columnFields.Fields[fieldColumName] = newField
			return false, false
		} else {
			forEachError = NewError().WithFunctionName(FunctionName).WithMessage("field column name not found in field group property").WithNestedError(ErrDatabaseGetColumnFieldsError)
			return true, true
		}
	})

	if forEachError != nil {
		return nil, forEachError
	}

	return n.columnFields, nil
}

func (n *GetColumnFields) WithAdd(value core.FieldGroupPropertiesMatch) *GetColumnFields {
	n.SetAdd(value)
	return n
}

func (n *GetColumnFields) SetAdd(value core.FieldGroupPropertiesMatch) {
	n.add = value
}

func (n *GetColumnFields) WithSkip(value core.FieldGroupPropertiesMatch) *GetColumnFields {
	n.SetSkip(value)
	return n
}

func (n *GetColumnFields) SetSkip(value core.FieldGroupPropertiesMatch) {
	n.skip = value
}

func (n *GetColumnFields) WithJoinDepth(value int64) *GetColumnFields {
	n.SetJoinDepth(value)
	return n
}

func (n *GetColumnFields) SetJoinDepth(value int64) {
	if n.joinDepth == nil {
		n.joinDepth = new(int64)
	}
	*n.joinDepth = value
	n.isJoinDepthValid = true
}

func (n *GetColumnFields) WithTableCollectionUID(value string) *GetColumnFields {
	n.SetTableCollectionUID(value)
	return n
}

func (n *GetColumnFields) SetTableCollectionUID(value string) {
	if n.tableCollectionUID == nil {
		n.tableCollectionUID = new(string)
	}
	*n.tableCollectionUID = value

	if len(*n.tableCollectionUID) > 0 {
		n.isTableCollectionUIDValid = true
	}
}

func (n *GetColumnFields) WithTableCollectionName(value string) *GetColumnFields {
	n.SetTableCollectionName(value)
	return n
}

func (n *GetColumnFields) SetTableCollectionName(value string) {
	if n.tableCollectionName == nil {
		n.tableCollectionName = new(string)
	}
	*n.tableCollectionName = value

	if len(*n.tableCollectionName) > 0 {
		n.isTableCollectionNameValid = true
	}
}

func (n *GetColumnFields) WithDefaultConverter(value schema.DefaultConverter) *GetColumnFields {
	n.SetDefaultConverter(value)
	return n
}

func (n *GetColumnFields) SetDefaultConverter(value schema.DefaultConverter) {
	n.defaultConverter = value
}

func NewGetColumnFields() *GetColumnFields {
	n := new(GetColumnFields)
	n.defaultConverter = schema.NewConversion()
	return n
}

/*
GetColumnFields Retrieve database field/column information from metadata model.

Usage:
 1. Instantiate using NewGetColumnFields.
 2. Set required parameters.
 3. Extract ColumnFields using GetColumnFields.Get.

Example:

	// Set metadata model
	var metadataModel gojsoncore.JsonObject

	gcf := NewGetColumnFields()

	// Set
	gcf.SetTableCollectionUID("_12xoP1y")
	// Or
	gcf.SetJoinDepth(1)
	gcf.SetTableCollectionName("User")

	columnFields, err := gcf.Get(testData.MetadataModel)
*/
type GetColumnFields struct {
	columnFields *ColumnFields

	tableCollectionUID        *string
	isTableCollectionUIDValid bool

	joinDepth                  *int64
	isJoinDepthValid           bool
	tableCollectionName        *string
	isTableCollectionNameValid bool

	// skip a field/group if its properties matches one of the entries values.
	skip core.FieldGroupPropertiesMatch

	// add process the field/group only if its properties matches one of the entries values.
	add core.FieldGroupPropertiesMatch

	defaultConverter schema.DefaultConverter
}

func NewColumnFields() *ColumnFields {
	n := new(ColumnFields)
	n.ColumnFieldsReadOrder = make(core.MetadataModelGroupReadOrderOfFields, 0)
	n.Fields = make(ColumnFieldsFields)
	return n
}

/*
ColumnFields represents columns/fields for a particular core.DatabaseTableCollectionName at a specific core.DatabaseJoinDepth.
*/
type ColumnFields struct {
	// ColumnFieldsReadOrder read order of columns/fields.
	ColumnFieldsReadOrder core.MetadataModelGroupReadOrderOfFields
	// Fields a map of  database columns/fields properties from metadata model.
	Fields ColumnFieldsFields
}

/*
ColumnFieldsFields a map of  database columns/fields properties from metadata model.

Key is the core.FieldGroupJsonPathKey suffix.
*/
type ColumnFieldsFields map[string]gojsoncore.JsonObject
