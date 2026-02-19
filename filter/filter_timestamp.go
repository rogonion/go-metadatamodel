package filter

import (
	"fmt"
	"reflect"
	"time"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
)

// IsTimestampConditionTrue checks if a timestamp condition is met.
func IsTimestampConditionTrue(ctx FilterContext, _ path.JSONPath, filterCondition string, valueFound reflect.Value, filterValue gojsoncore.JsonObject) (bool, error) {
	const FunctionName = "IsTimestampConditionTrue"

	if !valueFound.IsValid() {
		return false, nil
	}
	valueFoundInterface := valueFound.Interface()

	var dateTimeFormat string
	if filterDateTimeFormat, ok := filterValue[FilterConditionDateTimeFormat]; ok {
		if value, ok := filterDateTimeFormat.(string); ok {
			dateTimeFormat = value
		} else {
			if ctx.SilenceErrors() {
				return false, nil
			}
			return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' is not a string", FilterConditionDateTimeFormat)).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
		}
	} else {
		if ctx.SilenceErrors() {
			return false, nil
		}
		return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' not found", FilterConditionDateTimeFormat)).WithNestedError(ErrFilterConditionPropertyNotFound).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
	}

	var valueToCompare []time.Time
	if filterConditionValue, ok := filterValue[FilterConditionValue]; ok {
		if value, ok := filterConditionValue.(time.Time); ok {
			valueToCompare = append(valueToCompare, value)
		} else if value, ok := filterConditionValue.(*time.Time); ok {
			valueToCompare = append(valueToCompare, *value)
		} else if value, ok := filterConditionValue.(string); ok {
			if parsedTime, err := time.Parse(time.RFC3339Nano, value); err != nil {
				if ctx.SilenceErrors() {
					return false, nil
				}
				return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Error parsing filter condition value `%s` to time", value)).WithNestedError(err)
			} else {
				valueToCompare = append(valueToCompare, parsedTime)
			}
		} else {
			if ctx.SilenceErrors() {
				return false, nil
			}
			return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' to be time.Time or string", FilterConditionValue))
		}
	} else if filterConditionValues, ok := filterValue[FilterConditionValues]; ok {
		if value, ok := filterConditionValues.([]any); ok {
			for _, v := range value {
				if value, ok := v.(time.Time); ok {
					valueToCompare = append(valueToCompare, value)
				} else if value, ok := v.(*time.Time); ok {
					valueToCompare = append(valueToCompare, *value)
				} else if value, ok := v.(string); ok {
					if parsedTime, err := time.Parse(time.RFC3339Nano, value); err != nil {
						if ctx.SilenceErrors() {
							return false, nil
						}
						return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Error parsing filter condition value `%s` to time", value)).WithNestedError(err)
					} else {
						valueToCompare = append(valueToCompare, parsedTime)
					}
				} else {
					if ctx.SilenceErrors() {
						return false, nil
					}
					return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' to be time.Time or string", FilterConditionValue))
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
		return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter condition property '%s' not found", FilterConditionValue)).WithNestedError(ErrFilterConditionPropertyNotFound).WithData(gojsoncore.JsonObject{"FilterConditionValue": filterValue})
	}

	var valueFoundTime time.Time
	if value, ok := valueFoundInterface.(time.Time); ok {
		valueFoundTime = value
	} else if value, ok := valueFoundInterface.(*time.Time); ok {
		valueFoundTime = *value
	} else if value, ok := valueFoundInterface.(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339Nano, value); err != nil {
			return false, nil
		} else {
			valueFoundTime = parsedTime
		}
	} else {
		return false, nil
	}

	for _, value := range valueToCompare {
		switch filterCondition {
		case FilterConditionGreaterThan:
			switch dateTimeFormat {
			case core.FieldDatetimeFormatYYYYMMDDHHMM:
				vfYear, vfMonth, vfDay := valueFoundTime.Date()
				fvYear, fvMonth, fvDay := value.Date()
				if vfYear > fvYear {
					return true, nil
				}
				if vfYear == fvYear {
					if vfMonth > fvMonth {
						return true, nil
					}
					if vfMonth == fvMonth {
						if vfDay > fvDay {
							return true, nil
						}
						if vfDay == fvDay {
							if valueFoundTime.Hour() > value.Hour() {
								return true, nil
							}
							if valueFoundTime.Hour() == value.Hour() {
								if valueFoundTime.Minute() > value.Minute() {
									return true, nil
								}
							}
						}
					}
				}
			case core.FieldDatetimeFormatYYYYMMDD:
				vfYear, vfMonth, vfDay := valueFoundTime.Date()
				fvYear, fvMonth, fvDay := value.Date()
				if vfYear > fvYear {
					return true, nil
				}
				if vfYear == fvYear {
					if vfMonth > fvMonth {
						return true, nil
					}
					if vfMonth == fvMonth {
						if vfDay > fvDay {
							return true, nil
						}
					}
				}
			case core.FieldDatetimeFormatYYYYMM:
				vfYear, vfMonth, _ := valueFoundTime.Date()
				fvYear, fvMonth, _ := value.Date()
				if vfYear > fvYear {
					return true, nil
				}
				if vfYear == fvYear {
					if vfMonth > fvMonth {
						return true, nil
					}
				}
			case core.FieldDatetimeFormatHHMM:
				if valueFoundTime.Hour() > value.Hour() {
					return true, nil
				}
				if valueFoundTime.Hour() == value.Hour() {
					if valueFoundTime.Minute() > value.Minute() {
						return true, nil
					}
				}
			case core.FieldDatetimeFormatYYYY:
				return valueFoundTime.Year() > value.Year(), nil
			case core.FieldDatetimeFormatMM:
				_, vfMonth, _ := valueFoundTime.Date()
				_, fvMonth, _ := value.Date()
				if vfMonth > fvMonth {
					return true, nil
				}
			}
		case FilterConditionLessThan:
			switch dateTimeFormat {
			case core.FieldDatetimeFormatYYYYMMDDHHMM:
				vfYear, vfMonth, vfDay := valueFoundTime.Date()
				fvYear, fvMonth, fvDay := value.Date()
				if vfYear < fvYear {
					return true, nil
				}
				if vfYear == fvYear {
					if vfMonth < fvMonth {
						return true, nil
					}
					if vfMonth == fvMonth {
						if vfDay < fvDay {
							return true, nil
						}
						if vfDay == fvDay {
							if valueFoundTime.Hour() < value.Hour() {
								return true, nil
							}
							if valueFoundTime.Hour() == value.Hour() {
								if valueFoundTime.Minute() < value.Minute() {
									return true, nil
								}
							}
						}
					}
				}
			case core.FieldDatetimeFormatYYYYMMDD:
				vfYear, vfMonth, vfDay := valueFoundTime.Date()
				fvYear, fvMonth, fvDay := value.Date()
				if vfYear < fvYear {
					return true, nil
				}
				if vfYear == fvYear {
					if vfMonth < fvMonth {
						return true, nil
					}
					if vfMonth == fvMonth {
						if vfDay < fvDay {
							return true, nil
						}
					}
				}
			case core.FieldDatetimeFormatYYYYMM:
				vfYear, vfMonth, _ := valueFoundTime.Date()
				fvYear, fvMonth, _ := value.Date()
				if vfYear < fvYear {
					return true, nil
				}
				if vfYear == fvYear {
					if vfMonth < fvMonth {
						return true, nil
					}
				}
			case core.FieldDatetimeFormatHHMM:
				if valueFoundTime.Hour() < value.Hour() {
					return true, nil
				}
				if valueFoundTime.Hour() == value.Hour() {
					if valueFoundTime.Minute() < value.Minute() {
						return true, nil
					}
				}
			case core.FieldDatetimeFormatYYYY:
				return valueFoundTime.Year() < value.Year(), nil
			case core.FieldDatetimeFormatMM:
				_, vfMonth, _ := valueFoundTime.Date()
				_, fvMonth, _ := valueFoundTime.Date()
				if vfMonth < fvMonth {
					return true, nil
				}
			}
		case FilterConditionEqualTo:
			switch dateTimeFormat {
			case core.FieldDatetimeFormatYYYYMMDDHHMM:
				vfYear, vfMonth, vfDay := valueFoundTime.Date()
				fvYear, fvMonth, fvDay := value.Date()
				if vfYear == fvYear {
					if vfMonth == fvMonth {
						if vfDay == fvDay {
							if valueFoundTime.Hour() == value.Hour() {
								if valueFoundTime.Minute() == value.Minute() {
									return true, nil
								}
							}
						}
					}
				}
			case core.FieldDatetimeFormatYYYYMMDD:
				vfYear, vfMonth, vfDay := valueFoundTime.Date()
				fvYear, fvMonth, fvDay := value.Date()
				if vfYear == fvYear {
					if vfMonth == fvMonth {
						if vfDay == fvDay {
							return true, nil
						}
					}
				}
			case core.FieldDatetimeFormatYYYYMM:
				vfYear, vfMonth, _ := valueFoundTime.Date()
				fvYear, fvMonth, _ := value.Date()
				if vfYear == fvYear {
					if vfMonth == fvMonth {
						return true, nil
					}
				}
			case core.FieldDatetimeFormatHHMM:
				if valueFoundTime.Hour() == value.Hour() {
					if valueFoundTime.Minute() == value.Minute() {
						return true, nil
					}
				}
			case core.FieldDatetimeFormatYYYY:
				if valueFoundTime.Year() == value.Year() {
					return true, nil
				}
			case core.FieldDatetimeFormatMM:
				_, vfMonth, _ := valueFoundTime.Date()
				_, fvMonth, _ := value.Date()
				if vfMonth == fvMonth {
					return true, nil
				}
			}
		}
		if ctx.SilenceErrors() {
			return false, nil
		}
		return false, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Unsupported filter condition '%s'", filterCondition)).WithNestedError(ErrUnsupportedFilterConditionType).WithData(gojsoncore.JsonObject{"FilterValue": filterValue})
	}
	return false, nil
}
