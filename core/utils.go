package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
)

// PreparePathToValueInObject Prepares the path to value in an object based on the metadata-model FieldGroupJsonPathKey property of a field in a group.
//
// Parameters:
//
//   - path - path to value in object. Must begin with `$.GroupFields[*]`.
//     Examples: `$.GroupFields[*].field_1` results in `field_1` and `$.GroupFields[*].group_1.GroupFields[*].group_1_field` results in `group_1[*].group_1_field`.
//
//   - groupIndexes - Each element replaces array index placeholder (ArrayPathPlaceholder) found in path.
//
//     Must NOT be empty as the first element in groupIndexes removed as it matches the first `GroupFields[*]` in the path which is removed from the path since it indicates the root of the metadata-model.
//
//     Number of elements MUST match number of array index placeholders in path.
//
//     For example, with path like `$.GroupFields[*].group_1.GroupFields[*].group_1_field` the number of array indexes passed in groupIndexes MUST be 2.
//
// The first element in groupIndexes removed as it matches the first `GroupFields[*]` in the path which is removed from the path since it indicates the root of the metadata-model.
//
// For example, path `$.GroupFields[*].group_1.GroupFields[*].group_1_field` will be trimmed to `$.group_1[*].group_1_field` before groupIndexes are added.
//
// Returns path to value in object or error if the number of array index placeholders in path being more than the number of array indexes in groupIndexes.
func PreparePathToValueInObject(path string, groupIndexes []int) (string, error) {
	path = strings.Replace(path, fmt.Sprintf(".%s[*]", GroupFields), "", 1)
	path = string(GroupFieldsRegexSearch().ReplaceAll([]byte(path), []byte("")))
	groupIndexes = groupIndexes[1:]
	for _, groupIndex := range groupIndexes {
		path = strings.Replace(path, ArrayPathPlaceholder, fmt.Sprintf("[%v]", groupIndex), 1)
		groupIndexes = groupIndexes[1:]
	}

	if strings.Contains(path, ArrayPathPlaceholder) {
		return path, ErrPathContainsIndexPlaceholders
	}

	return path, nil
}

// FgGet2DConversion Retrieves the maximum no of columns/sets of columns that can be used to represent a field/group in a flat 2D table.
//
// Conditions for successful conversion:
//   - Must be a field not of type any.
//   - If field is a group, its group structure must not contain nested groups.
//
// Parameters:
//   - fg - field to get 2D properties.
//
// Returns the maximum number. If value is less than 0, it means the field/group cannot be represented.
func FgGet2DConversion(fg any) int {
	fgViewMaxNoOfValuesInSeparateColumns := -1
	if fgProperty, ok := fg.(gojsoncore.JsonObject); ok {
		if value, ok := fgProperty[FieldGroupViewMaxNoOfValuesInSeparateColumns].(bool); ok && value {
			if groupFieldsArray, ok := fgProperty[GroupFields].([]any); ok && len(groupFieldsArray) > 0 {
				if groupFieldsMap, ok := groupFieldsArray[0].(gojsoncore.JsonObject); ok {
					for _, value := range groupFieldsMap {
						if valueMap, ok := value.(gojsoncore.JsonObject); ok {
							if valueMap[GroupReadOrderOfFields] != nil && reflect.TypeOf(valueMap[GroupReadOrderOfFields]).Kind() == reflect.Slice {
								return fgViewMaxNoOfValuesInSeparateColumns
							}
						}
					}
				}
			}
			if vInt, ok := fgProperty[FieldGroupViewMaxNoOfValuesInSeparateColumns].(int); ok && vInt > 1 {
				fgViewMaxNoOfValuesInSeparateColumns = vInt
			} else {
				if vFloat, ok := fgProperty[FieldGroupViewMaxNoOfValuesInSeparateColumns].(float64); ok && vFloat > 1 {
					fgViewMaxNoOfValuesInSeparateColumns = int(vFloat)
				}
			}
		}
	}

	return fgViewMaxNoOfValuesInSeparateColumns
}

func IfKeySuffixMatchesValues(keyToCheck string, valuesToMatch []string) bool {
	for _, value := range valuesToMatch {
		if strings.HasSuffix(keyToCheck, value) {
			return true
		}
	}

	return false
}

func IsFieldAField(f any) bool {
	if field, ok := f.(gojsoncore.JsonObject); ok {
		return reflect.TypeOf(field[FieldDataType]).Kind() == reflect.String && reflect.TypeOf(field[FieldUI]).Kind() == reflect.String
	}

	return false
}

func IsFieldAGroup(f any) bool {
	if fgProperty, ok := f.(gojsoncore.JsonObject); ok {
		if gReadOrderOfFields, ok := fgProperty[GroupReadOrderOfFields].(gojsoncore.JsonArray); ok && len(gReadOrderOfFields) > 0 {
			if gFields, ok := fgProperty[GroupFields].(gojsoncore.JsonArray); ok && len(gFields) > 0 {
				if _, ok := gFields[0].(gojsoncore.JsonObject); ok {
					return true
				}
			}
		}
	}

	return false
}

