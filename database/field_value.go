package database

import (
	"fmt"
	"reflect"

	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
)

/*
Get retrieves value(s).

Parameters:
  - columnFieldName - Required. Will match core.DatabaseFieldColumnName field property.
  - suffixJsonPath - Additional path.JSONPath segment to append to core.FieldGroupJsonPathKey.
  - arrayIndexes - Replace core.ArrayPathPlaceholder in path.JSONPath from core.FieldGroupJsonPathKey with a specific set of array indexes.

Returns no of results found and last error.

If satisfied with no of results, call FieldValue.GetValueFoundInterface or FieldValue.GetValueFoundReflected to retrieve the value found.
*/
func (n *FieldValue) Get(columnFieldName string, suffixJsonPath path.JSONPath, arrayIndexes []int) (uint64, error) {
	const FunctionName = "Get"

	jsonPathKey, err := n.getJsonPathToValue(columnFieldName, suffixJsonPath, arrayIndexes)
	if err != nil {
		return 0, err
	}

	return n.sourceData.Get(jsonPathKey)
}

// GetValueFoundReflected returns the value found as a reflect.Value.
func (n *FieldValue) GetValueFoundReflected() reflect.Value {
	return n.sourceData.GetValueFoundReflected()
}

// GetValueFoundInterface returns the value found as an interface{}.
func (n *FieldValue) GetValueFoundInterface() any {
	return n.sourceData.GetValueFoundInterface()
}

/*
Set inserts or updates value(s).

Parameters:
  - columnFieldName - Required. Will match core.DatabaseFieldColumnName field property.
  - valueToSet - Value to be inserted.
  - suffixJsonPath - Additional path.JSONPath segment to append to core.FieldGroupJsonPathKey.
  - arrayIndexes - Replace core.ArrayPathPlaceholder in path.JSONPath from core.FieldGroupJsonPathKey with a specific set of array indexes.
*/
func (n *FieldValue) Set(columnFieldName string, valueToSet any, suffixJsonPath path.JSONPath, arrayIndexes []int) (uint64, error) {
	const FunctionName = "Set"

	jsonPathKey, err := n.getJsonPathToValue(columnFieldName, suffixJsonPath, arrayIndexes)
	if err != nil {
		return 0, err
	}

	if reflect.ValueOf(valueToSet).Kind() == reflect.Array || reflect.ValueOf(valueToSet).Kind() == reflect.Slice {
		return n.sourceData.Set(jsonPathKey, valueToSet)
	}
	return n.sourceData.Set(jsonPathKey, []any{valueToSet})
}

/*
Delete removes value(s).

Parameters:
  - columnFieldName - Required. Will match core.DatabaseFieldColumnName field property.
  - suffixJsonPath - Additional path.JSONPath segment to append to core.FieldGroupJsonPathKey. Should NOT begin with path.JsonpathDotNotation.
  - arrayIndexes - Replace core.ArrayPathPlaceholder in path.JSONPath from core.FieldGroupJsonPathKey with a specific set of array indexes.
*/
func (n *FieldValue) Delete(columnFieldName string, suffixJsonPath path.JSONPath, arrayIndexes []int) (uint64, error) {
	const FunctionName = "Delete"

	jsonPathKey, err := n.getJsonPathToValue(columnFieldName, suffixJsonPath, arrayIndexes)
	if err != nil {
		return 0, err
	}

	return n.sourceData.Delete(jsonPathKey)
}

func (n *FieldValue) getJsonPathToValue(columnFieldName string, suffixJsonPath path.JSONPath, arrayIndexes []int) (path.JSONPath, error) {
	const FunctionName = "getJsonPathToValue"

	if columnFieldName == "" {
		return "", NewError().WithFunctionName(FunctionName).WithMessage("column field name is empty").WithNestedError(ErrDatabaseFieldValueError)
	}

	if n.columnFields == nil {
		return "", NewError().WithFunctionName(FunctionName).WithMessage("column fields is nil").WithNestedError(ErrDatabaseFieldValueError)
	}

	columnField, ok := n.columnFields.Fields[columnFieldName]
	if !ok {
		return "", NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Field with %s '%s' not found", core.DatabaseFieldColumnName, columnFieldName)).WithNestedError(ErrDatabaseFieldValueError)
	}

	jsonPathKey, err := core.AsJSONPath(columnField[core.FieldGroupJsonPathKey])
	if err != nil {
		return "", NewError().WithFunctionName(FunctionName).WithMessage("AsJSONPath failed").WithNestedError(err)
	}

	if len(suffixJsonPath) > 0 {

	}

	jsonPathToValue, err := core.NewJsonPathToValue().WithSourceOfValueIsAnArray(n.objectSourceIsAnArray).WithRemoveGroupFields(true).Get(jsonPathKey+path.JSONPath(path.JsonpathDotNotation)+suffixJsonPath, arrayIndexes)
	if err != nil {
		return "", NewError().WithFunctionName(FunctionName).WithMessage("Get JsonPathToValue failed").WithNestedError(err)
	}

	return jsonPathToValue, nil
}

// WithColumnFields sets the column fields.
func (n *FieldValue) WithColumnFields(value *ColumnFields) *FieldValue {
	n.SetColumnFields(value)
	return n
}

// SetColumnFields sets the column fields.
func (n *FieldValue) SetColumnFields(value *ColumnFields) {
	n.columnFields = value
}

// WithSourceData sets the source data object.
func (n *FieldValue) WithSourceData(value *object.Object) *FieldValue {
	n.SetSourceData(value)
	return n
}

// SetSourceData sets the source data object.
func (n *FieldValue) SetSourceData(value *object.Object) {
	n.sourceData = value
	objectSourceReflect := n.sourceData.GetSourceReflected()
	if objectSourceReflect.Kind() == reflect.Slice || objectSourceReflect.Kind() == reflect.Array {
		n.objectSourceIsAnArray = true
	} else {
		n.objectSourceIsAnArray = false
	}
}

/*
NewFieldValue

Parameters:

  - sourceData - Refer to object.Object.

    Contains the actual source value to manipulated and does the actual manipulation of data.

  - columnFields - Obtain using GetColumnFields.

    Provides the path.JSONPath information needed to manipulate data using methods from object.Object.
*/
func NewFieldValue(sourceData *object.Object, columnFields *ColumnFields) *FieldValue {
	n := new(FieldValue)
	n.SetSourceData(sourceData)
	n.SetColumnFields(columnFields)
	return n
}

/*
FieldValue Get, Set, and Delete field values in FieldValue.sourceData.source using the field properties: core.DatabaseTableCollectionUid or core.DatabaseJoinDepth with core.DatabaseTableCollectionName, and core.DatabaseFieldColumnName.

Usage:
 1. Instantiate FieldValue using NewFieldValue.
 2. Set required parameters.
 3. Manipulate FieldValue.sourceData.source using FieldValue.Get, FieldValue.Set, or FieldValue.Delete.
 4. Get modified source using FieldValue.sourceData.GetSource.
*/
type FieldValue struct {
	// Use to get, set, and delete values in source. Refer to sourceData.
	sourceData *object.Object

	// To be set when FieldValue.SetSourceData is called.
	objectSourceIsAnArray bool

	// Extracted database fields. Refer to GetColumnFields.
	columnFields *ColumnFields
}
