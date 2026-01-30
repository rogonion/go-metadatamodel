package fieldcolumns

import (
	"fmt"
	"reflect"

	"github.com/brunoga/deep"
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/iter"
)

/*
getFieldColumnsFromMetadataModel used to test field extraction.
*/
func getFieldColumnsFromMetadataModel(metadataModel gojsoncore.JsonObject, sch schema.Schema, matchingProps core.FieldGroupPropertiesMatch) *ColumnFields {
	columnFields := NewColumnFields()

	currentFieldIndex := 0
	iter.ForEach(metadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
		if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
			if core.IsFieldAGroup(fieldGroup) {
				if value, ok := fieldGroup[core.GroupExtractAsSingleField].(bool); ok && value {
					goto ProcessAsField
				}

				if fgViewMaxNoOfValuesInSeparateColumns, _ := core.GetMaximumFlatNoOfColumns(fieldGroup); fgViewMaxNoOfValuesInSeparateColumns > 0 {
					if !core.DoesFieldGroupFieldsContainNestedGroupFields(fieldGroup) {
						if fgGroupFields, err := core.GetGroupFields(fieldGroup); err == nil {
							if fgGroupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fieldGroup); err == nil {
								for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
									for _, nFgKeySuffix := range fgGroupReadOrderOfFields {
										if value, err := core.AsJsonObject(fgGroupFields[nFgKeySuffix]); err == nil {
											if newField, err := deep.Copy(value); err == nil {
												if nJsonPathKey, err := core.AsJSONPath(newField[core.FieldGroupJsonPathKey]); err == nil {
													if matchingProps.IsValid() {
														if mp := matchingProps.MatchingProps(newField); len(mp) > 0 {
															core.MergeRightJsonObjectIntoLeft(newField, mp)
														}
													}
													newField[core.FieldViewValuesInSeparateColumnsHeaderIndex] = columnIndex

													if fgHeaderFormat, ok := newField[core.FieldViewValuesInSeparateColumnsHeaderFormat].(string); ok && fgHeaderFormat != "" {
														newField[core.FieldGroupName] = string(core.ArrayPathRegexSearch().ReplaceAll([]byte(fgHeaderFormat), fmt.Appendf(nil, "%d", columnIndex+1)))
													} else {
														newField[core.FieldGroupName] = fmt.Sprintf("%s %d", core.GetFieldGroupName(newField, ""), columnIndex+1)
													}

													fieldColumnPosition := &FieldColumnPosition{
														FieldGroupJsonPathKey:                       nJsonPathKey,
														GroupViewInSeparateColumns:                  true,
														GroupViewValuesInSeparateColumnsHeaderIndex: columnIndex,
														GroupViewParentJsonPathKey:                  jsonPathKey,
														FieldJsonPathKeySuffix:                      nFgKeySuffix,
													}
													columnFields.Fields[fieldColumnPosition.JSONPath()] = &ColumnField{
														Property:                               newField,
														IndexInOriginalReadOrderOfColumnFields: len(columnFields.OriginalReadOrderOfColumnFields),
													}
													columnFields.OriginalReadOrderOfColumnFields = append(columnFields.OriginalReadOrderOfColumnFields, fieldColumnPosition)

													if pathToSchema, err := core.NewJsonPathToValue().Get(nJsonPathKey, nil); err == nil {
														if fieldGroupSchema, err := schema.GetSchemaAtPath(pathToSchema, sch); err == nil {
															columnFields.Fields[fieldColumnPosition.JSONPath()].Schema = fieldGroupSchema
														}
													}

													columnFields.RepositionedReadOrderOfColumnFields = append(columnFields.RepositionedReadOrderOfColumnFields, currentFieldIndex)
													currentFieldIndex++
												}
											}
										}
									}
								}
								return false, true
							}
						}
					}
				}
				return false, false
			}

		ProcessAsField:
			if fgViewMaxNoOfValuesInSeparateColumns, _ := core.GetMaximumFlatNoOfColumns(fieldGroup); fgViewMaxNoOfValuesInSeparateColumns > 0 {
				for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
					if newField, err := deep.Copy(fieldGroup); err == nil {
						newField[core.FieldViewValuesInSeparateColumnsHeaderIndex] = columnIndex

						if fgHeaderFormat, ok := newField[core.FieldViewValuesInSeparateColumnsHeaderFormat].(string); ok && fgHeaderFormat != "" {
							newField[core.FieldGroupName] = string(core.ArrayPathRegexSearch().ReplaceAll([]byte(fgHeaderFormat), fmt.Appendf(nil, "%d", columnIndex+1)))
						} else {
							newField[core.FieldGroupName] = fmt.Sprintf("%s %d", core.GetFieldGroupName(newField, ""), columnIndex+1)
						}

						fieldColumnPosition := &FieldColumnPosition{
							FieldGroupJsonPathKey:                       jsonPathKey,
							FieldViewInSeparateColumns:                  true,
							FieldViewValuesInSeparateColumnsHeaderIndex: columnIndex,
						}
						columnFields.Fields[fieldColumnPosition.JSONPath()] = &ColumnField{
							Property:                               newField,
							IndexInOriginalReadOrderOfColumnFields: len(columnFields.OriginalReadOrderOfColumnFields),
						}
						columnFields.OriginalReadOrderOfColumnFields = append(columnFields.OriginalReadOrderOfColumnFields, fieldColumnPosition)

						if pathToSchema, err := core.NewJsonPathToValue().Get(jsonPathKey, nil); err == nil {
							if fieldGroupSchema, err := schema.GetSchemaAtPath(pathToSchema, sch); err == nil {
								columnFields.Fields[fieldColumnPosition.JSONPath()].Schema = fieldGroupSchema
							}
						}

						columnFields.RepositionedReadOrderOfColumnFields = append(columnFields.RepositionedReadOrderOfColumnFields, currentFieldIndex)
						currentFieldIndex++
					}
				}
			} else {
				fieldColumnPosition := &FieldColumnPosition{
					FieldGroupJsonPathKey: jsonPathKey,
				}
				if newField, err := deep.Copy(fieldGroup); err == nil {
					columnFields.Fields[fieldColumnPosition.JSONPath()] = &ColumnField{
						Property:                               newField,
						IndexInOriginalReadOrderOfColumnFields: len(columnFields.OriginalReadOrderOfColumnFields),
					}
					columnFields.OriginalReadOrderOfColumnFields = append(columnFields.OriginalReadOrderOfColumnFields, fieldColumnPosition)

					if pathToSchema, err := core.NewJsonPathToValue().Get(jsonPathKey, nil); err == nil {
						if fieldGroupSchema, err := schema.GetSchemaAtPath(pathToSchema, sch); err == nil {
							columnFields.Fields[fieldColumnPosition.JSONPath()].Schema = fieldGroupSchema
						}
					}

					columnFields.RepositionedReadOrderOfColumnFields = append(columnFields.RepositionedReadOrderOfColumnFields, currentFieldIndex)
					currentFieldIndex++
				}
			}
		}
		return false, false
	})

	return columnFields
}

func ExtractFieldColumnPosition(field gojsoncore.JsonObject) *FieldColumnPosition {
	if value, err := core.AsJsonObject(field[core.FieldColumnPosition]); err == nil {
		if fieldGroupJsonPathKey, err := core.AsJSONPath(value[core.FieldGroupJsonPathKey]); err == nil {
			fieldColumnPosition := &FieldColumnPosition{
				FieldGroupJsonPathKey: fieldGroupJsonPathKey,
			}
			if positionBefore, ok := value[core.FieldGroupPositionBefore].(bool); ok {
				fieldColumnPosition.FieldGroupPositionBefore = positionBefore
			}
			if chi, ok := value[core.FieldViewValuesInSeparateColumnsHeaderIndex]; ok {
				var columnHeaderIndex int
				if err := schema.NewConversion().Convert(chi, &schema.DynamicSchemaNode{Type: reflect.TypeOf(0), Kind: reflect.Int}, &columnHeaderIndex); err == nil {
					fieldColumnPosition.FieldViewValuesInSeparateColumnsHeaderIndex = columnHeaderIndex
				}
			}

			return fieldColumnPosition
		}
	}

	return nil
}
