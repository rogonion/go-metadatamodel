package filter

import (
	"fmt"
	"reflect"
	"strings"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
)

// IsTextConditionTrue checks if a text condition is met.
func IsTextConditionTrue(ctx FilterContext, _ path.JSONPath, filterCondition string, valueFound reflect.Value, filterValue gojsoncore.JsonObject) (bool, error) {
	const FunctionName = "IsTextConditionTrue"

	if !valueFound.IsValid() {
		return false, nil
	}
	valueFoundInterface := valueFound.Interface()

	caseInsensitive := false
	if value, ok := filterValue[FilterConditionCaseInsensitive]; ok {
		if value, ok := value.(bool); ok {
			caseInsensitive = value
		}
	}

	var valueToCompare []string
	if filterConditionValue, ok := filterValue[FilterConditionValue]; ok {
		if filterConditionValueString, ok := filterConditionValue.(string); ok {
			if caseInsensitive {
				valueToCompare = append(valueToCompare, strings.ToLower(filterConditionValueString))
			} else {
				valueToCompare = append(valueToCompare, filterConditionValueString)
			}
		} else {
			if ctx.SilenceErrors() {
				return false, nil
			}
			return false, NewError().WithMessage(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' is not a string", FilterConditionDateTimeFormat)).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
		}
	} else if filterConditionValues, ok := filterValue[FilterConditionValues]; ok {
		if value, ok := filterConditionValues.([]any); ok {
			for _, v := range value {
				if filterConditionValueString, ok := v.(string); ok {
					if caseInsensitive {
						valueToCompare = append(valueToCompare, strings.ToLower(filterConditionValueString))
					} else {
						valueToCompare = append(valueToCompare, filterConditionValueString)
					}
				} else {
					if ctx.SilenceErrors() {
						return false, nil
					}
					return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' is not a string", FilterConditionDateTimeFormat)).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
				}
			}
		} else {
			if ctx.SilenceErrors() {
				return false, nil
			}
			return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' is not a []any", FilterConditionValues)).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
		}
	} else {
		if ctx.SilenceErrors() {
			return false, nil
		}
		return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' not found", FilterConditionValue)).WithData(gojsoncore.JsonObject{"FilterValue": filterValue}).WithNestedError(ErrFilterConditionPropertyNotFound)
	}

	var valueFoundString string
	if value, ok := valueFoundInterface.(string); ok {
		valueFoundString = value
	} else {
		return false, nil
	}

	for _, value := range valueToCompare {
		switch filterCondition {
		case FilterConditionEqualTo:
			if value == valueFoundString {
				return true, nil
			}
		case FilterConditionBeginsWith:
			if strings.HasPrefix(valueFoundString, value) {
				return true, nil
			}
		case FilterConditionEndsWith:
			if strings.HasSuffix(valueFoundString, value) {
				return true, nil
			}
		case FilterConditionContains:
			if strings.Contains(valueFoundString, value) {
				return true, nil
			}
		default:
			if ctx.SilenceErrors() {
				return false, nil
			}
			return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Unsupported filter condition '%s'", filterCondition)).WithData(gojsoncore.JsonObject{"FilterValue": filterValue}).WithNestedError(ErrUnsupportedFilterConditionType)
		}
	}
	return false, nil
}
