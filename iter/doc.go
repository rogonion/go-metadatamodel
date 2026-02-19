/*
Package iter provides higher order functions like Filter, ForEach, and Map against a metadata model.

# Usage

## Filter

Filter recursively removes fields in a MetadataModel based on a callback.

	import (
		"strings"
		gojsoncore "github.com/rogonion/go-json/core"
		"github.com/rogonion/go-metadatamodel/core"
		"github.com/rogonion/go-metadatamodel/iter"
		"github.com/rogonion/go-metadatamodel/testdata"
	)

	var sourceMetadataModel any = testdata.AddressMetadataModel(nil)

	var updatedMetadataModel any = iter.Filter(sourceMetadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
		if fieldGroupName, ok := fieldGroup[core.FieldGroupName].(string); ok {
			if strings.HasSuffix(fieldGroupName, "Name") {
				return false, false
			}
		}
		return true, false
	})

## Map

Map recursively modifies fields in a MetadataModel.

	var sourceMetadataModel any = testdata.AddressMetadataModel(nil)

	var updatedMetadataModel any = iter.Map(sourceMetadataModel, func(fieldGroup gojsoncore.JsonObject) (any, bool) {
		if fieldGroupName, ok := fieldGroup[core.FieldGroupName].(string); ok {
			if strings.HasSuffix(fieldGroupName, "Code") {
				fieldGroup[core.FieldGroupName] = fieldGroupName + " Found"
			}
		}
		return fieldGroup, false
	})

## ForEach

ForEach recursively loops through fields in a MetadataModel.

	var sourceMetadataModel any = testdata.AddressMetadataModel(nil)
	var noOfIterations uint64 = 0

	iter.ForEach(sourceMetadataModel, func (fieldGroup gojsoncore.JsonObject)(bool, bool) {
		noOfIterations++
		return false, false
	})
*/
package iter
