package iter

import (
	"reflect"
	"strings"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/internal"
)

func TestIter_Filter(t *testing.T) {
	for testData := range filterDataTestData {
		res := Filter(testData.MetadataModel, testData.Callback)

		if reflect.DeepEqual(res, testData.Expected) != testData.ExpectedEqual {
			t.Error(
				"Test Title:", testData.TestTitle, "\n",
				"expected DeepEqual result to be equal to ExpectedEqual",
				"testData.ExpectedEqual=", testData.ExpectedEqual, "\n",
				"testData.MetadataModel=", gojsoncore.JsonStringifyMust(testData.MetadataModel), "\n",
				"res=", gojsoncore.JsonStringifyMust(res),
			)
		}
	}
}

type filterData struct {
	internal.TestData
	MetadataModel gojsoncore.JsonObject
	Callback      FilterCallback
	ExpectedEqual bool
	Expected      any
}

func filterDataTestData(yield func(data *filterData) bool) {
	if !yield(&filterData{
		TestData: internal.TestData{
			TestTitle: "Case 1: No filter (Keep everything)",
		},
		MetadataModel: UserInformationMetadataModel(),
		Callback:      func(fieldGroup gojsoncore.JsonObject) (bool, bool) { return true, false },
		ExpectedEqual: true,
		Expected:      UserInformationMetadataModel(),
	}) {
		return
	}

	if !yield(&filterData{
		TestData: internal.TestData{
			TestTitle: "Case 2: Remove fields whose suffix is 'Name'",
		},
		MetadataModel: UserInformationMetadataModel(),
		Callback: func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
			if fieldGroupName, ok := fieldGroup[core.FieldGroupName].(string); ok {
				if strings.HasSuffix(fieldGroupName, "Name") {
					return false, false
				}
			}
			return true, false
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
						"Details": gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey: FieldGroupJSONPathPrefixDepth0 + "Details",
							core.FieldGroupName:        "Address Details",
							core.GroupFields: func() gojsoncore.JsonArray {
								FieldGroupJSONPathPrefixDepth1 := FieldGroupJSONPathPrefixDepth0 + "Details" + core.GroupJsonPathPrefix
								return gojsoncore.JsonArray{
									gojsoncore.JsonObject{
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
							core.GroupReadOrderOfFields: gojsoncore.JsonArray{"ZipCode"},
						},
					},
				}
			}(),
			core.GroupReadOrderOfFields: gojsoncore.JsonArray{"ID", "Details"},
		},
	}) {
		return
	}

	if !yield(&filterData{
		TestData: internal.TestData{
			TestTitle: "Case 3: Remove an fields with nested groups like (\"Details\")",
		},
		MetadataModel: UserInformationMetadataModel(),
		Callback: func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
			if core.IsFieldAGroup(fieldGroup) {
				return false, false
			}
			return true, false
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
					},
				}
			}(),
			core.GroupReadOrderOfFields: gojsoncore.JsonArray{"ID", "Name"},
		},
	}) {
		return
	}
}
