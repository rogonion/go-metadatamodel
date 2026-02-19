package fieldcolumns

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/internal"
	"github.com/rogonion/go-metadatamodel/iter"
	"github.com/rogonion/go-metadatamodel/testdata"
)

func TestFieldColumns_Skip(t *testing.T) {
	for testData := range skipTestData {
		fcExtraction := NewColumnFieldsExtraction(testData.MetadataModel).WithSkip(testData.NestedSkip).WithAdd(testData.NestedAdd)
		columnFields, err := fcExtraction.Extract()
		if testData.ExpectedOk && err != nil {
			t.Error(
				testData.TestTitle, "\n",
				"expected extraction ok=", testData.ExpectedOk, "got error=", err, "\n",
			)
		}

		if columnFields == nil {
			continue
		}

		columnFields.Skip(testData.Skip, testData.Add)
		if !reflect.DeepEqual(columnFields.FieldsToSkip, testData.ExpectedFieldsToSkip) {
			t.Error(
				testData.TestTitle, "\n",
				"result skip ColumnFields not equal to testData.ExpectedFieldsToSkip\n",
				"ExpectedColumnFields=", gojsoncore.JsonStringifyMust(testData.ExpectedFieldsToSkip), "\n",
				"result=", gojsoncore.JsonStringifyMust(columnFields.FieldsToSkip),
			)
		}
	}
}

type skipData struct {
	internal.TestData
	MetadataModel        gojsoncore.JsonObject
	ExpectedOk           bool
	Skip                 core.FieldGroupPropertiesMatch
	NestedSkip           core.FieldGroupPropertiesMatch
	Add                  core.FieldGroupPropertiesMatch
	NestedAdd            core.FieldGroupPropertiesMatch
	ExpectedFieldsToSkip FieldsToSkip
}

