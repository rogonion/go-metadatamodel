/*
Package core contains code shared amongst packages within the module go-metadatamodel.

It contains type definitions, constants, as well as shared modules.

# Usage

## FieldGroupPropertiesMatch

Module can be used to define conditions for a set of properties in a field/group to skip/add when recursively going through a metadata model.

## JsonPathToValue

Module converts FieldGroupJsonPathKey to actual path.JSONPath in an object.

The result can then be used by object.Object to manipulate a source object.

It can do the following:
  - Remove GroupFields from path.JSONPath.
  - Replace ArrayPathPlaceholder with actual array index.

To begin using the module:

1. Create a new instance of the JsonPathToValue using the method NewJsonPathToValue.

The following parameters can be set using the builder method (prefixed `With`) or Set (prefixed `Set):
  - removeGroupFields - remove the snippet GroupFields from the path.JSONPath.
  - sourceOfValueIsAnArray - If source of value is NOT an array or slice (false), the first pair of GroupJsonPathPrefix is removed as the source is assumed to be an associative collection.

2. Retrieve actual path.JSONPath by calling JsonPathToValue.Get.

Example:

	jsonPathKey := "$.GroupFields[*].SiteAndGeoreferencing.GroupFields[*].Country"
	arrayIndexes := []int{1,0}

	jptv := core.NewJsonPathToValue().WithRemoveGroupFields(true).WithSourceOfValueIsAnArray(true)

	jsonPathToValue, err := jptv.Get(jsonPathKey, arrayIndexes)
	if err != nil {
		return "", NewError(FunctionName, "Get JsonPathToValue failed").WithNestedError(err)
	}

## FieldValue

Example:

	type Product struct {
		ID    []int
		Name  []string
		Price []float64
	}

	// Set Product schema. Useful for instantiating nested collections
	var sch schema.Schema

	// Source object
	var product *Product = &Product{
		ID: []int{1},
	}

	// source data to manipulate
	var obj *object.Object = object.NewObject(product).WithSchema(sch)

	// Set product metadata model
	var productMetadataModel gojsoncore.JsonObject

	var columnFields *database.ColumnFields
	var err error
	columnFields, err = NewGetColumnFields().WithJoinDepth(0).WithTableCollectionName("Product").Get(productMetadataModel)

	// Module to perform get,set, or delete
	var fieldValue *database.FieldValue = database.NewFieldValue(obj, columnFields)

	var res any
	var ok bool

	// Get value of column `ID`
	res, ok, err = fieldValue.Get("ID", "", nil)

	var noOfModifications uint64

	// Set value for column `Name`
	noOfModifications, err = fieldValue.Set("Name", "Twinkies", "", nil)

	// Delete value for column `Price`
	noOfModifications, err = fieldValue.Delete("Price", "", nil)

## Utils

Shared utility functions for manipulating and inspecting metadata models.

### MergeRightJsonObjectIntoLeft

Merges two JsonObjects.

	left := gojsoncore.JsonObject{"a": 1}
	right := gojsoncore.JsonObject{"b": 2}
	core.MergeRightJsonObjectIntoLeft(left, right)
	// left is now {"a": 1, "b": 2}

### IsFieldAGroup

	isGroup := core.IsFieldAGroup(fieldGroup)
*/
package core
