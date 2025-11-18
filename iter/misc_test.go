package iter

import (
	"github.com/brunoga/deep"
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
)

func UserInformationMetadataModel() gojsoncore.JsonObject {
	return deep.MustCopy(gojsoncore.JsonObject{
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
						core.FieldGroupName:        "Address Details",
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
	})
}