func skipTestData(yield func(data *skipData) bool) {
	expectedFieldsToSkip := make(FieldsToSkip)
	metadataModel := iter.Map(testdata.UserProfileMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
		if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
			if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "Age" {
				fieldGroup[core.FieldGroupViewDisable] = true
				expectedFieldsToSkip[jsonPathKey] = FieldToSkip()
			}
		}
		return fieldGroup, true
	}).(gojsoncore.JsonObject)
	testCaseIndex := 1
	if !yield(
		&skipData{
			TestData: internal.TestData{
				TestTitle: "Skip 'Age' field in UserProfile",
			},
			MetadataModel: metadataModel,
			ExpectedOk:    true,
			Skip: core.FieldGroupPropertiesMatch{
				core.FieldGroupViewDisable: true,
			},
			ExpectedFieldsToSkip: expectedFieldsToSkip,
		},
	) {
		return
	}

	expectedFieldsToSkip = make(FieldsToSkip)
	metadataModel = iter.Map(testdata.UserProfileMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
		if fieldGroupJsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
			if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "Address" {
				fieldGroup[core.FieldGroupViewValuesInSeparateColumns] = true
				fieldGroup[core.FieldGroupViewMaxNoOfValuesInSeparateColumns] = 3
				fieldGroup[core.FieldGroupViewDisable] = true
				if fgReadOrder, err := core.GetGroupReadOrderOfFields(fieldGroup); err == nil {
					if fgFields, err := core.GetGroupFields(fieldGroup); err == nil {
						for currentIndex := range 3 {
							for _, fgKeySuffix := range fgReadOrder {
								if field, err := core.AsJsonObject(fgFields[fgKeySuffix]); err == nil {
									if _, err := core.AsJSONPath(field[core.FieldGroupJsonPathKey]); err == nil {
										expectedFieldsToSkip[path.JSONPath(string(fieldGroupJsonPathKey)+path.JsonpathDotNotation+core.GroupFields+path.JsonpathLeftBracket+fmt.Sprintf("%d", currentIndex)+path.JsonpathRightBracket+path.JsonpathDotNotation+fgKeySuffix)] = FieldToSkip()
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
	testCaseIndex++
	if !yield(
		&skipData{
			TestData: internal.TestData{
				TestTitle: "Skip Pivoted Address Fields in UserProfile",
			},
			MetadataModel: metadataModel,
			ExpectedOk:    true,
			Skip: core.FieldGroupPropertiesMatch{
				core.FieldGroupViewDisable: true,
			},
			NestedSkip: core.FieldGroupPropertiesMatch{
				core.FieldGroupViewDisable: core.FuncFieldGroupPropertiesMatcherMatchingProps(func(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) gojsoncore.JsonObject {
					if viewDisable, ok := fieldGroup[core.FieldGroupViewDisable].(bool); ok && viewDisable {
						return gojsoncore.JsonObject{
							core.FieldGroupViewDisable: true,
						}
					}
					return nil
				}),
			},
			ExpectedFieldsToSkip: expectedFieldsToSkip,
		},
	) {
		return
	}

	{
		var profileJsonPathKey path.JSONPath
		expectedFieldsToSkip = make(FieldsToSkip)
		metadataModel = iter.Map(testdata.EmployeeMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
				if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "UserProfile" {
					profileJsonPathKey = jsonPathKey
					fieldGroup[core.FieldGroupViewDisable] = true
					return fieldGroup, false
				}

				if profileJsonPathKey != "" && !core.IsFieldAGroup(fieldGroup) && strings.HasPrefix(string(jsonPathKey), string(profileJsonPathKey)) {
					expectedFieldsToSkip[jsonPathKey] = FieldToSkip()
				}
			}

			return fieldGroup, false
		}).(gojsoncore.JsonObject)
		testCaseIndex++
		if !yield(
			&skipData{
				TestData: internal.TestData{
					TestTitle: "Skip Nested Profile Fields in Employee",
				},
				MetadataModel: metadataModel,
				ExpectedOk:    true,
				Skip: core.FieldGroupPropertiesMatch{
					core.FieldGroupViewDisable: true,
				},
				NestedSkip: core.FieldGroupPropertiesMatch{
					core.FieldGroupViewDisable: core.FuncFieldGroupPropertiesMatcherMatchingProps(func(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) gojsoncore.JsonObject {
						if viewDisable, ok := fieldGroup[core.FieldGroupViewDisable].(bool); ok && viewDisable {
							return gojsoncore.JsonObject{
								core.FieldGroupViewDisable: true,
							}
						}
						return nil
					}),
				},
				ExpectedFieldsToSkip: expectedFieldsToSkip,
			},
		) {
			return
		}
	}
}

func TestFieldColumns_Reposition(t *testing.T) {
	for testData := range repositionTestData {
		fcExtraction := NewColumnFieldsExtraction(testData.MetadataModel)
		columnFields, err := fcExtraction.Extract()
		if testData.ExpectedOk && err != nil {
			t.Error(
				testData.TestTitle, "\n",
				"expected extraction ok=", testData.ExpectedOk, "got error=", err, "\n",
			)
		}

		if columnFields == nil {
			continue
		}

		columnFields.Reposition()
		if !reflect.DeepEqual(columnFields.RepositionedReadOrderOfColumnFields, testData.ExpectedIndexOfReadOrderOfColumnFields) {
			t.Error(
				testData.TestTitle, "\n",
				"result repositioned ColumnFields not equal to testData.ExpectedIndexOfReadOrderOfColumnFields\n",
				"ExpectedColumnFields=", gojsoncore.JsonStringifyMust(testData.ExpectedIndexOfReadOrderOfColumnFields), "\n",
				"result=", gojsoncore.JsonStringifyMust(columnFields.RepositionedReadOrderOfColumnFields),
			)
		}
	}
}

type repositionData struct {
	internal.TestData
	MetadataModel                          gojsoncore.JsonObject
	ExpectedOk                             bool
	ExpectedIndexOfReadOrderOfColumnFields []int
}

func repositionTestData(yield func(data *repositionData) bool) {
	metadataModel := testdata.UserMetadataModel(nil)
	testCaseIndex := 1
	if !yield(
		&repositionData{
			TestData: internal.TestData{
				TestTitle: "User Metadata Model - Default Order",
			},
			MetadataModel:                          metadataModel,
			ExpectedOk:                             true,
			ExpectedIndexOfReadOrderOfColumnFields: []int{0, 1, 2},
		},
	) {
		return
	}

	metadataModel = testdata.EmployeeMetadataModel(nil)
	testCaseIndex++
	if !yield(
		&repositionData{
			TestData: internal.TestData{
				TestTitle: "Employee Metadata Model - Default Order",
			},
			MetadataModel:                          metadataModel,
			ExpectedOk:                             true,
			ExpectedIndexOfReadOrderOfColumnFields: []int{0, 1, 2, 3, 4, 5, 6},
		},
	) {
		return
	}

	{
		metadataModel = iter.Map(testdata.UserProfileMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if _, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
				if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "Address" {
					fieldGroup[core.FieldGroupViewValuesInSeparateColumns] = true
					fieldGroup[core.FieldGroupViewMaxNoOfValuesInSeparateColumns] = 3
					fieldGroup[core.FieldGroupViewDisable] = true
				}
			}
			return fieldGroup, true
		}).(gojsoncore.JsonObject)
		testCaseIndex++
		if !yield(
			&repositionData{
				TestData: internal.TestData{
					TestTitle: fmt.Sprintf("UserProfile - Pivoted Address ('%s')", core.FieldGroupViewValuesInSeparateColumns),
				},
				MetadataModel:                          metadataModel,
				ExpectedOk:                             true,
				ExpectedIndexOfReadOrderOfColumnFields: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
		) {
			return
		}
	}

	metadataModel = iter.Map(testdata.ProductMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
		if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "Name" {
			fieldGroup[core.FieldGroupViewValuesInSeparateColumns] = true
			fieldGroup[core.FieldGroupViewMaxNoOfValuesInSeparateColumns] = 3
		}
		return fieldGroup, true
	}).(gojsoncore.JsonObject)
	testCaseIndex++
	if !yield(
		&repositionData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Product - Pivoted Name ('%s')", core.FieldGroupViewValuesInSeparateColumns),
			},
			MetadataModel:                          metadataModel,
			ExpectedOk:                             true,
			ExpectedIndexOfReadOrderOfColumnFields: []int{0, 1, 2, 3, 4},
		},
	) {
		return
	}

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
					}
				}
			}
			return fieldGroup, true
		}).(gojsoncore.JsonObject)
		testCaseIndex++
		if !yield(
			&repositionData{
				TestData: internal.TestData{
					TestTitle: "Product - Reposition 'Price' after 'ID'",
				},
				MetadataModel:                          metadataModel,
				ExpectedOk:                             true,
				ExpectedIndexOfReadOrderOfColumnFields: []int{0, 4, 1, 2, 3},
			},
		) {
			return
		}

		metadataModel = iter.Map(metadataModel, func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if name, ok := fieldGroup[core.FieldGroupName].(string); ok {
				switch name {
				case "Name":
					fieldGroup[core.FieldColumnPosition] = gojsoncore.JsonObject{
						core.FieldGroupJsonPathKey:    idJsonPathKey,
						core.FieldGroupPositionBefore: true,
					}
				}
			}
			return fieldGroup, true
		}).(gojsoncore.JsonObject)
		testCaseIndex++
		if !yield(
			&repositionData{
				TestData: internal.TestData{
					TestTitle: "Product - Reposition 'Name' before 'ID'",
				},
				MetadataModel:                          metadataModel,
				ExpectedOk:                             true,
				ExpectedIndexOfReadOrderOfColumnFields: []int{1, 2, 3, 0, 4},
			},
		) {
			return
		}
	}

	{
		var ageJsonpathKey path.JSONPath
		metadataModel = iter.Map(testdata.UserProfileMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if name, ok := fieldGroup[core.FieldGroupName].(string); ok {
				switch name {
				case "Age":
					if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
						ageJsonpathKey = jsonPathKey
					}
				case "Address":
					fieldGroup[core.FieldGroupViewValuesInSeparateColumns] = true
					fieldGroup[core.FieldGroupViewMaxNoOfValuesInSeparateColumns] = 3
					fieldGroup[core.FieldGroupViewDisable] = true
					if ageJsonpathKey != "" {
						fieldGroup[core.FieldColumnPosition] = gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:    ageJsonpathKey,
							core.FieldGroupPositionBefore: true,
						}
					}
				}
			}
			return fieldGroup, true
		}).(gojsoncore.JsonObject)
		testCaseIndex++
		if !yield(
			&repositionData{
				TestData: internal.TestData{
					TestTitle: "UserProfile - Reposition 'Address' before 'Age'",
				},
				MetadataModel:                          metadataModel,
				ExpectedOk:                             true,
				ExpectedIndexOfReadOrderOfColumnFields: []int{0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1},
			},
		) {
			return
		}
	}

	{
		var skillsJsonPathKey path.JSONPath
		iter.ForEach(testdata.EmployeeMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
			if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "Skills" {
				if jsonPathKey, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
					skillsJsonPathKey = jsonPathKey
					return true, true
				}
			}
			return false, false
		})
		metadataModel = iter.Map(testdata.EmployeeMetadataModel(nil), func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if _, err := core.AsJSONPath(fieldGroup[core.FieldGroupJsonPathKey]); err == nil {
				if name, ok := fieldGroup[core.FieldGroupName].(string); ok && name == "UserProfile" {
					fieldGroup[core.FieldColumnPosition] = gojsoncore.JsonObject{
						core.FieldGroupJsonPathKey: skillsJsonPathKey,
					}
					return fieldGroup, true
				}
			}

			return fieldGroup, false
		}).(gojsoncore.JsonObject)
		testCaseIndex++
		if !yield(
			&repositionData{
				TestData: internal.TestData{
					TestTitle: "Employee - Reposition Nested 'Profile' after 'Skills'",
				},
				MetadataModel:                          metadataModel,
				ExpectedOk:                             true,
				ExpectedIndexOfReadOrderOfColumnFields: []int{0, 6, 1, 2, 3, 4, 5},
			},
		) {
			return
		}
	}
}
