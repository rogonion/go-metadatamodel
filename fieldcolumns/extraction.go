package fieldcolumns

import (
	"fmt"

	"github.com/brunoga/deep"
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/core"
)

/*
Extract retrieves recursively goes through Extraction.metadataModel in order and sets the following properties:
  - Extraction.columnFields
  - Extraction.RepositionFieldColumns
  - Extraction.readOrderOfColumnFields
*/
func (n *Extraction) Extract() (*ColumnFields, error) {
	n.columnFields = NewColumnFields()

	if err := n.recursiveExtract(n.metadataModel, gojsoncore.JsonObject{}, nil); err != nil {
		return nil, err
	}
	return n.columnFields, nil
}

func (n *Extraction) recursiveExtract(group any, matchingGroupProperties gojsoncore.JsonObject, nextGroupFieldPosition *FieldColumnPosition) error {
	const FunctionName = "recursiveExtract"

	fieldGroupProp, err := core.AsJsonObject(group)
	if err != nil {
		return NewError().WithFunctionName(FunctionName).WithMessage("group not JsonObject").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group, "MatchingProperties": matchingGroupProperties})
	}

	groupFields, err := core.GetGroupFields(fieldGroupProp)
	if err != nil {
		return NewError().WithFunctionName(FunctionName).WithMessage("get group fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group, "MatchingProperties": matchingGroupProperties})
	}

	groupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fieldGroupProp)
	if err != nil {
		return NewError().WithFunctionName(FunctionName).WithMessage("get group read order of fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group, "MatchingProperties": matchingGroupProperties})
	}

	for _, fgKeySuffix := range groupReadOrderOfFields {
		// Get field group
		fgProperty, err := core.AsJsonObject(groupFields[fgKeySuffix])
		if err != nil {
			return NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
		}

		// Get nested matching props
		fgMatchingProperties := make(gojsoncore.JsonObject)
		if matchingGroupProperties != nil {
			fgMatchingProperties, _ = deep.Copy(matchingGroupProperties)
		}

		if n.add.IsValid() {
			if matchingProps := n.add.MatchingProps(fgProperty); len(matchingProps) > 0 {
				core.MergeRightJsonObjectIntoLeft(fgMatchingProperties, matchingProps)
			}
		}
		if n.skip.IsValid() {
			if matchingProps := n.skip.MatchingProps(fgProperty); len(matchingProps) > 0 {
				core.MergeRightJsonObjectIntoLeft(fgMatchingProperties, matchingProps)
			}
		}

		nextFieldGroupPosition := ExtractFieldColumnPosition(fgProperty)
		if nextFieldGroupPosition == nil {
			nextFieldGroupPosition = nextGroupFieldPosition
		}

		fgJsonPathKey, err := core.AsJSONPath(fgProperty[core.FieldGroupJsonPathKey])
		if err != nil {
			return NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get FieldGroupJsonPathKey for field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
		}

		// Field is a group
		if core.IsFieldAGroup(fgProperty) {
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
									var newField gojsoncore.JsonObject
									if value, err := core.AsJsonObject(fgGroupFields[nFgKeySuffix]); err == nil {
										newField, err = n.createNewField(value, fgMatchingProperties)
										if err != nil {
											return err
										}
									} else {
										return NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"FgGroupFields": fgGroupFields})
									}
									n.addFlatColumnContext(newField, columnIndex)
									jsonPathKey, err := core.AsJSONPath(newField[core.FieldGroupJsonPathKey])
									if err != nil {
										return err
									}
									if err := n.appendField(newField, &FieldColumnPosition{
										FieldGroupJsonPathKey:                       jsonPathKey,
										GroupViewInSeparateColumns:                  true,
										GroupViewValuesInSeparateColumnsHeaderIndex: columnIndex,
										GroupViewParentJsonPathKey:                  fgJsonPathKey,
										FieldJsonPathKeySuffix:                      nFgKeySuffix,
									}); err != nil {
										return err
									}

									if nextFieldGroupPosition != nil {
										n.setRepositionForFieldColumn(nextFieldGroupPosition)
										nextFieldGroupPosition.FieldGroupJsonPathKey = jsonPathKey
										nextFieldGroupPosition.GroupViewInSeparateColumns = true
										nextFieldGroupPosition.GroupViewValuesInSeparateColumnsHeaderIndex = columnIndex
										nextFieldGroupPosition.GroupViewParentJsonPathKey = fgJsonPathKey
										nextFieldGroupPosition.FieldJsonPathKeySuffix = nFgKeySuffix
										nextFieldGroupPosition.FieldGroupPositionBefore = false
									}
								}
							}
							continue
						}
					}
				}
			}

			// Process group fields
			if err := n.recursiveExtract(fgProperty, fgMatchingProperties, nextFieldGroupPosition); err != nil {
				return err
			}
			continue
		}

		// Process field
	ProcessFgPropertyAsField:
		// Field WITH view in separate columns
		if fgViewMaxNoOfValuesInSeparateColumns, _ := core.GetMaximumFlatNoOfColumns(fgProperty); fgViewMaxNoOfValuesInSeparateColumns > 0 {
			for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
				newField, err := n.createNewField(fgProperty, fgMatchingProperties)
				if err != nil {
					return err
				}
				n.addFlatColumnContext(newField, columnIndex)
				if err := n.appendField(newField, &FieldColumnPosition{FieldGroupJsonPathKey: fgJsonPathKey, FieldViewInSeparateColumns: true, FieldViewValuesInSeparateColumnsHeaderIndex: columnIndex}); err != nil {
					return err
				}

				if nextFieldGroupPosition != nil {
					n.setRepositionForFieldColumn(nextFieldGroupPosition)
					nextFieldGroupPosition.FieldGroupJsonPathKey = fgJsonPathKey
					nextFieldGroupPosition.FieldViewInSeparateColumns = true
					nextFieldGroupPosition.FieldViewValuesInSeparateColumnsHeaderIndex = columnIndex
					nextFieldGroupPosition.FieldGroupPositionBefore = false
				}
			}
			continue
		}

		// Field WITHOUT view in separate columns
		newField, err := n.createNewField(fgProperty, fgMatchingProperties)
		if err != nil {
			return err
		}

		if err := n.appendField(newField, &FieldColumnPosition{FieldGroupJsonPathKey: fgJsonPathKey}); err != nil {
			return err
		}

		if nextFieldGroupPosition != nil {
			n.setRepositionForFieldColumn(nextFieldGroupPosition)
			nextFieldGroupPosition.FieldGroupJsonPathKey = fgJsonPathKey
			nextFieldGroupPosition.FieldViewInSeparateColumns = false
			nextFieldGroupPosition.FieldViewValuesInSeparateColumnsHeaderIndex = 0
			nextFieldGroupPosition.FieldGroupPositionBefore = false
		}
	}

	return nil
}

