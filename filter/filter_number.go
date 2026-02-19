package filter

import (
	"fmt"
	"reflect"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-json/schema"
)

// IsNumberConditionTrue checks if a number condition is met.
func IsNumberConditionTrue(ctx FilterContext, _ path.JSONPath, filterCondition string, valueFound reflect.Value, filterValue gojsoncore.JsonObject) (bool, error) {
	const FunctionName = "IsNumberConditionTrue"

	if !valueFound.IsValid() {
		return false, nil
	}
	valueFoundInterface := valueFound.Interface()
	var valueToCompare []float64
	conversion := schema.NewConversion()
	float64Schema := &schema.DynamicSchemaNode{Type: reflect.TypeOf(float64(0)), Kind: reflect.Float64}

	if filterConditionValue, ok := filterValue[FilterConditionValue]; ok {
		if filterConditionValueFloat, ok := filterConditionValue.(float64); ok {
			valueToCompare = append(valueToCompare, filterConditionValueFloat)
		} else {
			if err := conversion.Convert(filterConditionValue, float64Schema, &filterConditionValueFloat); err != nil {
				if ctx.SilenceErrors() {
					return false, nil
				}
				return false, NewError().WithFunctionName(FunctionName).WithMessage("convert filterConditionValueFloat to float64 failed").WithNestedError(err)
			} else {
				valueToCompare = append(valueToCompare, filterConditionValueFloat)
			}
		}
	} else if filterConditionValues, ok := filterValue[FilterConditionValues]; ok {
		if value, ok := filterConditionValues.([]any); ok {
			for _, v := range value {
				if filterConditionValueFloat, ok := v.(float64); ok {
					valueToCompare = append(valueToCompare, filterConditionValueFloat)
				} else {
					var newValue float64
					if err := conversion.Convert(filterConditionValue, float64Schema, &newValue); err != nil {
						if ctx.SilenceErrors() {
							return false, nil
						}
						return false, NewError().WithFunctionName(FunctionName).WithMessage("convert filterConditionValueFloat to float64 failed").WithNestedError(err)
					} else {
						valueToCompare = append(valueToCompare, newValue)
					}
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

	var valueFoundFloat float64
	if value, ok := valueFoundInterface.(float64); ok {
		valueFoundFloat = value
	} else {
		if err := conversion.Convert(valueFoundInterface, float64Schema, &valueFoundFloat); err != nil {
			return false, nil
		}
	}

	for _, value := range valueToCompare {
		switch filterCondition {
		case FilterConditionEqualTo:
			if valueFoundFloat == value {
				return true, nil
			}
		case FilterConditionGreaterThan:
			if valueFoundFloat > value {
				return true, nil
			}
		case FilterConditionLessThan:
			if valueFoundFloat < value {
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
