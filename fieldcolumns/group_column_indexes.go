package fieldcolumns

import (
	"fmt"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
)

func (n *GroupsColumnsIndexesRetrieval) recursiveSetPrimaryKeysFromGroupFields(group gojsoncore.JsonObject) error {
	const FunctionName = "recursiveSetPrimaryKeysFromGroupFields"

	groupFields, err := core.GetGroupFields(group)
	if err != nil {
		return NewError().WithFunctionName(FunctionName).WithMessage("get group fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	groupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(group)
	if err != nil {
		return NewError().WithFunctionName(FunctionName).WithMessage("get group read order of fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	for _, fgKeySuffix := range groupReadOrderOfFields {
		fgProperty, err := core.AsJsonObject(groupFields[fgKeySuffix])
		if err != nil {
			return NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
		}

		fgJsonPathKey, err := core.AsJSONPath(fgProperty[core.FieldGroupJsonPathKey])
		if err != nil {
			return NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get FieldGroupJsonPathKey for field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
		}

		if value, ok := fgProperty[core.FieldGroupIsPrimaryKey].(bool); !ok || !value {
			continue
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
											if err := n.appendFieldColumnIndex(&FieldColumnPosition{
												FieldGroupJsonPathKey:                       nJsonPathKey,
												GroupViewParentJsonPathKey:                  fgJsonPathKey,
												GroupViewInSeparateColumns:                  true,
												GroupViewValuesInSeparateColumnsHeaderIndex: columnIndex,
												FieldJsonPathKeySuffix:                      nFgKeySuffix,
											}, true, false); err != nil {
												return err
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

			if err := n.recursiveSetPrimaryKeysFromGroupFields(fgProperty); err != nil {
				return err
			}
			continue
		}

		// Process field
	ProcessFgPropertyAsField:
		// Field WITH view in separate columns
		if fgViewMaxNoOfValuesInSeparateColumns, _ := core.GetMaximumFlatNoOfColumns(fgProperty); fgViewMaxNoOfValuesInSeparateColumns > 0 {
			for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
				if err := n.appendFieldColumnIndex(&FieldColumnPosition{
					FieldGroupJsonPathKey:                       fgJsonPathKey,
					FieldViewInSeparateColumns:                  true,
					FieldViewValuesInSeparateColumnsHeaderIndex: columnIndex,
				}, true, false); err != nil {
					return err
				}
			}
			continue
		}

		// Field WITHOUT view in separate columns
		if err := n.appendFieldColumnIndex(&FieldColumnPosition{FieldGroupJsonPathKey: fgJsonPathKey}, true, false); err != nil {
			return err
		}
	}

	return nil
}

// Get retrieves the column indexes (Primary and All) for a specific group within the metadata model.
func (n *GroupsColumnsIndexesRetrieval) Get(group gojsoncore.JsonObject) (*GroupColumnIndexes, error) {
	const FunctionName = "Get"
	n.primary = make([]int, 0)
	n.all = make([]int, 0)

	groupFields, err := core.GetGroupFields(group)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get group fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	groupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(group)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get group read order of fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	for _, fgKeySuffix := range groupReadOrderOfFields {
		fgProperty, err := core.AsJsonObject(groupFields[fgKeySuffix])
		if err != nil {
			return nil, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
		}

		isPrimary := false
		if value, ok := fgProperty[core.FieldGroupIsPrimaryKey].(bool); ok && value {
			isPrimary = value
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
											if err := n.appendFieldColumnIndex(&FieldColumnPosition{
												FieldGroupJsonPathKey:                       nJsonPathKey,
												GroupViewParentJsonPathKey:                  fgJsonPathKey,
												GroupViewInSeparateColumns:                  true,
												GroupViewValuesInSeparateColumnsHeaderIndex: columnIndex,
												FieldJsonPathKeySuffix:                      nFgKeySuffix,
											}, isPrimary, true); err != nil {
												return nil, err
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

			if isPrimary {
				if err := n.recursiveSetPrimaryKeysFromGroupFields(fgProperty); err != nil {
					return nil, err
				}
			}
			continue
		}

		// Process field
	ProcessFgPropertyAsField:
		// Field WITH view in separate columns
		if fgViewMaxNoOfValuesInSeparateColumns, _ := core.GetMaximumFlatNoOfColumns(fgProperty); fgViewMaxNoOfValuesInSeparateColumns > 0 {
			for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
				if err := n.appendFieldColumnIndex(&FieldColumnPosition{
					FieldGroupJsonPathKey:                       fgJsonPathKey,
					FieldViewInSeparateColumns:                  true,
					FieldViewValuesInSeparateColumnsHeaderIndex: columnIndex,
				}, isPrimary, true); err != nil {
					return nil, err
				}
			}
			continue
		}

		// Field WITHOUT view in separate columns
		if err := n.appendFieldColumnIndex(&FieldColumnPosition{FieldGroupJsonPathKey: fgJsonPathKey}, isPrimary, true); err != nil {
			return nil, err
		}
	}

	return &GroupColumnIndexes{
		Primary: n.primary,
		All:     n.all,
	}, nil
}

func (n *GroupsColumnsIndexesRetrieval) appendFieldColumnIndex(fieldColumnPosition *FieldColumnPosition, isPrimary bool, isAll bool) error {
	const FunctionName = "appendFieldColumnIndex"

	if columnField, ok := n.columnFields.GetColumnFieldByFieldGroupJsonPathKey(fieldColumnPosition.JSONPath()); ok {
		if isAll {
			n.all = append(n.all, columnField.IndexInUnskippedColumnFields)
		}
		if isPrimary {
			n.primary = append(n.primary, columnField.IndexInUnskippedColumnFields)
		}
	} else {
		return NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("column '%s' not found in columnFields", fieldColumnPosition.JSONPath()))
	}

	return nil
}

// NewGroupsColumnsIndexesRetrieval creates a new GroupsColumnsIndexesRetrieval instance.
func NewGroupsColumnsIndexesRetrieval(columnFields *ColumnFields) *GroupsColumnsIndexesRetrieval {
	return &GroupsColumnsIndexesRetrieval{
		columnFields: columnFields,
	}
}

// GroupsColumnsIndexesRetrieval is a helper to retrieve column indexes for a group.
type GroupsColumnsIndexesRetrieval struct {
	columnFields *ColumnFields
	primary      []int
	all          []int
}

/*
GroupColumnIndexes fast lookup index to quickly indentify which columns to read at current group
*/
type GroupColumnIndexes struct {
	// Primary key columns.
	Primary []int

	// All columns.
	All []int
}
