package core

import (
	"reflect"

	gojsoncore "github.com/rogonion/go-json/core"
)

/*
FieldGroupPropertiesMatch A set of gojsoncore.JsonObject properties where the key is the property and the value is the value to match.

If map value is not of type FieldGroupPropertiesMatcher, reflect.DeepEqual will be used to check if property matches.
*/
type FieldGroupPropertiesMatch map[string]any

/*
FieldGroupPropertiesMatcher for complex property matching logic.
*/
type FieldGroupPropertiesMatcher interface {
	// Match
	//
	//Return `true` for property match.
	Match(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) bool
}

type FuncFieldGroupPropertiesMatcher func(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) bool

func (f FuncFieldGroupPropertiesMatcher) Match(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) bool {
	return f(fieldGroupPropertyValue, fieldGroup)
}

/*
IsValid check if FieldGroupPropertiesMatch is empty i.e., is map is not empty.

Returns `true` if map is not empty.
*/
func (n FieldGroupPropertiesMatch) IsValid() bool {
	return len(n) > 0
}

/*
Match returns true to indicate if property in fieldGroup matches entry in FieldGroupPropertiesMatch.

Check if FieldGroupPropertiesMatch IsValid beforehand.
*/
func (n FieldGroupPropertiesMatch) Match(fieldGroup gojsoncore.JsonObject) bool {
	for fieldGroupPropertyKey, valueToMatch := range n {
		fieldGroupProp := fieldGroup[fieldGroupPropertyKey]

		switch matcher := valueToMatch.(type) {
		case FieldGroupPropertiesMatcher:
			if matcher.Match(fieldGroupProp, fieldGroup) {
				return true
			}
		default:
			if reflect.DeepEqual(fieldGroupProp, valueToMatch) {
				return true
			}
		}
	}

	return false
}
