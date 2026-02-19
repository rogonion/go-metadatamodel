/*
Package unflattener provides functionality to convert a 2D array/slice (FlattenedTable) back into a slice of complex objects based on a Metadata Model.

It acts as the inverse of the `flattener` package. It reconstructs the hierarchical structure from flat data by using Primary Keys (defined in the Metadata Model) to group related rows together.

It can perform the following tasks:
  - Reconstruct nested objects and arrays from a flat table.
  - Handle one-to-many relationships by grouping rows that share the same parent key.
  - Handle pivoted columns (horizontal expansion) by mapping them back to their array representation.
  - Write the reconstructed objects into a destination `object.Object`.

# Usage

	import (
		gojsoncore "github.com/rogonion/go-json/core"
		"github.com/rogonion/go-json/object"
		"github.com/rogonion/go-metadatamodel/unflattener"
	)

## Initialization

1. Create a new instance of the Unflattener using `NewUnflattener`. You also need a `Signature` generator which handles key creation.

	var metadataModel gojsoncore.JsonObject // ... load metadata model
	signature := unflattener.NewSignature()
	unflattener := unflattener.NewUnflattener(metadataModel, signature)

2. Prepare the destination object. This should be a slice of pointers to your struct type.

	var dest []*MyStruct
	destObj := object.NewObject().WithSourceInterface(&dest)
	unflattener.WithDestination(destObj)

## Unflattening Data

1. Prepare your source data as a `flattener.FlattenedTable` (a `[][]reflect.Value`).

2. Call `Unflatten`.

	// sourceTable is [][]reflect.Value
	err := unflattener.Unflatten(sourceTable)

	// dest now contains the reconstructed object graph.
*/
package unflattener