/*
setRepositionForFieldColumn

Set after call to Extraction.appendField for the same field.
*/
func (n *Extraction) setRepositionForFieldColumn(value *FieldColumnPosition) {
	value.SourceIndex = len(n.columnFields.ReadOrderOfColumnFields) - 1
	n.columnFields.RepositionFieldColumns = append(n.columnFields.RepositionFieldColumns, *value)
}

/*
appendField

Add new field to Extraction.columnFields.
*/
func (n *Extraction) appendField(field gojsoncore.JsonObject, fieldColumnPosition *FieldColumnPosition) error {
	const FunctionName = "appendField"

	newFieldColumn := &ColumnField{
		Property:                       field,
		IndexInReadOrderOfColumnFields: len(n.columnFields.ReadOrderOfColumnFields),
	}
	if n.schema != nil {
		if jsonPathToValue, err := core.NewJsonPathToValue().Get(fieldColumnPosition.FieldGroupJsonPathKey, nil); err == nil {
			if fieldSchema, err := schema.GetSchemaAtPath(jsonPathToValue, n.schema); err == nil {
				newFieldColumn.Schema = fieldSchema
			}
		}
	}

	n.columnFields.Fields[fieldColumnPosition.JSONPath()] = newFieldColumn
	n.columnFields.ReadOrderOfColumnFields = append(n.columnFields.ReadOrderOfColumnFields, fieldColumnPosition)
	n.columnFields.CurrentIndexOfReadOrderOfColumnFields = append(n.columnFields.CurrentIndexOfReadOrderOfColumnFields, newFieldColumn.IndexInReadOrderOfColumnFields)
	return nil
}

