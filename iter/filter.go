package iter

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
)

/*
Filter recursively remove fields in MetadataModel.

Parameters:
  - group - metadata model. Should be of type gojsoncore.JsonObject.
  - Callback - called for each field in a metadata model.

Returns modified group.
*/
func Filter(group any, callback FilterCallback) any {
	fieldGroupProp, err := core.AsJsonObject(group)
	if err != nil {
		return group
	}

	groupFields, err := core.GetGroupFields(fieldGroupProp)
	if err != nil {
		return group
	}

	groupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fieldGroupProp)
	if err != nil {
		return group
	}

	for fgKeySuffixIndex, fgKeySuffix := range groupReadOrderOfFields {
		fgProperty, err := core.AsJsonObject(groupFields[fgKeySuffix])
		if err != nil {
			return group
		}

		retainFieldGroup, skipFieldGroupPropertyFields := callback(fgProperty)
		if !retainFieldGroup {
			fieldGroupProp[core.GroupReadOrderOfFields].(gojsoncore.JsonArray)[fgKeySuffixIndex] = ""
			delete(fieldGroupProp[core.GroupFields].(gojsoncore.JsonArray)[0].(gojsoncore.JsonObject), fgKeySuffix)
		} else {
			if core.IsFieldAGroup(fgProperty) {
				if !skipFieldGroupPropertyFields {
					fieldGroupProp[core.GroupFields].(gojsoncore.JsonArray)[0].(gojsoncore.JsonObject)[fgKeySuffix] = Filter(fgProperty, callback)
				}
			}
		}
	}

	newGroupReadOrderOfFields := make(gojsoncore.JsonArray, 0)
	for _, fgJsonPathKeySuffix := range fieldGroupProp[core.GroupReadOrderOfFields].(gojsoncore.JsonArray) {
		if fgJsonPathKeySuffixString, ok := fgJsonPathKeySuffix.(string); ok && len(fgJsonPathKeySuffixString) > 0 {
			newGroupReadOrderOfFields = append(newGroupReadOrderOfFields, fgJsonPathKeySuffixString)
		}
	}
	fieldGroupProp[core.GroupReadOrderOfFields] = newGroupReadOrderOfFields

	return fieldGroupProp
}

/*
FilterCallback called  for each field in a metadata model

Parameters:
  - fieldGroup - current Field/Group property.

Returns:
 1. `false` to signal fieldGroup should be removed.
 2. `true` to skip processing fieldGroup fields if it contains them.
*/
type FilterCallback func(fieldGroup gojsoncore.JsonObject) (bool, bool)