func GetJsonPathToValue(fgKey string, removeGroupFields bool, arrayIndexPlaceholder string) string {
	if removeGroupFields {
		fgKey = strings.Replace(fgKey, path.JsonpathDotNotation+GroupFields+ArrayPathPlaceholder, "", 1)
		fgKey = string(GroupFieldsRegexSearch().ReplaceAll([]byte(fgKey), []byte("")))
	}
	fgKey = string(ArrayPathRegexSearch().ReplaceAll([]byte(fgKey), []byte(arrayIndexPlaceholder)))
	return fgKey
}

type I2DFieldViewPosition struct {
	FgKey                                   string
	FViewValuesInSeparateColumnsHeaderIndex *int
	FieldPositionBefore                     *bool
}

func Get2DFieldViewPosition(f gojsoncore.JsonObject) *I2DFieldViewPosition {
	if value, ok := f[Field2dViewPosition].(gojsoncore.JsonObject); ok {
		if jsonData, err := json.Marshal(value); err != nil {
			return nil
		} else {
			var fieldViewPosition *I2DFieldViewPosition
			if err := json.Unmarshal(jsonData, fieldViewPosition); err != nil {
				return nil
			}
			return fieldViewPosition
		}
	}

	return nil
}

func Is2DFieldViewPositionValid(f any) bool {
	if field, ok := f.(gojsoncore.JsonObject); ok {
		if field2dPosition, ok := field[Field2dViewPosition].(gojsoncore.JsonObject); ok {
			return reflect.TypeOf(field2dPosition[FieldGroupJsonPathKey]).Kind() == reflect.String
		}
	}

	return false
}

func AsJSONPath(v any) (path.JSONPath, error) {
	if value, err := gojsoncore.As[path.JSONPath](v); err == nil {
		return value, nil
	}

	if value, err := gojsoncore.As[string](v); err == nil {
		return path.JSONPath(value), nil
	}

	return "", ErrArgumentInvalid
}

func AsJsonObject(v any) (gojsoncore.JsonObject, error) {
	if v, ok := v.(gojsoncore.JsonObject); ok {
		return v, nil
	}

	if v, ok := v.(map[string]interface{}); ok {
		return v, nil
	}

	return nil, ErrArgumentInvalid
}

func AsJsonArray(v any) (gojsoncore.JsonArray, error) {
	if v, ok := v.(gojsoncore.JsonArray); ok {
		return v, nil
	}

	if v, ok := v.([]any); ok {
		return v, nil
	}

	return nil, ErrArgumentInvalid
}

func GetGroupReadOrderOfFields(fg any) (MetadataModelGroupReadOrderOfFields, error) {
	if fgProperty, ok := fg.(gojsoncore.JsonObject); ok {
		if gReadOrderOfFields, ok := fgProperty[GroupReadOrderOfFields].(gojsoncore.JsonArray); ok {
			res := make(MetadataModelGroupReadOrderOfFields, len(gReadOrderOfFields))
			for index, gReadOrderOfField := range gReadOrderOfFields {
				if str, ok := gReadOrderOfField.(string); ok {
					res[index] = str
				} else {
					return nil, fmt.Errorf("gReadOrderOfField is not string: %w", ErrArgumentInvalid)
				}
			}
			return res, nil
		}
		return nil, fmt.Errorf("gReadOrderOfFields is not []any: %w", ErrArgumentInvalid)
	}

	return nil, ErrArgumentInvalid
}

func GetGroupFields(fg any) (gojsoncore.JsonObject, error) {
	if fgProperty, ok := fg.(gojsoncore.JsonObject); ok {
		if gFields, ok := fgProperty[GroupFields].(gojsoncore.JsonArray); ok && len(gFields) > 0 {
			if gFieldsMap, ok := gFields[0].(gojsoncore.JsonObject); ok {
				return gFieldsMap, nil
			}
			return nil, fmt.Errorf("gFieldsMap is not gojsoncore.JsonObject: %w", ErrArgumentInvalid)
		}
		return nil, fmt.Errorf("gFields is not []any: %w", ErrArgumentInvalid)
	}

	return nil, fmt.Errorf("fgProperty is not gojsoncore.JsonObject: %w", ErrArgumentInvalid)
}

func GetFieldGroupName(fg any, defaultValue string) string {
	if fieldGroup, ok := fg.(gojsoncore.JsonObject); ok {
		if fieldGroupName, ok := fieldGroup[FieldGroupName].(string); ok && len(fieldGroupName) > 0 {
			return fieldGroupName
		}

		if fieldGroupKey, ok := fieldGroup[FieldGroupJsonPathKey].(string); ok && len(fieldGroupKey) > 0 {
			fieldGroupKeyParts := strings.Split(fieldGroupKey, ".")
			if len(fieldGroupKeyParts) > 0 {
				return fieldGroupKeyParts[len(fieldGroupKeyParts)-1]
			}
		}
	}

	if len(defaultValue) == 0 {
		return "#unnamed"
	}

	return defaultValue
}

func GetFieldGroupJsonPathKeySuffix(fg any) string {
	if fieldGroup, ok := fg.(gojsoncore.JsonObject); ok {
		if fieldGroupKey, ok := fieldGroup[FieldGroupJsonPathKey].(string); ok && len(fieldGroupKey) > 0 {
			fieldGroupKeyParts := strings.Split(fieldGroupKey, ".")
			if len(fieldGroupKeyParts) > 0 {
				return fieldGroupKeyParts[len(fieldGroupKeyParts)-1]
			}
		}
	}
	return ""
}
