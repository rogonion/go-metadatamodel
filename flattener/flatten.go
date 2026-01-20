package flattener

import (
	"errors"
	"fmt"
	"reflect"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/fieldcolumns"
)

/*
WriteToDestination A general purpose function for writing current FlattenedTable to object.Object.
*/
func (n *Flattener) WriteToDestination(destination *object.Object) error {
	const FunctionName = "WriteToDestination"

	if len(n.currentSourceObjectResult) == 0 {
		return nil
	}

	// 1. Determine the Read Order
	// If columnFields is nil (raw dump mode), we create a linear 1:1 map.
	// If columnFields exists, we get the calculated read order (Repositioned & Skipped).
	var readOrder []int
	if n.columnFields != nil {
		readOrder = n.columnFields.GetCurrentIndexOfReadOrderOfFields()
	} else {
		// Fallback: Create a sequential list [0, 1, 2, ... N]
		// This ensures raw dumps work even without metadata manipulation.
		if len(n.currentSourceObjectResult) > 0 {
			cols := len(n.currentSourceObjectResult[0])
			readOrder = make([]int, cols)
			for i := 0; i < cols; i++ {
				readOrder[i] = i
			}
		}
	}

	// 2. Iterate and Write
	for rowIndex, row := range n.currentSourceObjectResult {
		for destColIndex, sourceColIndex := range readOrder {

			// Safety Check: Bounds validation
			if sourceColIndex < 0 || sourceColIndex >= len(row) {
				return NewError().WithFunctionName(FunctionName).
					WithMessage(fmt.Sprintf("source column index %d out of bounds", sourceColIndex)).
					WithData(gojsoncore.JsonObject{"RowIndex": rowIndex, "RowLength": len(row)})
			}

			cellValue := row[sourceColIndex]

			// Write to Destination using the Destination Index
			targetPath := path.JSONPath(fmt.Sprintf("$[%d][%d]", rowIndex, destColIndex))

			if _, err := destination.SetReflect(targetPath, cellValue); err != nil {
				return NewError().WithFunctionName(FunctionName).
					WithMessage(fmt.Sprintf("write failed at row %d col %d", rowIndex, destColIndex)).
					WithNestedError(err)
			}
		}
	}

	return nil
}

/*
recursiveConvert now takes the current table state and returns the mutated table.
*/
func (n *Flattener) recursiveConvert(groupConversion *FieldGroupConversion, linearCollectionIndexes []int, incomingRows FlattenedTable) (FlattenedTable, error) {
	const FunctionName = "recursiveConvert"

	// Working set starts as a copy of incoming rows
	currentRows := n.copyTable(incomingRows)

	// If incoming rows is empty (start of process), initialize with one empty row
	// so Cartesian products work.
	if len(currentRows) == 0 {
		currentRows = FlattenedTable{{}}
	}

	for _, fgConversion := range groupConversion.GroupFields {
		fieldData := DefaultEmptyColumn()
		if jsonPath, err := core.NewJsonPathToValue().Get(fgConversion.FieldColumnPosition.JSONPath(), linearCollectionIndexes); err != nil {
			return nil, NewError().WithFunctionName(FunctionName).WithMessage("get path to value at field failed").WithNestedError(err)
		} else {
			noOfResults, err := n.currentSourceObject.Get(jsonPath)
			if err != nil {
				return nil, NewError().WithFunctionName(FunctionName).WithMessage("get value at jsonPath failed").WithNestedError(err)
			}
			if noOfResults > 0 {
				fieldData = n.currentSourceObject.GetValueFoundReflected()
			}
		}

		// 1. Handle Nested Groups / Recursion
		if len(fgConversion.GroupFields) > 0 {
			if fieldData.Kind() == reflect.Slice || fieldData.Kind() == reflect.Array {
				// BRANCHING LOGIC (The Deep Copy replacement)
				// We are splitting the universe here: Row 1 vs Row 2 vs Row 3
				collectedBranchResults := make(FlattenedTable, 0)

				for i := 0; i < fieldData.Len(); i++ {
					// IMPORTANT: Pass a CLEAN COPY of currentRows to this branch
					// This ensures iteration i=0 doesn't mess up i=1
					branchInput := n.copyTable(currentRows)

					branchResult, err := n.recursiveConvert(fgConversion, append(linearCollectionIndexes, i), branchInput)
					if err != nil {
						return nil, err
					}
					collectedBranchResults = append(collectedBranchResults, branchResult...)
				}
				// The result of processing this field is the union of all branches
				currentRows = collectedBranchResults
				continue
			}

			// Non-array nested group (Single Object)
			var err error
			currentRows, err = n.recursiveConvert(fgConversion, append(linearCollectionIndexes, 0), currentRows)
			if err != nil {
				return nil, err
			}
			continue
		}

		// 2. Handle Leaf Fields (Merge into current rows)
		// This is a Sibling field, so we merge it into EVERY row in currentRows (Cartesian)
		cellValue := n.getCellValueFromFieldData(fieldData)
		currentRows = n.mergeCellValueIntoRows(currentRows, cellValue)
	}

	return currentRows, nil
}

