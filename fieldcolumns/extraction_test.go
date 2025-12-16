package fieldcolumns

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/internal"
	"github.com/rogonion/go-metadatamodel/iter"
	"github.com/rogonion/go-metadatamodel/testdata"
)

func TestFieldColumns_Extraction(t *testing.T) {
	for testData := range extractionTestData {
		fcExtraction := NewColumnFieldsExtraction(testData.MetadataModel).WithSchema(testData.Schema).WithSkip(testData.NestedSkip).WithAdd(testData.NestedAdd)
		columnFields, err := fcExtraction.Extract()
		if testData.ExpectedOk && err != nil {
			t.Error(
				testData.TestTitle, "\n",
				"expected ok=", testData.ExpectedOk, "got error=", err, "\n",
			)
		}
		if err != nil && testData.LogErrorsIfExpectedNotOk {
			var fieldColumnsError *core.Error
			if errors.As(err, &fieldColumnsError) {
				t.Error(
					testData.TestTitle, "\n",
					"-----Error Details-----", "\n",
					fieldColumnsError.String(), "\n",
					"-----------------------",
				)
			}
		}

		if !reflect.DeepEqual(columnFields, testData.ExpectedColumnFields) {
			t.Error(
				testData.TestTitle, "\n",
				"result ColumnFields not equal to testData.ExpectedColumnFields\n",
				"ExpectedColumnFields=", gojsoncore.JsonStringifyMust(testData.ExpectedColumnFields), "\n",
				"result=", gojsoncore.JsonStringifyMust(columnFields),
			)
		}
	}
}

type extractionData struct {
	internal.TestData
	MetadataModel        gojsoncore.JsonObject
	Schema               schema.Schema
	NestedSkip           core.FieldGroupPropertiesMatch
	NestedAdd            core.FieldGroupPropertiesMatch
	ExpectedOk           bool
	ExpectedColumnFields *ColumnFields
}

