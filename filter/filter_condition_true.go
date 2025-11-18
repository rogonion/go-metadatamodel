package filter

import (
	"fmt"
	"reflect"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/core"
)

func IsConditionTrue(ctx FilterContext, fieldGroupJsonPathKey path.JSONPath, filterCondition string, valueFound reflect.Value, filterValue gojsoncore.JsonObject) (bool, error) {
	const FunctionName = "IsConditionTrue"

	switch filterCondition {
	case FilterConditionNoOfEntriesGreaterThan, FilterConditionNoOfEntriesLessThan, FilterConditionNoOfEntriesEqualTo:
		return IsNumberOfEntriesConditionTrue(ctx, fieldGroupJsonPathKey, filterCondition, valueFound, filterValue)
	default:
		var assumedFieldType string
		if defFieldType, ok := filterValue[FilterConditionAssumedFieldType]; ok {
			if value, ok := defFieldType.(string); ok {
				assumedFieldType = value
			} else {
				if ctx.SilenceErrors() {
					return false, nil
				}
				return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' is not a string", FilterConditionAssumedFieldType)).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
			}
		} else {
			if ctx.SilenceErrors() {
				return false, nil
			}
			return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' not found", FilterConditionAssumedFieldType)).WithNestedError(ErrFilterConditionPropertyNotFound).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
		}

		if valueFound.Kind() == reflect.Slice || valueFound.Kind() == reflect.Array {
			for i := 0; i < valueFound.Len(); i++ {
				vFound := valueFound.Index(i)

				var conditionTrue bool
				var err error
				switch assumedFieldType {
				case core.FieldTypeText:
					conditionTrue, err = IsTextConditionTrue(ctx, fieldGroupJsonPathKey, filterCondition, vFound, filterValue)
				case core.FieldTypeNumber:
					conditionTrue, err = IsNumberConditionTrue(ctx, fieldGroupJsonPathKey, filterCondition, vFound, filterValue)
				case core.FieldTypeTimestamp:
					conditionTrue, err = IsTimestampConditionTrue(ctx, fieldGroupJsonPathKey, filterCondition, vFound, filterValue)
				default:
					conditionTrue, err = IsDefaultEqualTrue(ctx, fieldGroupJsonPathKey, filterCondition, vFound, filterValue)
				}

				if err != nil {
					if ctx.SilenceErrors() {
						return false, nil
					}
					return false, err
				}

				if conditionTrue {
					return true, nil
				}
			}
			return false, nil
		} else {
			switch assumedFieldType {
			case core.FieldTypeText:
				return IsTextConditionTrue(ctx, fieldGroupJsonPathKey, filterCondition, valueFound, filterValue)
			case core.FieldTypeNumber:
				return IsNumberConditionTrue(ctx, fieldGroupJsonPathKey, filterCondition, valueFound, filterValue)
			case core.FieldTypeTimestamp:
				return IsTimestampConditionTrue(ctx, fieldGroupJsonPathKey, filterCondition, valueFound, filterValue)
			default:
				return IsDefaultEqualTrue(ctx, fieldGroupJsonPathKey, filterCondition, valueFound, filterValue)
			}
		}
	}
}

func IsNumberOfEntriesConditionTrue(ctx FilterContext, _ path.JSONPath, filterCondition string, valueFound reflect.Value, filterValue gojsoncore.JsonObject) (bool, error) {
	const FunctionName = "IsNumberOfEntriesConditionTrue"

	if !valueFound.IsValid() {
		return false, nil
	}

	if valueFound.Kind() != reflect.Slice && valueFound.Kind() != reflect.Array {
		return false, nil
	}
	valueFoundLen := valueFound.Len()

	var valueToCompare []int
	conversion := schema.NewConversion()
	intSchema := &schema.DynamicSchemaNode{Type: reflect.TypeOf(0), Kind: reflect.Int}

	if filterConditionValue, ok := filterValue[FilterConditionValue]; ok {
		if filterConditionValueInt, ok := filterConditionValue.(int); ok {
			valueToCompare = append(valueToCompare, filterConditionValueInt)
		} else {
			if err := conversion.Convert(filterConditionValue, intSchema, &valueToCompare); err != nil {
				if ctx.SilenceErrors() {
					return false, nil
				}
				return false, NewError().WithFunctionName(FunctionName).WithMessage("convert filterConditionValueInt to int failed").WithNestedError(err)
			}
		}
	} else if filterConditionValues, ok := filterValue[FilterConditionValues]; ok {
		if value, ok := filterConditionValues.([]any); ok {
			for _, v := range value {
				if filterConditionValueInt, ok := v.(int); ok {
					valueToCompare = append(valueToCompare, filterConditionValueInt)
				} else {
					if err := conversion.Convert(filterConditionValue, intSchema, &valueToCompare); err != nil {
						if ctx.SilenceErrors() {
							return false, nil
						}
						return false, NewError().WithFunctionName(FunctionName).WithMessage("convert filterConditionValueInt to int failed").WithNestedError(err)
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
		return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' not found", FilterConditionValue)).WithNestedError(ErrFilterConditionPropertyNotFound).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
	}

	for _, value := range valueToCompare {
		switch filterCondition {
		case FilterConditionNoOfEntriesEqualTo:
			if value == valueFoundLen {
				return true, nil
			}
		case FilterConditionNoOfEntriesGreaterThan:
			if value > valueFoundLen {
				return true, nil
			}
		case FilterConditionNoOfEntriesLessThan:
			if value < valueFoundLen {
				return true, nil
			}
		default:
			if ctx.SilenceErrors() {
				return false, nil
			}
			return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Unsupported filter condition '%s'", filterCondition)).WithNestedError(ErrUnsupportedFilterConditionType).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
		}
	}
	return false, nil
}

func IsDefaultEqualTrue(ctx FilterContext, _ path.JSONPath, filterCondition string, valueFound reflect.Value, filterValue gojsoncore.JsonObject) (bool, error) {
	if !valueFound.IsValid() {
		return false, nil
	}

	const FunctionName = "IsDefaultEqualTrue"

	var valueToCompare []any
	if filterConditionValue, ok := filterValue[QueryConditionValue]; ok {
		valueToCompare = append(valueToCompare, filterConditionValue)
	} else if filterConditionValues, ok := filterValue[FilterConditionValues]; ok {
		if value, ok := filterConditionValues.([]any); ok {
			valueToCompare = value
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
		return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' or '%s' not found", FilterConditionValue, FilterConditionValues)).WithNestedError(ErrFilterConditionPropertyNotFound).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
	}

	switch filterCondition {
	case FilterConditionEqualTo:
		for _, value := range valueToCompare {
			if reflect.DeepEqual(value, valueFound.Interface()) {
				return true, nil
			}
		}
		return false, nil
	default:
		if ctx.SilenceErrors() {
			return false, nil
		}
		return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Unsupported filter condition '%s'", filterCondition)).WithNestedError(ErrUnsupportedFilterConditionType).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
	}
}