/*
Helper to copy the table structure (to replace old deep.Copy logic)
*/
func (n *Flattener) copyTable(source FlattenedTable) FlattenedTable {
	if len(source) == 0 {
		return make(FlattenedTable, 0)
	}
	newTable := make(FlattenedTable, len(source))
	for i, row := range source {
		// Copy the row slice so appending to it doesn't affect the original
		newRow := make(FlattenedRow, len(row))
		copy(newRow, row)
		newTable[i] = newRow
	}
	return newTable
}

/*
getCellValueFromFieldData enforces the rule that every cell value is a slice/array.
*/
func (n *Flattener) getCellValueFromFieldData(fieldData reflect.Value) reflect.Value {
	if !fieldData.IsValid() {
		return DefaultEmptyColumn()
	}

	if fieldData.Kind() == reflect.Slice || fieldData.Kind() == reflect.Array {
		return fieldData
	}

	newFieldData := reflect.MakeSlice(reflect.SliceOf(fieldData.Type()), 1, 1)
	newFieldData.Index(0).Set(fieldData)
	return newFieldData
}

/*
Merges a single cell value (which is a slice) into all existing rows.

If the cell value is meant to explode rows (e.g. if the cell itself contains multiple items
that should be separate rows), logic differs, but based on your "normalizeToSlice",
we treat the cell as 1 unit of data for that column.
*/
func (n *Flattener) mergeCellValueIntoRows(rows FlattenedTable, cellValue reflect.Value) FlattenedTable {
	newRows := make(FlattenedTable, len(rows))

	for i, row := range rows {
		// Create new row with existing data + new cell
		// Note: row is already copied if coming from copyTable, but safer to append
		newRow := make(FlattenedRow, len(row)+1)
		copy(newRow, row)
		newRow[len(row)] = cellValue
		newRows[i] = newRow
	}

	return newRows
}

