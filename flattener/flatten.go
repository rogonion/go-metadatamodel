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

func (n *Flattener) recursiveConvert(currentFlattenedObject FlattenedTable, groupConversion *FieldGroupConversion, arrayIndexPlaceholders []int) (FlattenedTable, error) {
}

func (n *Flattener) Convert(object object.Object) error {
	n.currentSourceObjectIsAnArray = object.GetSourceReflected().Kind() == reflect.Slice || n.currentSourceObject.GetSourceReflected().Kind() == reflect.Array
	if n.currentFlattenedTable == nil {
		n.currentFlattenedTable = make(FlattenedTable, 0)
	}

	if n.fieldGroupConversion == nil {
		if value, err := n.recursiveInitFieldGroupConversion(n.metadataModel, path.JSONPath(path.JsonpathKeyRoot)); err != nil {
			return err
		} else {
			n.fieldGroupConversion = value
		}
	}

	if n.currentSourceObjectIsAnArray {
		object.ForEach(path.JSONPath(path.JsonpathKeyRoot+path.JsonpathDotNotation+path.JsonpathLeftBracket+path.JsonpathKeyIndexAll+path.JsonpathRightBracket), func(jsonPath path.RecursiveDescentSegment, value reflect.Value) bool {

		})
	} else {
		n.currentSourceObject = object
		if value, err := n.recursiveConvert(make(FlattenedTable, 0), n.fieldGroupConversion, make([]int, 0)); err != nil {
			return err
		} else {
			n.currentFlattenedTable = append(n.currentFlattenedTable, value...)
		}
	}
}

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
		FieldGroupJsonPathKey: groupJsonPathKey,
		GroupFields:           make([]*FieldGroupConversion, 0),
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

			if core.DoesFieldGroupFieldsContainNestedGroupFields(fgProperty) {
				// Process group as field with set of fields in separate columns
				if fgViewMaxNoOfValuesInSeparateColumns, _ := core.GetMaximumFlatNoOfColumns(fgProperty); fgViewMaxNoOfValuesInSeparateColumns > 0 {
					if fgGroupFields, err := core.GetGroupFields(fgProperty); err == nil {
						if fgGroupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fgProperty); err == nil {
							for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
								for _, nFgKeySuffix := range fgGroupReadOrderOfFields {
									if nFgProperty, err := core.AsJsonObject(fgGroupFields[nFgKeySuffix]); err == nil {
										if nJsonPathKey, err := core.AsJSONPath(nFgProperty[core.FieldGroupJsonPathKey]); err == nil {
											if field, ok := n.columnFields.Fields[gojsoncore.Ptr(fieldcolumns.ReadOrderOfColumnField{FieldGroupJsonPathKey: nJsonPathKey, ViewInSeparateColumns: true, IndexOfValueInSeparateColumns: columnIndex}).JSONPath()]; ok {
												if !field.Skip {
													newFieldConversion := &FieldGroupConversion{
														FieldGroupJsonPathKey:                       nJsonPathKey,
														FieldGroupViewValueInSeparateColumns:        true,
														FieldGroupIndexInViewValueInSeparateColumns: columnIndex,
													}

													if joinSymbol, ok := field.Property[core.FieldMultipleValuesJoinSymbol].(string); ok && joinSymbol != "" {
														newFieldConversion.FieldJoinSymbol = joinSymbol
													} else {
														newFieldConversion.FieldJoinSymbol = ","
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
				if field, ok := n.columnFields.Fields[gojsoncore.Ptr(fieldcolumns.ReadOrderOfColumnField{FieldGroupJsonPathKey: fgJsonPathKey, ViewInSeparateColumns: true, IndexOfValueInSeparateColumns: columnIndex}).JSONPath()]; ok {
					if !field.Skip {
						newFieldConversion := &FieldGroupConversion{
							FieldGroupJsonPathKey:                       fgJsonPathKey,
							FieldGroupViewValueInSeparateColumns:        true,
							FieldGroupIndexInViewValueInSeparateColumns: columnIndex,
						}

						if joinSymbol, ok := field.Property[core.FieldMultipleValuesJoinSymbol].(string); ok && joinSymbol != "" {
							newFieldConversion.FieldJoinSymbol = joinSymbol
						} else {
							newFieldConversion.FieldJoinSymbol = ","
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
					FieldGroupJsonPathKey: fgJsonPathKey,
				}

				if joinSymbol, ok := field.Property[core.FieldMultipleValuesJoinSymbol].(string); ok && joinSymbol != "" {
					newFieldConversion.FieldJoinSymbol = joinSymbol
				} else {
					newFieldConversion.FieldJoinSymbol = ","
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

func (n *Flattener) ResetCurrentFlattenedTable() {
	n.currentFlattenedTable = make(FlattenedTable, 0)
}

func (n *Flattener) WithFlattenFieldArray(value bool) *Flattener {
	n.SetFlattenFieldArray(value)
	return n
}

func (n *Flattener) SetFlattenFieldArray(value bool) *Flattener {
	n.flattenFieldArray = value
	return n
}

func (n *Flattener) WithColumnFields(value *fieldcolumns.ColumnFields) *Flattener {
	n.SetColumnFields(value)
	return n
}

func (n *Flattener) SetColumnFields(value *fieldcolumns.ColumnFields) {
	n.columnFields = value
	n.currentReadOrderOfColumnFields = n.columnFields.GetCurrentIndexOfReadOrderOfFields()
}

func (n *Flattener) WithMetadataModel(value gojsoncore.JsonObject) *Flattener {
	n.SetMetadataModel(value)
	return n
}

func (n *Flattener) SetMetadataModel(value gojsoncore.JsonObject) {
	n.metadataModel = value
}

func NewFlattener() *Flattener {
	return &Flattener{
		currentFlattenedTable: make(FlattenedTable, 0),
	}
}

type FlattenedRow []any

type FlattenedTable []FlattenedRow

type Flattener struct {
	metadataModel gojsoncore.JsonObject

	// columnFields extracted fields as table columns from metadataModel
	columnFields *fieldcolumns.ColumnFields

	// currentReadOrderOfColumnFields will be extracted using columnFields.GetCurrentIndexOfReadOrderOfFields.
	currentReadOrderOfColumnFields []int

	// currentSourceObject object to convert to 2D slice
	currentSourceObject object.Object

	// currentSourceObjectIsAnArray determines how currentSourceObject will be processed.
	currentSourceObjectIsAnArray bool

	// currentFlattenedTable holds current result of flattening currentSourceObject.
	currentFlattenedTable FlattenedTable

	// fieldGroupConversion current tree of fields/groups to read when processing currentSourceObject
	fieldGroupConversion *FieldGroupConversion

	//flattenFieldArray if field value is a slice/array, extract it to a single value before insertion.
	flattenFieldArray bool
}

/*
FieldGroupConversion represents tree of field/groups to read for Flattener.
*/
type FieldGroupConversion struct {
	FieldGroupJsonPathKey                       path.JSONPath
	FieldGroupViewValueInSeparateColumns        bool
	FieldGroupIndexInViewValueInSeparateColumns int

	FieldJoinSymbol string

	GroupFields            []*FieldGroupConversion
	GroupReadOrderOfFields core.MetadataModelGroupReadOrderOfFields
}
