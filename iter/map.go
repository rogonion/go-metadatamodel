package iter

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
)

/*
Map recursively modify fields in MetadataModel.

Parameters:
  - group - A metadata model. Should be of type gojsoncore.JsonObject.
  - Callback - Called for each field in a metadata model.
*/
func Map(metadataModelGroup any, callback MapCallback) any {
	fieldGroupProp, err := core.AsJsonObject(metadataModelGroup)
	if err != nil {
		return metadataModelGroup
	}

	groupFields, err := core.GetGroupFields(fieldGroupProp)
	if err != nil {
		return metadataModelGroup
	}

	groupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fieldGroupProp)
	if err != nil {
		return metadataModelGroup
	}

	for _, fgKeySuffix := range groupReadOrderOfFields {
		fgProperty, err := core.AsJsonObject(groupFields[fgKeySuffix])
		if err != nil {
			return metadataModelGroup
		}

		skipFieldGroupPropertyFields := false
		fieldGroupProp[core.GroupFields].(gojsoncore.JsonArray)[0].(gojsoncore.JsonObject)[fgKeySuffix], skipFieldGroupPropertyFields = callback(fgProperty)
		if core.IsFieldAGroup(fgProperty) {
			if !skipFieldGroupPropertyFields {
				fieldGroupProp[core.GroupFields].(gojsoncore.JsonArray)[0].(gojsoncore.JsonObject)[fgKeySuffix] = Map(fgProperty, callback)
			}
		}
	}

	return fieldGroupProp
}

/*
MapCallback called  for each field in a metadata model.

Use to modify fields/group.

Parameters:
  - fieldGroup - current Field/Group property.

Return:
 1. fieldGroup whether modified or not.
 2. `true` to skip processing fieldGroup fields if it contains them.
*/
type MapCallback func(fieldGroup gojsoncore.JsonObject) (any, bool)
