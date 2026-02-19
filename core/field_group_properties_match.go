package core

import (
	"reflect"

	gojsoncore "github.com/rogonion/go-json/core"
)

/*
FieldGroupPropertiesMatch A set of gojsoncore.JsonObject properties where the key is the property and the value is the value to match.

If map value is not of type FieldGroupPropertiesFirstMatcher or FieldGroupPropertiesMatchingProps, reflect.DeepEqual will be used to check if property matches.
*/
type FieldGroupPropertiesMatch map[string]any

/*
IsValid check if FieldGroupPropertiesMatch is empty i.e., is map is not empty.

Returns `true` if map is not empty.
*/
func (n FieldGroupPropertiesMatch) IsValid() bool {
	return len(n) > 0
}

/*
FieldGroupPropertiesFirstMatcher for complex property matching logic.

Use for simple first match.
*/
type FieldGroupPropertiesFirstMatcher interface {
	// FirstMatch
	//
	// Return `true` for first property match.
	FirstMatch(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) bool
}

// FuncFieldGroupPropertiesMatcherFirst is a function type that implements FieldGroupPropertiesFirstMatcher.
type FuncFieldGroupPropertiesMatcherFirst func(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) bool

// FirstMatch calls the underlying function to check for a match.
func (f FuncFieldGroupPropertiesMatcherFirst) FirstMatch(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) bool {
	return f(fieldGroupPropertyValue, fieldGroup)
}

/*
FirstMatch returns true to indicate if property in fieldGroup matches entry in FieldGroupPropertiesMatch.

Check if FieldGroupPropertiesMatch IsValid beforehand.
*/
func (n FieldGroupPropertiesMatch) FirstMatch(fieldGroup gojsoncore.JsonObject) bool {
	for fieldGroupPropertyKey, valueToMatch := range n {
		fieldGroupProp := fieldGroup[fieldGroupPropertyKey]

		switch matcher := valueToMatch.(type) {
		case FieldGroupPropertiesFirstMatcher:
			if matcher.FirstMatch(fieldGroupProp, fieldGroup) {
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

/*
FieldGroupPropertiesMatchingProps for complex property matching logic.

Use if you want to retrieve the set of props that satisfied the match.
*/
type FieldGroupPropertiesMatchingProps interface {
	// MatchingProps
	//
	// Return set of properties that match.
	MatchingProps(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) gojsoncore.JsonObject
}

// FuncFieldGroupPropertiesMatcherMatchingProps is a function type that implements FieldGroupPropertiesMatchingProps.
type FuncFieldGroupPropertiesMatcherMatchingProps func(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) gojsoncore.JsonObject

// MatchingProps calls the underlying function to retrieve matching properties.
func (f FuncFieldGroupPropertiesMatcherMatchingProps) MatchingProps(fieldGroupPropertyValue any, fieldGroup gojsoncore.JsonObject) gojsoncore.JsonObject {
	return f(fieldGroupPropertyValue, fieldGroup)
}

/*
MatchingProps returns the set of properties in fieldGroup that satisfied the match.

Check if FieldGroupPropertiesMatch IsValid beforehand.
*/
func (n FieldGroupPropertiesMatch) MatchingProps(fieldGroup gojsoncore.JsonObject) gojsoncore.JsonObject {
	matchingProps := make(gojsoncore.JsonObject)

	for fieldGroupPropertyKey, valueToMatch := range n {
		fieldGroupProp := fieldGroup[fieldGroupPropertyKey]

		switch matcher := valueToMatch.(type) {
		case FieldGroupPropertiesMatchingProps:
			if res := matcher.MatchingProps(fieldGroupProp, fieldGroup); len(res) > 0 {
				for k, v := range res {
					matchingProps[k] = v
				}
			}
		default:
			if reflect.DeepEqual(fieldGroupProp, valueToMatch) {
				matchingProps[fieldGroupPropertyKey] = fieldGroupProp
			}
		}
	}

	return matchingProps
}
