package core

import (
	"fmt"
	"reflect"
	"strings"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-json/schema"
)

// MergeRightJsonObjectIntoLeft merges the right JsonObject into the left JsonObject.
func MergeRightJsonObjectIntoLeft(left gojsoncore.JsonObject, right gojsoncore.JsonObject) {
	if left != nil {
		for k, v := range right {
			left[k] = v
		}
	}
}

// DoesFieldGroupFieldsContainNestedGroupFields checks if the field group contains nested group fields.
func DoesFieldGroupFieldsContainNestedGroupFields(fg any) bool {
	if groupFields, err := GetGroupFields(fg); err == nil {
		if groupReadOrderOfFields, err := GetGroupReadOrderOfFields(fg); err == nil {
			for _, fgJsonPathKeySuffix := range groupReadOrderOfFields {
				if _, err := GetGroupReadOrderOfFields(groupFields[fgJsonPathKeySuffix]); err == nil {
					return true
				}
			}
		}
	}
	return false
}

/*
GetMaximumFlatNoOfColumns Retrieves the maximum no of columns/sets of columns that can be used to represent a field/group in a flat 2D table.

The value will only be extracted if fg is valid and if it contains GroupFields then each of them does not contain nested GroupFields.

Returns:
  - The maximum number. If value is less than 0, it means the field/group cannot be represented.
  - Error if fg and its properties are not valid structure-wise.
*/
func GetMaximumFlatNoOfColumns(fg any) (int, error) {
	fgViewMaxNoOfValuesInSeparateColumns := -1
	if fgProperty, err := AsJsonObject(fg); err == nil {
		if value, ok := fgProperty[FieldGroupViewValuesInSeparateColumns].(bool); ok && value {
			if fgFields, err := GetGroupFields(fg); err == nil {
				for _, value := range fgFields {
					if IsFieldAGroup(value) {
						return fgViewMaxNoOfValuesInSeparateColumns, nil
					}
				}
			}

			if v, ok := fgProperty[FieldGroupViewMaxNoOfValuesInSeparateColumns]; ok {
				if err = schema.NewConversion().Convert(v, &schema.DynamicSchemaNode{Type: reflect.TypeOf(0), Kind: reflect.Int}, &fgViewMaxNoOfValuesInSeparateColumns); err != nil {
					return fgViewMaxNoOfValuesInSeparateColumns, err
				}
			}
		}
	} else {
		return fgViewMaxNoOfValuesInSeparateColumns, err
	}

	return fgViewMaxNoOfValuesInSeparateColumns, nil
}

// IfKeySuffixMatchesValues checks if the key suffix matches any of the values.
func IfKeySuffixMatchesValues(keyToCheck string, valuesToMatch []string) bool {
	for _, value := range valuesToMatch {
		if strings.HasSuffix(keyToCheck, value) {
			return true
		}
	}

	return false
}

// IsFieldAField checks if the input is a field (has FieldDataType and FieldUI).
func IsFieldAField(f any) bool {
	if field, err := AsJsonObject(f); err == nil {
		if _, ok := field[FieldDataType].(string); ok {
			if _, ok := field[FieldUI].(bool); ok {
				return true
			}
		}
	}

	return false
}

// IsFieldAGroup checks if the input is a group (has GroupReadOrderOfFields and GroupFields).
func IsFieldAGroup(f any) bool {
	if _, err := GetGroupReadOrderOfFields(f); err == nil {
		if _, err := GetGroupFields(f); err == nil {
			return true
		}
	}

	return false
}

// AsJSONPath converts the input to a JSONPath.
func AsJSONPath(v any) (path.JSONPath, error) {
	if value, ok := v.(path.JSONPath); ok {
		return value, nil
	}

	if value, ok := v.(string); ok {
		return path.JSONPath(value), nil
	}

	return "", ErrArgumentInvalid
}

// AsJsonObject converts the input to a JsonObject.
func AsJsonObject(v any) (gojsoncore.JsonObject, error) {
	if v, ok := v.(gojsoncore.JsonObject); ok {
		return v, nil
	}

	if v, ok := v.(map[string]interface{}); ok {
		return v, nil
	}

	return nil, ErrArgumentInvalid
}

// AsJsonArray converts the input to a JsonArray.
func AsJsonArray(v any) (gojsoncore.JsonArray, error) {
	if v, ok := v.(gojsoncore.JsonArray); ok {
		return v, nil
	}

	if v, ok := v.([]any); ok {
		return v, nil
	}

	return nil, ErrArgumentInvalid
}

// AsGroupReadOrderOfFields converts the input to MetadataModelGroupReadOrderOfFields.
func AsGroupReadOrderOfFields(v any) (MetadataModelGroupReadOrderOfFields, error) {
	if value, ok := v.(MetadataModelGroupReadOrderOfFields); ok {
		return value, nil
	}

	if value, ok := v.([]string); ok {
		return value, nil
	}

	if value, err := AsJsonArray(v); err == nil {
		res := make(MetadataModelGroupReadOrderOfFields, len(value))
		for index, gReadOrderOfField := range value {
			if str, ok := gReadOrderOfField.(string); ok {
				res[index] = str
			} else {
				return nil, fmt.Errorf("gReadOrderOfField is not string: %w", ErrArgumentInvalid)
			}
		}
		return res, nil
	}

	return nil, ErrArgumentInvalid
}

// GetGroupReadOrderOfFields retrieves the group read order of fields.
func GetGroupReadOrderOfFields(fg any) (MetadataModelGroupReadOrderOfFields, error) {
	if fgProperty, err := AsJsonObject(fg); err == nil {
		return AsGroupReadOrderOfFields(fgProperty[GroupReadOrderOfFields])
	}

	return nil, ErrArgumentInvalid
}

// GetGroupFields retrieves the group fields.
func GetGroupFields(fg any) (gojsoncore.JsonObject, error) {
	if fgProperty, err := AsJsonObject(fg); err == nil {
		if gFields, err := AsJsonArray(fgProperty[GroupFields]); err == nil {
			if gFieldsMap, err := AsJsonObject(gFields[0]); err == nil {
				return gFieldsMap, nil
			} else {
				return nil, fmt.Errorf("gFields[0]: %w", err)
			}
		} else {
			return nil, fmt.Errorf("gFields: %w", err)
		}
	} else {
		return nil, fmt.Errorf("fg: %w", err)
	}
}

// GetFieldGroupName retrieves the field group name.
func GetFieldGroupName(fg any, defaultValue string) string {
	if fieldGroup, err := AsJsonObject(fg); err == nil {
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

// GetFieldGroupJsonPathKeySuffix retrieves the suffix of the field group JSON path key.
func GetFieldGroupJsonPathKeySuffix(fg any) string {
	if fieldGroup, err := AsJsonObject(fg); err == nil {
		if fieldGroupKey, err := AsJSONPath(fieldGroup[FieldGroupJsonPathKey]); err == nil {
			if segments := fieldGroupKey.Parse(); len(segments) > 0 {
				return segments[len(segments)-1].String()
			}
		}
	}
	return ""
}
