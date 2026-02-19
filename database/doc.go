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

You can identify the target Table Collection in two ways:
  - **By Unique ID**: Use `WithTableCollectionUID`. This is the most precise method if your metadata model assigns unique IDs to collections.
  - **By Name & Join Depth**: Use `WithTableCollectionName` AND `WithJoinDepth`. This is useful when IDs are not available or when traversing a join hierarchy where the same table name might appear at different depths.

Example:

	import (
		gojsoncore "github.com/rogonion/go-json/core"
		"github.com/rogonion/go-metadatamodel/database"
	)

	// Set metadata model
	var metadataModel gojsoncore.JsonObject

	gcf := database.NewGetColumnFields()

	// Option A: Set by UID
	gcf.SetTableCollectionUID("_12xoP1y")
	// Option B: Set by Name and Depth
	// gcf.SetJoinDepth(1)
	// gcf.SetTableCollectionName("User")

	columnFields, err := gcf.Get(metadataModel)

## FieldValue

Module can be used to FieldValue.Get, FieldValue.Set, and FieldValue.Delete value(s) in an sourceData using metadata model and its database properties.

Example:

	import (
		gojsoncore "github.com/rogonion/go-json/core"
		"github.com/rogonion/go-json/object"
		"github.com/rogonion/go-json/schema"
		"github.com/rogonion/go-metadatamodel/database"
	)

	type Product struct {
		ID    []int
		Name  []string
		Price []float64
	}

	// Set Product schema. Useful for instantiating nested collections
	var sch *schema.DynamicSchemaNode

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
	columnFields, err = database.NewGetColumnFields().WithJoinDepth(0).WithTableCollectionName("Product").Get(productMetadataModel)

	// Module to perform get,set, or delete
	var fieldValue *database.FieldValue = database.NewFieldValue(sourceData, columnFields)

	var noOfResults uint64

	// Get value of column `ID`
	// Returns number of results found (uint64) and error
	noOfResults, err = fieldValue.Get("ID", "", nil)

	var noOfModifications uint64

	// Set value for column `Name`
	noOfModifications, err = fieldValue.Set("Name", "Twinkies", "", nil)

	// Delete value for column `Price`
	noOfModifications, err = fieldValue.Delete("Price", "", nil)
*/
package database
