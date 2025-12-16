/*
Package fieldcolumns can be used to extract fields in a metadata model into a structure that resembles columns in a table.

It can perform the following tasks:
  - Extracts field properties into an ordered slice of fields, resembling columns in a table -> ColumnFields.
  - Set the new read order of column fields after repositioning -> ColumnFields.Reposition.
  - Set column fields to skip based on core.FieldGroupPropertiesMatch -> ColumnFields.Skip.

# Usage

## Extraction

Module can be used to recursively extract fields in a metadata model into ColumnFields.

1. Create a new instance of the Extraction struct using NewColumnFieldsExtraction.

The following parameters can be set using the builder method (prefixed `With`) or Set (prefixed `Set):
  - schema - Set using Extraction.WithSchema or Extraction.SetSchema.
  - skip - Set using Extraction.WithSkip or Extraction.SetSkip.
  - add - Set using Extraction.WithAdd or Extraction.SetAdd.

2. Begin field data extraction using Extraction.Extract.

Example:

	// Set metadata model
	var metadataModel core.JsonObject

	fcExtraction := NewColumnFieldsExtraction(metadataModel)
	var columnFields *ColumnFields
	var err error
	columnFields, err = fcExtraction.Extract()

	// Check err and that columnFields is not nil

## ColumnFields.Reposition

After extracting metadata model fields into ColumnFields, you can reposition them based on the ColumnFields.RepositionFieldColumns information set during Extraction.Extract.

1. Call ColumnFields.Reposition to update the ColumnFields.CurrentIndexOfReadOrderOfColumnFields.

	var columnFields *ColumnFields
	columnFields.Reposition()

2. Optionally, update ColumnField.IndexInRepositionedColumnFields of each ColumnFields.Fields.

	var columnFields *ColumnFields
	columnFields.UpdateIndexInRepositionedColumnFieldsInColumnField()

## ColumnFields.Skip

Set the ColumnField.Skip of each ColumnFields.Fields.

Useful for automated skips of processing fields if they match core.FieldGroupPropertiesMatch.

Example:

	var columnFields *ColumnFields

	// if field property does not match, skip it
	var add core.FieldGroupPropertiesMatch

	// if a field property match, skip
	var skip core.FieldGroupPropertiesMatch

	columnFields.Skip(skip, add)
*/
package fieldcolumns
