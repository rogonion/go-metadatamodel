package core

import (
	"fmt"
	"strings"

	"github.com/rogonion/go-json/path"
)

/*
Get

Parameters:

  - jsonPathKey - From property FieldGroupJsonPathKey in metadata model field/group.

  - arrayIndexes - Set of actual indexes of arrays or slice to replace ArrayPathPlaceholder in path.JSONPath in order.

    If empty or uninitialized, it will default to `0`.
*/
func (n *JsonPathToValue) Get(jsonPathKey path.JSONPath, arrayIndexes []int) (path.JSONPath, error) {
	jsonPathKeyStr := string(jsonPathKey)
	noOfArrayPathPlaceholders := len(ArrayPathRegexSearch().FindAllStringSubmatch(jsonPathKeyStr, -1))

	if n.removeGroupFields {
		if n.sourceOfValueIsAnArray {
			jsonPathKeyStr = strings.Replace(jsonPathKeyStr, path.JsonpathDotNotation+GroupFields, "", 1)
		} else {
			jsonPathKeyStr = strings.Replace(jsonPathKeyStr, path.JsonpathDotNotation+GroupFields+ArrayPathPlaceholder, "", 1)
		}
		jsonPathKeyStr = string(GroupFieldsRegexSearch().ReplaceAll([]byte(jsonPathKeyStr), []byte("")))
	}

	if n.replaceArrayPathPlaceholderWithActualIndexes || len(arrayIndexes) > 0 {
		if len(arrayIndexes) == 0 {
			arrayIndexes = make([]int, noOfArrayPathPlaceholders)
		}

		for _, index := range arrayIndexes {
			jsonPathKeyStr = strings.Replace(jsonPathKeyStr, ArrayPathPlaceholder, fmt.Sprintf("%s%d%s", path.JsonpathLeftBracket, index, path.JsonpathRightBracket), 1)
			arrayIndexes = arrayIndexes[1:]
		}

		if strings.Contains(jsonPathKeyStr, ArrayPathPlaceholder) {
			return path.JSONPath(jsonPathKeyStr), ErrPathContainsIndexPlaceholders
		}
	}

	return path.JSONPath(jsonPathKeyStr), nil
}

// WithReplaceArrayPathPlaceholderWithActualIndexes sets whether to replace the array path placeholder with actual indexes.
func (n *JsonPathToValue) WithReplaceArrayPathPlaceholderWithActualIndexes(value bool) *JsonPathToValue {
	n.SetReplaceArrayPathPlaceholderWithActualIndexes(value)
	return n
}

// SetReplaceArrayPathPlaceholderWithActualIndexes sets whether to replace the array path placeholder with actual indexes.
func (n *JsonPathToValue) SetReplaceArrayPathPlaceholderWithActualIndexes(value bool) {
	n.replaceArrayPathPlaceholderWithActualIndexes = value
}

// WithSourceOfValueIsAnArray sets whether the source of the value is an array.
func (n *JsonPathToValue) WithSourceOfValueIsAnArray(value bool) *JsonPathToValue {
	n.SetSourceOfValueIsAnArray(value)
	return n
}

// SetSourceOfValueIsAnArray sets whether the source of the value is an array.
func (n *JsonPathToValue) SetSourceOfValueIsAnArray(value bool) {
	n.sourceOfValueIsAnArray = value
}

// WithRemoveGroupFields sets whether to remove group fields from the path.
func (n *JsonPathToValue) WithRemoveGroupFields(value bool) *JsonPathToValue {
	n.SetRemoveGroupFields(value)
	return n
}

// SetRemoveGroupFields sets whether to remove group fields from the path.
func (n *JsonPathToValue) SetRemoveGroupFields(value bool) {
	n.removeGroupFields = value
}

/*
NewJsonPathToValue

- Set JsonPathToValue.replaceArrayPathPlaceholderWithActualIndexes to `true`.

- Set JsonPathToValue.removeGroupFields to `true`.
*/
func NewJsonPathToValue() *JsonPathToValue {
	n := new(JsonPathToValue)
	n.SetReplaceArrayPathPlaceholderWithActualIndexes(true)
	n.SetRemoveGroupFields(true)
	return n
}

/*
JsonPathToValue retrieves the actual path.JSONPath to target value in a source object.

It can do the following:
  - Remove GroupFields from path.JSONPath.
  - Replace ArrayPathPlaceholder with actual array index.

Usage:
  - Instantiate the JsonPathToValue using NewJsonPathToValue.
  - Retrieve path.JSONPath using JsonPathToValue.Get.
*/
type JsonPathToValue struct {
	// Remove GroupFields from path.JSONPath.
	removeGroupFields bool

	// If source of value is NOT an array or slice, the first pair of GroupJsonPathPrefix is removed as the source is assumed to be an associative collection.
	sourceOfValueIsAnArray bool

	// Replace ArrayPathPlaceholder wildcard with an actual array index.
	replaceArrayPathPlaceholderWithActualIndexes bool
}
