package iter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/internal"
)

func TestIter_Map(t *testing.T) {
	for testData := range mapDataTestData {
		res := Map(testData.MetadataModel, testData.Callback)

		if reflect.DeepEqual(res, testData.Expected) != testData.ExpectedEqual {
			t.Error(
				testData.TestTitle, "\n",
				"expected DeepEqual result to be equal to ExpectedEqual",
				"testData.ExpectedEqual=", testData.ExpectedEqual, "\n",
				"testData.MetadataModel=", gojsoncore.JsonStringifyMust(testData.MetadataModel), "\n",
				"res=", gojsoncore.JsonStringifyMust(res),
			)
		}
	}
}

type mapData struct {
	internal.TestData
	MetadataModel gojsoncore.JsonObject
	Callback      MapCallback
	ExpectedEqual bool
	Expected      any
}

func mapDataTestData(yield func(data *mapData) bool) {
	testCaseIndex := 1
	if !yield(&mapData{
		TestData: internal.TestData{
			TestTitle: fmt.Sprintf("Test Case %d: No Mapping (Keep everything)", testCaseIndex),
		},
		MetadataModel: UserInformationMetadataModel(),
		Callback:      func(fieldGroup gojsoncore.JsonObject) (any, bool) { return fieldGroup, false },
		ExpectedEqual: true,
		Expected:      UserInformationMetadataModel(),
	}) {
		return
	}

	testCaseIndex++
	if !yield(&mapData{
		TestData: internal.TestData{
			TestTitle: fmt.Sprintf("Test Case %d: Append ' Found' to fields whose suffix is 'Name'", testCaseIndex),
		},
		MetadataModel: UserInformationMetadataModel(),
		Callback: func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if fieldGroupName, ok := fieldGroup[core.FieldGroupName].(string); ok {
				if strings.HasSuffix(fieldGroupName, "Name") {
					fieldGroup[core.FieldGroupName] = fieldGroupName + " Found"
				}
			}
			return fieldGroup, false
		},
		ExpectedEqual: true,
		Expected: gojsoncore.JsonObject{
			core.FieldGroupJsonPathKey: path.JsonpathKeyRoot,
			core.FieldGroupName:        "Root Group",
			core.GroupFields: func() gojsoncore.JsonArray {
				FieldGroupJSONPathPrefixDepth0 := path.JsonpathKeyRoot + core.GroupJsonPathPrefix
				return gojsoncore.JsonArray{
					gojsoncore.JsonObject{
						"ID": gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:  FieldGroupJSONPathPrefixDepth0 + "ID",
							core.FieldGroupName:         "Primary ID",
							core.FieldDataType:          core.FieldTypeText,
							core.FieldUI:                core.FieldUiText,
							core.FieldGroupIsPrimaryKey: true,
						},
						"Name": gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey: FieldGroupJSONPathPrefixDepth0 + "Name",
							core.FieldGroupName:        "User Name Found",
							core.FieldDataType:         core.FieldTypeText,
							core.FieldUI:               core.FieldUiText,
						},
						"Details": gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey: FieldGroupJSONPathPrefixDepth0 + "Details",
							core.FieldGroupName:        "Address Details",
							core.GroupFields: func() gojsoncore.JsonArray {
								FieldGroupJSONPathPrefixDepth1 := FieldGroupJSONPathPrefixDepth0 + "Details" + core.GroupJsonPathPrefix
								return gojsoncore.JsonArray{
									gojsoncore.JsonObject{
										"City": gojsoncore.JsonObject{
											core.FieldGroupJsonPathKey: FieldGroupJSONPathPrefixDepth1 + "City",
											core.FieldGroupName:        "City Name Found",
											core.FieldDataType:         core.FieldTypeText,
											core.FieldUI:               core.FieldUiText,
											core.DatabaseJoinDepth:     1,
										},
										"ZipCode": gojsoncore.JsonObject{
											core.FieldGroupJsonPathKey: FieldGroupJSONPathPrefixDepth1 + "ZipCode",
											core.FieldGroupName:        "Postal Code",
											core.FieldDataType:         core.FieldTypeNumber,
											core.FieldUI:               core.FieldUiNumber,
											core.DatabaseJoinDepth:     1,
										},
									},
								}
							}(),
							core.GroupReadOrderOfFields: gojsoncore.JsonArray{"City", "ZipCode"},
						},
					},
				}
			}(),
			core.GroupReadOrderOfFields: gojsoncore.JsonArray{"ID", "Name", "Details"},
		},
	}) {
		return
	}

	testCaseIndex++
	if !yield(&mapData{
		TestData: internal.TestData{
			TestTitle: fmt.Sprintf("Test Case %d: Append ' Found' to fields with nested groups like 'Details'", testCaseIndex),
		},
		MetadataModel: UserInformationMetadataModel(),
		Callback: func(fieldGroup gojsoncore.JsonObject) (any, bool) {
			if core.IsFieldAGroup(fieldGroup) {
				if fieldGroupName, ok := fieldGroup[core.FieldGroupName].(string); ok {
					fieldGroup[core.FieldGroupName] = fieldGroupName + " Found"
				}
			}
			return fieldGroup, false
		},
		ExpectedEqual: true,
		Expected: gojsoncore.JsonObject{
			core.FieldGroupJsonPathKey: path.JsonpathKeyRoot,
			core.FieldGroupName:        "Root Group",
			core.GroupFields: func() gojsoncore.JsonArray {
				FieldGroupJSONPathPrefixDepth0 := path.JsonpathKeyRoot + core.GroupJsonPathPrefix
				return gojsoncore.JsonArray{
					gojsoncore.JsonObject{
						"ID": gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:  FieldGroupJSONPathPrefixDepth0 + "ID",
							core.FieldGroupName:         "Primary ID",
							core.FieldDataType:          core.FieldTypeText,
							core.FieldUI:                core.FieldUiText,
							core.FieldGroupIsPrimaryKey: true,
						},
						"Name": gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey: FieldGroupJSONPathPrefixDepth0 + "Name",
							core.FieldGroupName:        "User Name",
							core.FieldDataType:         core.FieldTypeText,
							core.FieldUI:               core.FieldUiText,
						},
						"Details": gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey: FieldGroupJSONPathPrefixDepth0 + "Details",
							core.FieldGroupName:        "Address Details Found",
							core.GroupFields: func() gojsoncore.JsonArray {
								FieldGroupJSONPathPrefixDepth1 := FieldGroupJSONPathPrefixDepth0 + "Details" + core.GroupJsonPathPrefix
								return gojsoncore.JsonArray{
									gojsoncore.JsonObject{
										"City": gojsoncore.JsonObject{
											core.FieldGroupJsonPathKey: FieldGroupJSONPathPrefixDepth1 + "City",
											core.FieldGroupName:        "City Name",
											core.FieldDataType:         core.FieldTypeText,
											core.FieldUI:               core.FieldUiText,
											core.DatabaseJoinDepth:     1,
										},
										"ZipCode": gojsoncore.JsonObject{
											core.FieldGroupJsonPathKey: FieldGroupJSONPathPrefixDepth1 + "ZipCode",
											core.FieldGroupName:        "Postal Code",
											core.FieldDataType:         core.FieldTypeNumber,
											core.FieldUI:               core.FieldUiNumber,
											core.DatabaseJoinDepth:     1,
										},
									},
								}
							}(),
							core.GroupReadOrderOfFields: gojsoncore.JsonArray{"City", "ZipCode"},
						},
					},
				}
			}(),
			core.GroupReadOrderOfFields: gojsoncore.JsonArray{"ID", "Name", "Details"},
		},
	}) {
		return
	}
}