/*
Flatten processes the sourceObject and populates the internal Flattener.currentSourceObjectResult.

Once the process is successful, you can call Flattener.GetResult to retrieve the FlattenedTable.
*/
func (n *Flattener) Flatten(sourceObject *object.Object) error {
	const FunctionName = "Flatten"

	n.currentSourceObjectIsALinearCollection = sourceObject.GetSourceReflected().Kind() == reflect.Slice || sourceObject.GetSourceReflected().Kind() == reflect.Array

	if n.columnFields == nil {
		if columnFields, err := fieldcolumns.NewColumnFieldsExtraction(n.metadataModel).Extract(); err != nil {
			return NewError().WithFunctionName(FunctionName).WithMessage("extract default columnFields failed").WithNestedError(err)
		} else {
			columnFields.Reposition()
			n.columnFields = columnFields
		}
	}

	if n.fieldGroupConversion == nil {
		if value, err := n.recursiveInitFieldGroupConversion(n.metadataModel, path.JSONPath(path.JsonpathKeyRoot)); err != nil {
			return err
		} else {
			n.fieldGroupConversion = value
		}
	}

	var err error
	if n.currentSourceObjectIsALinearCollection {
		sourceObject.ForEach(path.JSONPath(path.JsonpathKeyRoot+path.JsonpathDotNotation+path.JsonpathLeftBracket+path.JsonpathKeyIndexAll+path.JsonpathRightBracket), func(jsonPath path.RecursiveDescentSegment, value reflect.Value) bool {
			n.currentSourceObject = object.NewObject().WithSourceReflected(value)
			var resultTable FlattenedTable
			resultTable, err = n.recursiveConvert(n.fieldGroupConversion, make([]int, 0), make(FlattenedTable, 0))
			if err == nil {
				n.currentSourceObjectResult = append(n.currentSourceObjectResult, resultTable...)
			}
			return false
		})
	} else {
		n.currentSourceObject = sourceObject
		var resultTable FlattenedTable
		resultTable, err = n.recursiveConvert(n.fieldGroupConversion, make([]int, 0), make(FlattenedTable, 0))
		if err == nil {
			n.currentSourceObjectResult = append(n.currentSourceObjectResult, resultTable...)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

/*
GetResult returns the raw FlattenedTable.
This allows the user to read the results directly.
*/
func (n *Flattener) GetResult() FlattenedTable {
	return n.currentSourceObjectResult
}

/*
Reset clears the internal state, allowing the Flattener to be reused for a new batch of data without carrying over previous results.
*/
func (n *Flattener) Reset() {
	// Re-initialize with capacity 0 but keeping the pointer if possible,
	// or just make a fresh slice. Fresh slice is safer for GC.
	n.currentSourceObjectResult = make(FlattenedTable, 0)
}

/*
Generates a tree of the source object with the necessary data required to perform the conversion at each recursive step.

Ignores if a field needs to be skipped or reordered to ensure proper merging of FlattenedTable generated at each recursive step.
*/
func (n *Flattener) recursiveInitFieldGroupConversion(group any, groupJsonPathKey path.JSONPath) (*FieldGroupConversion, error) {
	const FunctionName = "recursiveInitFieldGroupConversion"

	fieldGroupProp, err := core.AsJsonObject(group)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("group not JsonObject").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	groupFields, err := core.GetGroupFields(fieldGroupProp)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get group fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	groupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fieldGroupProp)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get group read order of fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	groupConversion := &FieldGroupConversion{
		FieldColumnPosition: &fieldcolumns.FieldColumnPosition{
			FieldGroupJsonPathKey: groupJsonPathKey,
		},
		GroupFields: make([]*FieldGroupConversion, 0),
	}

	for _, fgKeySuffix := range groupReadOrderOfFields {
		// Get field group
		fgProperty, err := core.AsJsonObject(groupFields[fgKeySuffix])
		if err != nil {
			return nil, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
		}

		fgJsonPathKey, err := core.AsJSONPath(fgProperty[core.FieldGroupJsonPathKey])
		if err != nil {
			return nil, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get FieldGroupJsonPathKey for field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
		}

		// Field is a Group
		if ok := core.IsFieldAGroup(fgProperty); ok {
			// Process group as a field
			if extractAsSingleField, ok := fgProperty[core.GroupExtractAsSingleField].(bool); ok && extractAsSingleField {
				goto ProcessFgPropertyAsField
			}

			if !core.DoesFieldGroupFieldsContainNestedGroupFields(fgProperty) {
				// Process group as field with set of fields in separate columns
				if fgViewMaxNoOfValuesInSeparateColumns, _ := core.GetMaximumFlatNoOfColumns(fgProperty); fgViewMaxNoOfValuesInSeparateColumns > 0 {
					if fgGroupFields, err := core.GetGroupFields(fgProperty); err == nil {
						if fgGroupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fgProperty); err == nil {
							for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
								for _, nFgKeySuffix := range fgGroupReadOrderOfFields {
									if nFgProperty, err := core.AsJsonObject(fgGroupFields[nFgKeySuffix]); err == nil {
										if nJsonPathKey, err := core.AsJSONPath(nFgProperty[core.FieldGroupJsonPathKey]); err == nil {
											fieldColumnPosition := &fieldcolumns.FieldColumnPosition{
												FieldGroupJsonPathKey:                       nJsonPathKey,
												GroupViewParentJsonPathKey:                  fgJsonPathKey,
												GroupViewInSeparateColumns:                  true,
												GroupViewValuesInSeparateColumnsHeaderIndex: columnIndex,
												FieldJsonPathKeySuffix:                      nFgKeySuffix,
											}
											if field, ok := n.columnFields.Fields[fieldColumnPosition.JSONPath()]; ok {
												if !field.Skip {
													newFieldConversion := &FieldGroupConversion{
														FieldColumnPosition: fieldColumnPosition,
													}

													groupConversion.GroupFields = append(groupConversion.GroupFields, newFieldConversion)
												}
											}
										}
									}
								}
							}
							continue
						}
					}
				}
			}

			if fgGroupConversion, err := n.recursiveInitFieldGroupConversion(fgProperty, fgJsonPathKey); err != nil {
				if !errors.Is(err, ErrNoGroupFields) {
					return nil, err
				}
			} else {
				groupConversion.GroupFields = append(groupConversion.GroupFields, fgGroupConversion)
			}

			continue
		}

		// Process field
	ProcessFgPropertyAsField:
		// Field WITH view in separate columns
		if fgViewMaxNoOfValuesInSeparateColumns, _ := core.GetMaximumFlatNoOfColumns(fgProperty); fgViewMaxNoOfValuesInSeparateColumns > 0 {
			for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
				if field, ok := n.columnFields.Fields[gojsoncore.Ptr(fieldcolumns.FieldColumnPosition{FieldGroupJsonPathKey: fgJsonPathKey, FieldViewInSeparateColumns: true, FieldViewValuesInSeparateColumnsHeaderIndex: columnIndex}).JSONPath()]; ok {
					if !field.Skip {
						newFieldConversion := &FieldGroupConversion{
							FieldColumnPosition: &fieldcolumns.FieldColumnPosition{
								FieldGroupJsonPathKey:                       fgJsonPathKey,
								FieldViewInSeparateColumns:                  true,
								FieldViewValuesInSeparateColumnsHeaderIndex: columnIndex,
							},
						}

						groupConversion.GroupFields = append(groupConversion.GroupFields, newFieldConversion)
					}
				}
			}
			continue
		}

		// Field WITHOUT view in separate columns
		if field, ok := n.columnFields.Fields[fgJsonPathKey]; ok {
			if !field.Skip {
				newFieldConversion := &FieldGroupConversion{
					FieldColumnPosition: &fieldcolumns.FieldColumnPosition{
						FieldGroupJsonPathKey: fgJsonPathKey,
					},
				}

				groupConversion.GroupFields = append(groupConversion.GroupFields, newFieldConversion)
			}
		}
	}

	if len(groupConversion.GroupFields) == 0 {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("no group fields to extract found").WithData(gojsoncore.JsonObject{"Group": group}).WithNestedError(ErrNoGroupFields)
	}

	return groupConversion, nil
}

