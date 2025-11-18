/*
Package database can be used to work with data whose metadata model represents a relational database structure.

The following Field/Group properties are important:
  - core.DatabaseTableCollectionUid
  - core.DatabaseJoinDepth
  - core.DatabaseTableCollectionName
  - core.DatabaseFieldColumnName

# Usage

## GetColumnFields

Module can be used to extract database field(s) information into ColumnFields.

To begin using the module:

1. Create a new instance of the GetColumnFields struct. You can use the convenience method NewGetColumnFields which:
  - Sets GetColumnFields.defaultConverter to schema.NewConversion

The following parameters can be set using the builder method (prefixed `With`) or Set (prefixed `Set):
  - tableCollectionUID - Set using GetColumnFields.WithTableCollectionUID or GetColumnFields.SetTableCollectionUID.
  - joinDepth - Set using GetColumnFields.WithJoinDepth or GetColumnFields.SetJoinDepth.
  - tableCollectionName - Set using GetColumnFields.WithTableCollectionName or GetColumnFields.SetTableCollectionName.
  - skip - Set using GetColumnFields.WithSkip or GetColumnFields.SetSkip.
  - add - Set using GetColumnFields.WithAdd or GetColumnFields.SetAdd.

2. Extract the database fields using GetColumnFields.Get.

Example:

	// Set metadata model
	var metadataModel core.JsonObject

	gcf := NewGetColumnFields()

	// Set
	gcf.SetTableCollectionUID("_12xoP1y")
	// Or
	gcf.SetJoinDepth(1)
	gcf.SetTableCollectionName("User")

	columnFields, err := gcf.Get(testData.MetadataModel)

## FieldValue

Module can be used to FieldValue.Get, FieldValue.Set, and FieldValue.Delete value(s) in an sourceData using metadata model and its database properties.

Example:

	import (
		gojsoncore "github.com/rogonion/go-json/core"
	)

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
	var sourceData *object.Object = object.NewObject(product).WithSchema(sch)

	// Set product metadata model
	var productMetadataModel gojsoncore.JsonObject

	var columnFields *database.ColumnFields
	var err error
	columnFields, err = NewGetColumnFields().WithJoinDepth(0).WithTableCollectionName("Product").Get(productMetadataModel)

	// Module to perform get,set, or delete
	var fieldValue *database.FieldValue = database.NewFieldValue(sourceData, columnFields)

	var res any
	var ok bool

	// Get value of column `ID`
	res, ok, err = fieldValue.Get("ID", "", nil)

	var noOfModifications uint64

	// Set value for column `Name`
	noOfModifications, err = fieldValue.Set("Name", "Twinkies", "", nil)

	// Delete value for column `Price`
	noOfModifications, err = fieldValue.Delete("Price", "", nil)
*/
package database
