/*
Package flattener provides functionality to convert deeply nested data structures into flat 2D tables based on a Metadata Model.

It acts as a "Projection" engine, mapping complex hierarchical data (like JSON or nested structs) into a linear format suitable for CSV exports, data analysis, or grid views.

It can perform the following tasks:
  - Recursively flatten nested objects and arrays into a 2D `FlattenedTable`.
  - Handle one-to-many relationships by generating Cartesian products (row explosion).
  - Handle specific fields by pivoting them into horizontal columns (horizontal expansion).
  - Write the flattened results into a destination `object.Object` using schema-based type conversion.
  - Support batch processing via the `Reset` method.

# Usage

## Initialization

1. Create a new instance of the Flattener using `NewFlattener`, providing the Metadata Model that defines the structure.

	// Set metadata model
	var metadataModel gojsoncore.JsonObject
	flattener := NewFlattener(metadataModel)

2. Optionally, configure column behavior (skipping/reordering) using `WithColumnFields`.

	var columnFields *fieldcolumns.ColumnFields
	// ... initialize columnFields ...
	flattener.WithColumnFields(columnFields)

## Flattening Data

Use `Flatten` to process a source object. This method appends the results to the Flattener's internal state, allowing for batch accumulation.

	var sourceData any // map[string]any or struct
	sourceObj := object.NewObject().WithSourceInterface(sourceData)

	err := flattener.Flatten(sourceObj)

## Retrieving Results

There are two ways to retrieve the flattened data:

1. **Direct Access:** Use `GetResult` to get the raw `FlattenedTable` (a `[][]reflect.Value`).

	table := flattener.GetResult()

2. **Write to Destination:** Use `WriteToDestination` to map the results into a target `object.Object`. This applies column reordering/skipping defined in `ColumnFields` and utilizes the destination's schema for type conversion.

	// Destination could be a 2D slice, a CSV writer wrapper, etc.
	destObj := object.NewObject().WithSourceInterface(make([][]any, 0))
	err := flattener.WriteToDestination(destObj)

## Batch Processing

To process large datasets in chunks, use the `Reset` method to clear the internal state without re-allocating the Flattener.

	for _, batch := range hugeDataset {
	    flattener.Flatten(object.NewObject().WithSourceInterface(batch))
	    flattener.WriteToDestination(finalOutput)
	    flattener.Reset() // Clear internal table for next batch
	}
*/
package flattener