func (n *Flattener) WithColumnFields(value *fieldcolumns.ColumnFields) *Flattener {
	n.SetColumnFields(value)
	return n
}

func (n *Flattener) SetColumnFields(value *fieldcolumns.ColumnFields) {
	n.columnFields = value
}

func NewFlattener(metadataModel gojsoncore.JsonObject) *Flattener {
	return &Flattener{
		currentSourceObjectResult: make(FlattenedTable, 0),
		metadataModel:             metadataModel,
	}
}

/*
FlattenedRow represents a single row in a table.

The flattener will attempt to enforce that each cell (column in row) is either a slice or array for uniformity.

Default, empty or uninitialized cells are represented as an empty slice of type any ([]any{}).
*/
type FlattenedRow []reflect.Value

/*
FlattenedTable represents a 2D linear collection.

This is the result of flattening an object.
*/
type FlattenedTable []FlattenedRow

/*
Flattener converts a deeply nested mix of associative collections (like an array of objects) into a 2 dimension linear collection (like a 2D array).
*/
type Flattener struct {
	metadataModel gojsoncore.JsonObject

	// columnFields extracted fields as table columns from metadataModel
	columnFields *fieldcolumns.ColumnFields

	// currentSourceObject individual associative collection object to convert to a 2D slice. Result appended to currentSourceObjectResult.
	currentSourceObject *object.Object

	// currentSourceObjectIsALinearCollection determines how currentSourceObject will be processed.
	//
	// If true, object to flatten is assumed to be a collection of associative collections (maps and structs) thus each at the top-level will be loaded into currentSourceObject for flattening individually.
	currentSourceObjectIsALinearCollection bool

	// currentSourceObjectResult holds the current result of flattening the currentSourceObject.
	//
	// Will be written into destination.
	currentSourceObjectResult FlattenedTable

	// fieldGroupConversion data (tree of fields/groups) to use when converting the currentSourceObject.
	fieldGroupConversion *FieldGroupConversion
}

/*
FieldGroupConversion represents tree of field/groups to read for Flattener.
*/
type FieldGroupConversion struct {
	FieldColumnPosition *fieldcolumns.FieldColumnPosition

	GroupFields []*FieldGroupConversion
}

func DefaultEmptyColumn() reflect.Value {
	return reflect.ValueOf([]any(nil))
}
