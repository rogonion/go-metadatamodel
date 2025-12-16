package iter

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
)

/*
ForEach recursively loop through fields in MetadataModel.

Parameters:
  - group - metadata model. Should be of type gojsoncore.JsonObject.
  - Callback - called for each field in a metadata model.
*/
func ForEach(group any, callback ForeachCallback) {
	fieldGroupProp, err := core.AsJsonObject(group)
	if err != nil {
		return
	}

	groupFields, err := core.GetGroupFields(fieldGroupProp)
	if err != nil {
		return
	}

	groupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fieldGroupProp)
	if err != nil {
		return
	}

	for _, fgKeySuffix := range groupReadOrderOfFields {
		fgProperty, err := core.AsJsonObject(groupFields[fgKeySuffix])
		if err != nil {
			return
		}

		terminateLoop, skipFieldGroupPropertyFields := callback(fgProperty)
		if terminateLoop {
			return
		}

		if core.IsFieldAGroup(fgProperty) {
			if !skipFieldGroupPropertyFields {
				ForEach(fgProperty, callback)
			}
		}
	}
}

/*
ForeachCallback called  for each field in a metadata model

Parameters:
  - fieldGroup - current Field/Group property.

Return:
 1. `true` to signal loop should be terminated.
 2. `true` to skip processing fieldGroup fields if it contains them.
*/
type ForeachCallback func(fieldGroup gojsoncore.JsonObject) (bool, bool)
