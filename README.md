# go-metadatamodel

`go-metadatamodel` is a Go library designed to manipulate, transform, and query complex data structures using a declarative Metadata Model. It provides tools for flattening/unflattening nested data, filtering, extracting column definitions, and handling database-like operations on in-memory objects.

## Sections

- Prerequisites
- Installation
- Environment Setup
- Modules
    - Database
    - Field Columns
    - Filter
    - Flattener
    - Iteration
    - Unflattener

## Prerequisites

Ensure you have the following installed on your system. This project supports Linux, macOS, and Windows (via WSL2).

<table>
  <thead>
    <tr>
      <th>Tool</th>
      <th>Description</th>
      <th>Link</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Go</td>
      <td>Programming language used for the project.</td>
      <td><a href="https://go.dev/">Official Website</a></td>
    </tr>
    <tr>
      <td>Task</td>
      <td>
        <p>Task runner / build tool.</p> 
        <p>You can use the provided shell script <a href="taskw">taskw</a> that automatically downloads the binary and install it in the <code>.task</code> folder.</p>
      </td>
      <td><a href="https://taskfile.dev/">Official Website</a></td>
    </tr>
    <tr>
      <td>Docker / Podman</td>
      <td>Optional container engine for isolated development environment.</td>
      <td><a href="https://www.docker.com/">Docker</a> / <a href="https://podman.io/">Podman</a></td>
    </tr>
  </tbody>
</table>

After building the dev container, below is a sample script that runs the container and mounts the project directory into the container:

```shell
#!/bin/bash

CONTAINER_ENGINE="podman"
CONTAINER="projects-go-metadatamodel"
NETWORK="systemd-leap"
NETWORK_ALIAS="projects-go-metadatamodel"
CONTAINER_UID=1000
IMAGE="localhost/projects/go-metadatamodel:latest"
SSH_PORT="127.0.0.1:2200" # for local proxy vscode ssh access
PROJECT_DIRECTORY="$(pwd)"

# Check if container exists (Running or Stopped)
if $CONTAINER_ENGINE ps -a --format '{{.Names}}' | grep -q "^$CONTAINER$"; then
    echo "   Found existing container: $CONTAINER"
    # Check if it is currently running
    if $CONTAINER_ENGINE ps --format '{{.Names}}' | grep -q "^$CONTAINER$"; then
        echo "✅ Container is already running."
    else
        echo "🔄 Container stopped. Starting it..."
        $CONTAINER_ENGINE start $CONTAINER
        echo "✅ Started."
    fi
else
    # Container doesn't exist -> Create and Run it
    echo "🆕 Container not found. Creating new..."
    $CONTAINER_ENGINE run -d \
    # start container from scratch
    # `sudo` is used because systemd-leap network was created in `sudo`
    # Ensure container image exists in `sudo`
    # Not needed if target network is not in `sudo`
    sudo podman run -d \
        --name $CONTAINER \
        --network $NETWORK \
        --network-alias $NETWORK_ALIAS \
        --user $CONTAINER_UID:$CONTAINER_UID \
        -p $SSH_PORT:22 \
        -v $PROJECT_DIRECTORY:/home/dev/go-metadatamodel:Z \
        $IMAGE
    echo "✅ Created and Started."
fi
```

## Installation

```shell
go get github.com/rogonion/go-metadatamodel
```

## Environment Setup

This project uses `Taskfile` to manage the development environment and tasks.

<table>
  <thead>
    <tr>
      <th>Task</th>
      <th>Description</th>
      <th>Usage</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>env:build</code></td>
      <td>
        <p>Build the dev container image.</p>
        <p>Image runs an ssh server one can connect to with vscode.</p>
      </td>
      <td><code>task env:build</code></td>
    </tr>
    <tr>
      <td><code>env:info</code></td>
      <td>Show current environment configuration.</td>
      <td><code>task env:info</code></td>
    </tr>
    <tr>
      <td><code>deps</code></td>
      <td>Download and tidy dependencies.</td>
      <td><code>task deps</code></td>
    </tr>
    <tr>
      <td><code>test</code></td>
      <td>Run tests. Supports optional <code>TARGET</code> variable.</td>
      <td><code>task test</code><br><code>task test TARGET=./database</code></td>
    </tr>
  </tbody>
</table>

## Modules

### Database

This module can be used to work with data (get, set, delete) whose metadata model represents a relational
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

Module uses go-json for actual data manipulation.

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
columnFields, err = database.NewGetColumnFields().WithJoinDepth(0).WithTableCollectionName("Product").Get(productMetadataModel)

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

// Set other properties using builder pattern 'With' or 'Set'. Refer to filter.DataFilter structure.
var filterData *filter.DataFilter = filter.NewFilterData(sourceData, metadataModel)

var filterExcludeIndexes []int
var err error

filterExcludeIndexes, err = filterData.Filter(queryCondition, "", "")

```

### Flattener

This module converts deeply nested data structures into flat 2D tables based on a Metadata Model.

It acts as a "Projection" engine, capable of:
- Recursively flattening nested objects/arrays.
- Generating Cartesian products for one-to-many relationships.
- Pivoting specific fields into horizontal columns.

Example usage:

```go
package main

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-metadatamodel/flattener"
)

// ... setup metadataModel ...

// 1. Initialize
f := flattener.NewFlattener(metadataModel)

// 2. Flatten Source
sourceObj := object.NewObject().WithSourceInterface(myData)
err := f.Flatten(sourceObj)

// 3. Get Results (Raw Table)
table := f.GetResult()

// 4. Or Write to Destination (Object)
destObj := object.NewObject().WithSourceInterface(make([][]any, 0))
err = f.WriteToDestination(destObj)
```

### Iteration

This module provides higher-order functions processing the fields in a metadata model.

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

### Unflattener

This module converts a 2D array/slice (FlattenedTable) back into a slice of complex objects.

It is the inverse of the Flattener, reconstructing hierarchies using Primary Keys defined in the Metadata Model.

Example usage:

```go
package main

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-metadatamodel/unflattener"
)

// ... setup metadataModel ...

// 1. Initialize with Signature generator
sig := unflattener.NewSignature()
u := unflattener.NewUnflattener(metadataModel, sig)

// 2. Prepare Destination
var dest []*MyStruct
destObj := object.NewObject().WithSourceInterface(&dest)
u.WithDestination(destObj)

// 3. Unflatten
// sourceTable is [][]reflect.Value (from flattener)
err := u.Unflatten(sourceTable)
```