func extractionTestData(yield func(data *extractionData) bool) {
	metadataModel := testdata.UserMetadataModel(nil)
	sch := testdata.UserSchema()
	testCaseIndex := 1
	columnFields := getFieldColumnsFromMetadataModel(metadataModel, sch, nil)
	if !yield(
		&extractionData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d", testCaseIndex),
			},
			MetadataModel:        metadataModel,
			Schema:               sch,
			ExpectedOk:           true,
			ExpectedColumnFields: columnFields,
		},
	) {
		return
	}

	metadataModel = testdata.EmployeeMetadataModel(nil)
	sch = testdata.EmployeeSchema()
	columnFields = getFieldColumnsFromMetadataModel(metadataModel, sch, nil)
	testCaseIndex++
	if !yield(
		&extractionData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d", testCaseIndex),
			},
			MetadataModel:        metadataModel,
			Schema:               sch,
			ExpectedOk:           true,
			ExpectedColumnFields: columnFields,
		},
	) {
		return
	}

	metadataModel = iter.Map(testdata.ProductMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
		if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "Name" {
			fieldGroup[core.FieldGroupViewValuesInSeparateColumns] = true
			fieldGroup[core.FieldGroupViewMaxNoOfValuesInSeparateColumns] = 3
		}
		return fieldGroup, true
	}).(gojsoncore.JsonObject)
	sch = testdata.ProductSchema()
	columnFields = getFieldColumnsFromMetadataModel(metadataModel, sch, nil)
	testCaseIndex++
	if !yield(
		&extractionData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d: Test property '%s' for group field that is a field", testCaseIndex, core.FieldGroupViewValuesInSeparateColumns),
			},
			MetadataModel:        metadataModel,
			Schema:               sch,
			ExpectedOk:           true,
			ExpectedColumnFields: columnFields,
		},
	) {
		return
	}

	repositionFieldColumns := make(RepositionFieldColumns, 0)
	{
		var idJsonPathKey path.JSONPath
		metadataModel = iter.Map(metadataModel, func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if name, ok := fieldGroup[core.FieldGroupName].(string); ok {
				switch name {
				case "ID":
					if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
						idJsonPathKey = jsonPathKey
					}
				case "Price":
					if idJsonPathKey != "" {
						fieldGroup[core.FieldColumnPosition] = gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey: idJsonPathKey,
						}
						repositionFieldColumns = append(repositionFieldColumns, FieldColumnPosition{
							SourceIndex:           4,
							FieldGroupJsonPathKey: idJsonPathKey,
						})
					}
				}
			}
			return fieldGroup, true
		}).(gojsoncore.JsonObject)
		columnFields = getFieldColumnsFromMetadataModel(metadataModel, sch, nil)
		columnFields.RepositionFieldColumns = repositionFieldColumns
		testCaseIndex++
		if !yield(
			&extractionData{
				TestData: internal.TestData{
					TestTitle: fmt.Sprintf("Test Case %d: Test column positioning for 'Price' after 'ID'", testCaseIndex),
				},
				MetadataModel:        metadataModel,
				Schema:               sch,
				ExpectedOk:           true,
				ExpectedColumnFields: columnFields,
			},
		) {
			return
		}

		priceAfterIDReposition := repositionFieldColumns[0]
		repositionFieldColumns = make(RepositionFieldColumns, 0)
		metadataModel = iter.Map(metadataModel, func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if name, ok := fieldGroup[core.FieldGroupName].(string); ok {
				switch name {
				case "Name":
					fieldGroup[core.FieldColumnPosition] = gojsoncore.JsonObject{
						core.FieldGroupJsonPathKey:    idJsonPathKey,
						core.FieldGroupPositionBefore: true,
					}
					if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
						if fgViewMaxNoOfValuesInSeparateColumns, ok := fieldGroup[core.FieldGroupViewMaxNoOfValuesInSeparateColumns].(int); ok {
							currentIndex := 1
							nextFieldColumnPosition := FieldColumnPosition{
								FieldGroupJsonPathKey:    idJsonPathKey,
								FieldGroupPositionBefore: true,
							}
							for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
								nextFieldColumnPosition.SourceIndex = currentIndex
								repositionFieldColumns = append(repositionFieldColumns, nextFieldColumnPosition)
								nextFieldColumnPosition = FieldColumnPosition{
									FieldGroupJsonPathKey:                       jsonPathKey,
									FieldViewInSeparateColumns:                  true,
									FieldViewValuesInSeparateColumnsHeaderIndex: columnIndex,
								}
								currentIndex++
							}
						}
					}
				}
			}
			return fieldGroup, true
		}).(gojsoncore.JsonObject)
		columnFields = getFieldColumnsFromMetadataModel(metadataModel, sch, nil)
		repositionFieldColumns = append(repositionFieldColumns, priceAfterIDReposition)
		columnFields.RepositionFieldColumns = repositionFieldColumns
		testCaseIndex++
		if !yield(
			&extractionData{
				TestData: internal.TestData{
					TestTitle: fmt.Sprintf("Test Case %d: Test column positioning for 'Name' before 'ID'", testCaseIndex),
				},
				MetadataModel:        metadataModel,
				Schema:               sch,
				ExpectedOk:           true,
				ExpectedColumnFields: columnFields,
			},
		) {
			return
		}
	}

	{
		var addressJsonPathKey path.JSONPath
		metadataModel = iter.Map(testdata.UserProfileMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
				if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "Address" {
					addressJsonPathKey = jsonPathKey
					fieldGroup[core.FieldGroupViewValuesInSeparateColumns] = true
					fieldGroup[core.FieldGroupViewMaxNoOfValuesInSeparateColumns] = 3
					fieldGroup[core.FieldGroupViewDisable] = true
				}
			}
			return fieldGroup, true
		}).(gojsoncore.JsonObject)
		sch = testdata.UserProfileSchema()
		columnFields = getFieldColumnsFromMetadataModel(metadataModel, sch, nil)
		testCaseIndex++
		if !yield(
			&extractionData{
				TestData: internal.TestData{
					TestTitle: fmt.Sprintf("Test Case %d: Test property '%s' for group field that is a group", testCaseIndex, core.FieldGroupViewValuesInSeparateColumns),
				},
				MetadataModel:        metadataModel,
				Schema:               sch,
				ExpectedOk:           true,
				ExpectedColumnFields: columnFields,
			},
		) {
			return
		}
		columnFields = getFieldColumnsFromMetadataModel(metadataModel, sch, core.FieldGroupPropertiesMatch{
			core.FieldGroupViewDisable: core.FuncFieldGroupPropertiesMatcherMatchingProps(func(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) gojsoncore.JsonObject {
				if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
					if strings.HasPrefix(string(jsonPathKey), string(addressJsonPathKey)) {
						return gojsoncore.JsonObject{
							core.FieldGroupViewDisable: true,
						}
					}
				}
				return nil
			}),
		})
		testCaseIndex++
		if !yield(
			&extractionData{
				TestData: internal.TestData{
					TestTitle: fmt.Sprintf("Test Case %d: Test matching props by disabling all fields in Address just by setting '%s' to true on Address alone", testCaseIndex, core.FieldGroupViewDisable),
				},
				MetadataModel: metadataModel,
				Schema:        sch,
				ExpectedOk:    true,
				NestedSkip: core.FieldGroupPropertiesMatch{
					core.FieldGroupViewDisable: core.FuncFieldGroupPropertiesMatcherMatchingProps(func(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) gojsoncore.JsonObject {
						if fieldGroup[core.FieldGroupViewDisable] == true {
							return gojsoncore.JsonObject{
								core.FieldGroupViewDisable: true,
							}
						}
						return nil
					}),
				},
				ExpectedColumnFields: columnFields,
			},
		) {
			return
		}

		repositionFieldColumns = make(RepositionFieldColumns, 0)
		{
			var ageJsonpathKey path.JSONPath
			metadataModel = iter.Map(metadataModel, func(fieldGroup gojsoncore.JsonObject) (any, bool) {
				if name, ok := fieldGroup[core.FieldGroupName].(string); ok {
					if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
						switch name {
						case "Age":
							ageJsonpathKey = jsonPathKey
						case "Address":
							if ageJsonpathKey != "" {
								fieldGroup[core.FieldColumnPosition] = gojsoncore.JsonObject{
									core.FieldGroupJsonPathKey:    ageJsonpathKey,
									core.FieldGroupPositionBefore: true,
								}
								if fgGroupFields, err := core.GetGroupFields(fieldGroup); err == nil {
									if fgGroupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fieldGroup); err == nil {
										if fgViewMaxNoOfValuesInSeparateColumns, ok := fieldGroup[core.FieldGroupViewMaxNoOfValuesInSeparateColumns].(int); ok {
											currentIndex := 2
											nextFieldColumnPosition := FieldColumnPosition{
												FieldGroupJsonPathKey:    ageJsonpathKey,
												FieldGroupPositionBefore: true,
											}
											for columnIndex := range fgViewMaxNoOfValuesInSeparateColumns {
												for _, nFgKeySuffix := range fgGroupReadOrderOfFields {
													if nFieldGroup, err := core.AsJsonObject(fgGroupFields[nFgKeySuffix]); err == nil {
														if nJsonPathKey, err := core.AsJSONPath(nFieldGroup[core.FieldGroupJsonPathKey]); err == nil {
															nextFieldColumnPosition.SourceIndex = currentIndex
															repositionFieldColumns = append(repositionFieldColumns, nextFieldColumnPosition)
															nextFieldColumnPosition = FieldColumnPosition{
																FieldGroupJsonPathKey:                       nJsonPathKey,
																GroupViewInSeparateColumns:                  true,
																GroupViewValuesInSeparateColumnsHeaderIndex: columnIndex,
																GroupViewParentJsonPathKey:                  jsonPathKey,
																FieldJsonPathKeySuffix:                      nFgKeySuffix,
															}
														}
													}
													currentIndex++
												}
											}
										}

									}
								}
							}
						}
					}
				}
				return fieldGroup, true
			}).(gojsoncore.JsonObject)
			columnFields = getFieldColumnsFromMetadataModel(metadataModel, sch, nil)
			columnFields.RepositionFieldColumns = repositionFieldColumns
			testCaseIndex++
			if !yield(
				&extractionData{
					TestData: internal.TestData{
						TestTitle: fmt.Sprintf("Test Case %d: Test column positioning for 'Address' before 'Age'", testCaseIndex),
					},
					MetadataModel:        metadataModel,
					Schema:               sch,
					ExpectedOk:           true,
					ExpectedColumnFields: columnFields,
				},
			) {
				return
			}
		}
	}

	{
		repositionFieldColumns = make(RepositionFieldColumns, 0)
		var skillsJsonPathKey path.JSONPath
		var profileJsonPathKey path.JSONPath
		iter.ForEach(testdata.EmployeeMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
			if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "Skills" {
				if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
					skillsJsonPathKey = jsonPathKey
					return true, true
				}
			}
			return false, false
		})
		currentIndex := 1
		var nextFieldColumnPosition FieldColumnPosition
		metadataModel = iter.Map(testdata.EmployeeMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
				if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "UserProfile" {
					profileJsonPathKey = jsonPathKey
					fieldGroup[core.FieldColumnPosition] = gojsoncore.JsonObject{
						core.FieldGroupJsonPathKey: skillsJsonPathKey,
					}
					nextFieldColumnPosition = FieldColumnPosition{
						FieldGroupJsonPathKey: skillsJsonPathKey,
					}
					return fieldGroup, false
				}

				if profileJsonPathKey != "" && !core.IsFieldAGroup(fieldGroup) && strings.HasPrefix(string(jsonPathKey), string(profileJsonPathKey)) {
					nextFieldColumnPosition.SourceIndex = currentIndex
					repositionFieldColumns = append(repositionFieldColumns, nextFieldColumnPosition)
					nextFieldColumnPosition = FieldColumnPosition{
						FieldGroupJsonPathKey: jsonPathKey,
					}
					currentIndex++
				}
			}

			return fieldGroup, false
		}).(gojsoncore.JsonObject)
		columnFields = getFieldColumnsFromMetadataModel(metadataModel, sch, nil)
		columnFields.RepositionFieldColumns = repositionFieldColumns
		testCaseIndex++
		if !yield(
			&extractionData{
				TestData: internal.TestData{
					TestTitle: fmt.Sprintf("Test Case %d: Test column positioning for nested group 'Profile' after 'Skills'", testCaseIndex),
				},
				MetadataModel:        metadataModel,
				Schema:               sch,
				ExpectedOk:           true,
				ExpectedColumnFields: columnFields,
			},
		) {
			return
		}
	}
}
