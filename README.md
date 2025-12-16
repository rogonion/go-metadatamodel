# go-metadatamodel

## Sections

- [Installation](#installation)
- [Modules](#modules)
    - [Field Columns](#field-columns)
    - [Filter](#filter)
    - [Database](#database)
    - [Iteration](#iteration)

## Installation

```shell
go get github.com/rogonion/go-metadatamodel
```

## Modules

### Field Columns

This [module](fieldcolumns) can be used to extract fields from a metadata model into a structure that resembles columns in a table.

It can:
- Extract field properties into an ordered slice of fields, resembling columns in a table -> ColumnFields.
- Set the new read order of column fields after repositioning -> ColumnFields.Reposition.
- Set column fields to skip based on core.FieldGroupPropertiesMatch -> ColumnFields.Skip.

Example usage:

```go
package main

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/fieldcolumns"
)

// Set metadata model
var metadataModel gojsoncore.JsonObject

// Module for extracting fields
var fcExtraction *fieldcolumns.Extraction = fieldcolumns.NewColumnFieldsExtraction(metadataModel)

var columnFields *fieldcolumns.ColumnFields
var err error
columnFields, err = fcExtraction.Extract()

// Check err and that columnFields is not nil

// Using ColumnFields.RepositionFieldColumns, reorder ColumnFields.CurrentIndexOfReadOrderOfColumnFields
columnFields.Reposition()

// update ColumnField.IndexInRepositionedColumnFields of each ColumnFields.Fields
columnFields.UpdateIndexInRepositionedColumnFieldsInColumnField()

// if field property does not match, skip it
var add core.FieldGroupPropertiesMatch

// if a field property matches, skip
var skip core.FieldGroupPropertiesMatch

// update ColumnField.Skip of each ColumnFields.Fields
columnFields.Skip(skip, add)

```

### Filter

This [module](filter) can be used to filter through data with a metadata model structure.

Designed to support both simple queries and deeply nested logical operator queries which are extensible and customizable.

Below is a sample query condition structure:

```json
{
  /* The current query condition context: Can be 'LogicalOperator' for nesting or 'FieldGroup' for the actual filter condition */
  "Type": "LogicalOperator",
  "Negate": false,
  "LogicalOperator": "And",
  // Can be 'And' or 'Or'. Default 'And'.
  "Value": [
    {
      "Type": "LogicalOperator",
      "LogicalOperator": "Or",
      "Value": [
        {
          "Type": "FieldGroup",
          "Negate": false,
          "LogicalOperator": "And",
          "Value": {
            "$.GroupFields[*].Bio": {
              "EqualTo": {
                "AssumedFieldType": "Any",
                "Values": [
                  true,
                  "Yes"
                ]
              }
            }
          }
        },
        {
          "Type": "FieldGroup",
          "Negate": false,
          "LogicalOperator": "And",
          "Value": {
            "$.GroupFields[*].Bio": {
              "EqualTo": {
                "AssumedFieldType": "Text",
                "Negate": true,
                "Value": "no"
              }
            },
            "$.GroupFields[*].Occ": {
              "EqualTo": {
                "AssumedFieldType": "Text",
                "Negate": true,
                "Value": "no"
              }
            }
          }
        }
      ]
    },
    {
      "Type": "FieldGroup",
      "Value": {
        "$.GroupFields[*].SiteAndGeoreferencing.GroupFields[*].Country": {
          "FullTextSearchQuery": {
            "AssumedFieldType": "Text",
            "Value": "Kenya",
            "ExactMatch": true
          }
        }
      }
    },
    {
      "Type": "FieldGroup",
      "Negate": false,
      "LogicalOperator": "And",
      "Value": {
        "$.GroupFields[*].SiteAndGeoreferencing.GroupFields[*].Sites.GroupFields[*].Coordinates.GroupFields[*].Latitude": {
          /*
          Default processed as 'And' logical operation.          
          The key should be a unique filter condition in the system while the value is an object with the relevant filter condition value.
          */
          "GreaterThan": {
            "AssumedFieldType": "Number",
            "Value": 20.00
          },
          "LessThan": {
            "AssumedFieldType": "Number",
            "Value": 21.00
          }
        }
      }
    }
  ]
}
```

Example usage:

```go
package main

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-metadatamodel/filter"
)

// Set metadata model
var metadataModel gojsoncore.JsonObject

// Set source data
var sourceData *object.Object

// Set query condition
var queryCondition gojsoncore.JsonObject

// Set other properties using builder pattern 'With' or 'Set'. Refer to filter.FilterData structure.
var filterData *filter.DataFilter = filter.NewFilterData(sourceData, metadataModel)

var filterExcludeIndexes []int
var err error

filterExcludeIndexes, err = filterData.Filter(queryCondition, "", "")

```

### Database

This [module](database) can be used to work with data (get, set, delete) whose metadata model represents a relational
database structure
using the following field/group properties:

- core.DatabaseTableCollectionUid
- core.DatabaseJoinDepth
- core.DatabaseTableCollectionName
- core.DatabaseFieldColumnName

Example usage:

#### Get Column Fields

Retrieve column fields information from a metadata model.

```go
package main

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/database"
)

// Set metadata model
var metadataModel gojsoncore.JsonObject

var gcf *database.GetColumnFields = database.NewGetColumnFields()

// Set
gcf.SetTableCollectionUID("_12xoP1y")
// Or
gcf.SetJoinDepth(1)
gcf.SetTableCollectionName("User")

var columnFields *database.ColumnFields
var err error
columnFields, err = gcf.Get(metadataModel)

```

#### Field Value

A set of methods to get, set, and delete value(s) in a source object using database properties in the metadata model.

Module uses [go-json](https://github.com/rogonion/go-json) for actual data manipulation.

```go
package main

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
var sch schema.Schema

// Source object
var product *Product = &Product{
	ID: []int{1},
}

// source data to manipulate
var obj *object.Object = object.NewObject().WithSourceInterface(product).WithSchema(sch)

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

```

### Iteration

This [module](iter) provides higher-order functions processing the fields in a metadata model.

Provides the following methods:

- Filter
- For Each
- Map

Example usage:

#### Filter

```go
package main

import (
	"strings"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/iter"
	"github.com/rogonion/go-metadatamodel/testdata"
)

var sourceMetadataModel any = testdata.AddressMetadataModel(nil)

var updatedMetadataModel any = iter.Filter(sourceMetadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
	if fieldGroupName, ok := fieldGroup[core.FieldGroupName].(string); ok {
		if strings.HasSuffix(fieldGroupName, "Name") {
			return false, false
		}
	}
	return true, false
})

```

#### Map

```go
package main

import (
	"strings"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/iter"
	"github.com/rogonion/go-metadatamodel/testdata"
)

var sourceMetadataModel any = testdata.AddressMetadataModel(nil)

var updatedMetadataModel any = iter.Map(sourceMetadataModel, func(fieldGroup gojsoncore.JsonObject) (any, bool) {
	if fieldGroupName, ok := fieldGroup[core.FieldGroupName].(string); ok {
		if strings.HasSuffix(fieldGroupName, "Code") {
			fieldGroup[core.FieldGroupName] = fieldGroupName + " Found"
		}
	}
	return fieldGroup, false
})

```

#### For Each

```go
package main

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/iter"
	"github.com/rogonion/go-metadatamodel/testdata"
)

var sourceMetadataModel any = testdata.AddressMetadataModel(nil)

var noOfIterations uint64 = 0

iter.ForEach(sourceMetadataModel, func (fieldGroup gojsoncore.JsonObject)(bool, bool) {
	noOfIterations++
	return false, false
})

```