/*
createNewField

1. Deep copy field to not affect original field inside metadata model.
2. Append matchingGroupProperties to new field.
*/
func (n *Extraction) createNewField(original gojsoncore.JsonObject, matchingGroupProperties gojsoncore.JsonObject) (gojsoncore.JsonObject, error) {
	const FunctionName = "createNewField"

	newField, err := deep.Copy(original)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("deep copy original field failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Original": original})
	}

	core.MergeRightJsonObjectIntoLeft(newField, matchingGroupProperties)

	return newField, nil
}

/*
addFlatColumnContext

For field with core.FieldGroupViewValuesInSeparateColumns property:
1. Set core.FieldViewValuesInSeparateColumnsHeaderIndex in flat view columns.
2. Add core.FieldViewValuesInSeparateColumnsHeaderIndex to core.FieldGroupName.
*/
func (n *Extraction) addFlatColumnContext(field gojsoncore.JsonObject, columnIndex int) {
	field[core.FieldViewValuesInSeparateColumnsHeaderIndex] = columnIndex

	if fgHeaderFormat, ok := field[core.FieldViewValuesInSeparateColumnsHeaderFormat].(string); ok && fgHeaderFormat != "" {
		field[core.FieldGroupName] = string(core.ArrayPathRegexSearch().ReplaceAll([]byte(fgHeaderFormat), fmt.Appendf(nil, "%d", columnIndex+1)))
	} else {
		field[core.FieldGroupName] = fmt.Sprintf("%s %d", core.GetFieldGroupName(field, ""), columnIndex+1)
	}
}

func (n *Extraction) WithSchema(value schema.Schema) *Extraction {
	n.SetSchema(value)
	return n
}

func (n *Extraction) SetSchema(value schema.Schema) {
	n.schema = value
}

func (n *Extraction) WithAdd(value core.FieldGroupPropertiesMatch) *Extraction {
	n.SetAdd(value)
	return n
}

func (n *Extraction) SetAdd(value core.FieldGroupPropertiesMatch) {
	n.add = value
}

func (n *Extraction) WithSkip(value core.FieldGroupPropertiesMatch) *Extraction {
	n.SetSkip(value)
	return n
}

func (n *Extraction) SetSkip(value core.FieldGroupPropertiesMatch) {
	n.skip = value
}

func NewColumnFieldsExtraction(metadataModel gojsoncore.JsonObject) *Extraction {
	n := new(Extraction)
	n.metadataModel = metadataModel
	return n
}

type Extraction struct {
	metadataModel gojsoncore.JsonObject

	// schema that represents data for metadataModel in object form.
	//
	// During columnFields extraction, schema for specific field may be extracted as well into ColumnField.Schema.
	schema schema.Schema

	// columnFields extracted field columns.
	columnFields *ColumnFields

	// skip as each field/group is being processed, if one of the properties matches, then all nested fields/groups properties should contain the result of core.FieldGroupPropertiesMatchingProps.
	//
	// Useful for scenarios such as disabling all nested fields/groups in a group with property core.FieldGroupViewDisable set to `true` in the parent group alone.
	skip core.FieldGroupPropertiesMatch

	// add as each field/group is being processed, if one of the properties matches, then all nested fields/groups properties should contain the result of core.FieldGroupPropertiesMatchingProps.
	//
	// Useful for scenarios such as disabling all nested fields/groups in a group with property core.FieldGroupViewDisable set to `true` in the parent group alone.
	add core.FieldGroupPropertiesMatch
}
