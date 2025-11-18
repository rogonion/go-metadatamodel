package filter

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
)

/*
Filter

Parameters:
  - queryConditions
  - rootJsonPathKey - Set sub-set of metadata model with DataFilter.metadataModel as root context.
  - rootJsonPathToValue - Path to data in DataFilter.sourceData that will act as root context.

Returns:
 1. Array of indexes that DID NOT pass the filter test.
 2. An error especially if queryConditions is not valid or nil if DataFilter.silenceAllErrors is `true`.
*/
func (n *DataFilter) Filter(queryCondition gojsoncore.JsonObject, rootJsonPathKey path.JSONPath, rootJsonPathToValue path.JSONPath) ([]int, error) {
	const FunctionName = "Filter"

	if len(rootJsonPathKey) == 0 {
		rootJsonPathKey = path.JSONPath(path.JsonpathKeyRoot)
	}
	n.rootJsonPathKey = rootJsonPathKey
	if len(rootJsonPathToValue) == 0 {
		rootJsonPathToValue = path.JSONPath(path.JsonpathKeyRoot)
	}
	n.rootJsonPathToValue = rootJsonPathToValue

	if noOfResults, err := n.sourceData.Get(n.rootJsonPathToValue); noOfResults == 0 {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get root value yielded 0 results").WithNestedError(err)
	}

	if n.sourceData.GetValueFoundReflected().Kind() != reflect.Slice && n.sourceData.GetValueFoundReflected().Kind() != reflect.Array {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("root value should be slice or array")
	}

	filterExcludeIndexes := make([]int, 0)
	var returnErr error
	object.NewObject().WithSourceReflected(n.sourceData.GetValueFoundReflected()).ForEach(path.JSONPath(path.JsonpathKeyRoot+core.ArrayPathPlaceholder), func(jsonPath path.RecursiveDescentSegment, value reflect.Value) bool {
		if len(queryCondition) == 0 {
			return false
		}

		lastPathSegment := jsonPath[len(jsonPath)-1]

		if !lastPathSegment.IsIndex {
			if n.silenceAllErrors {
				return false
			}
			returnErr = NewError().WithFunctionName(FunctionName).WithMessage("in root value loop, last path segment is not an index").WithData(gojsoncore.JsonObject{"Path": jsonPath})
			return true
		}

		//fmt.Println("--------------")
		//fmt.Println("Index", lastPathSegment)
		//fmt.Println("Value", gojsoncore.JsonStringifyMust(value.Interface()))

		ok, err := n.isQueryConditionTrue(queryCondition, value, jsonPath)
		if err != nil {
			if n.silenceAllErrors {
				return false
			}
			returnErr = err
			return true
		}

		if !ok {
			filterExcludeIndexes = append(filterExcludeIndexes, lastPathSegment.Index)
		}

		//fmt.Println("--------------")

		return false
	})
	return filterExcludeIndexes, returnErr
}

func (n *DataFilter) isQueryConditionTrue(queryCondition gojsoncore.JsonObject, currentValue reflect.Value, jsonPath path.RecursiveDescentSegment) (bool, error) {
	const FunctionName = "isQueryConditionTrue"

	var queryConditionType string
	if value, ok := queryCondition[QueryConditionType].(string); ok {
		queryConditionType = value
	} else {
		return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Key '%s' is not valid", QueryConditionType)).WithData(gojsoncore.JsonObject{"QueryCondition": queryCondition}))
	}

	switch queryConditionType {
	case QuerySectionTypeLogicalOperator:
		return n.isRecursiveLogicalOperatorTrue(queryCondition, currentValue, jsonPath)
	case QuerySectionTypeFieldGroup:
		return n.isRecursiveFieldGroupTrue(queryCondition, currentValue, jsonPath)
	default:
		return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Unknown query condition type: %s", queryConditionType)).WithData(gojsoncore.JsonObject{"QueryCondition": queryCondition}))
	}
}

func (n *DataFilter) isRecursiveLogicalOperatorTrue(queryCondition gojsoncore.JsonObject, currentValue reflect.Value, jsonPath path.RecursiveDescentSegment) (bool, error) {
	const FunctionName = "isRecursiveLogicalOperatorTrue"

	negate := false
	if value, ok := queryCondition[QueryConditionNegate].(bool); ok {
		negate = value
	}

	logicalOperator, err := GetQuerySectionTypeLogicalOperator(queryCondition)
	if err != nil {
		return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage("Invalid logical operator").WithNestedError(err))
	}

	var conditions gojsoncore.JsonArray
	if value, err := core.AsJsonArray(queryCondition[QueryConditionValue]); err == nil {
		conditions = value
	} else {
		return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Key '%s' is not valid", QueryConditionValue)).WithData(gojsoncore.JsonObject{"QueryCondition": queryCondition}).WithNestedError(err))
	}

	conditionsResults := make([]bool, 0)
	for _, condition := range conditions {
		conditionJsonObject, err := core.AsJsonObject(condition)
		if err != nil {
			return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage("condition not JsonObject").WithNestedError(err))
		}

		conditionTrue, err := n.isQueryConditionTrue(conditionJsonObject, currentValue, jsonPath)
		if err != nil {
			return n.returnErrorOrFalse(err)
		}

		if conditionTrue {
			if logicalOperator == QuerySectionTypeLogicalOperatorOr {
				if negate {
					return false, nil
				}
				return true, nil
			}
		} else {
			if logicalOperator == QuerySectionTypeLogicalOperatorAnd {
				if negate {
					return true, nil
				}
				return false, nil
			}
		}
		conditionsResults = append(conditionsResults, conditionTrue)
	}

	if logicalOperator == QuerySectionTypeLogicalOperatorOr {
		if slices.Contains(conditionsResults, true) {
			if negate {
				return false, nil
			}
			return true, nil
		} else {
			if negate {
				return true, nil
			}
			return false, nil
		}
	} else {
		if slices.Contains(conditionsResults, false) {
			if negate {
				return true, nil
			}
			return false, nil
		} else {
			if negate {
				return false, nil
			}
			return true, nil
		}
	}
}

func (n *DataFilter) isRecursiveFieldGroupTrue(queryCondition gojsoncore.JsonObject, currentValue reflect.Value, jsonPath path.RecursiveDescentSegment) (bool, error) {
	const FunctionName = "isRecursiveFieldGroupTrue"

	negate := false
	if value, ok := queryCondition[QueryConditionNegate].(bool); ok {
		negate = value
	}

	logicalOperator, err := GetQuerySectionTypeLogicalOperator(queryCondition)
	if err != nil {
		return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage("Invalid logical operator").WithData(gojsoncore.JsonObject{"QueryCondition": queryCondition}).WithNestedError(err))
	}

	conditions := make(gojsoncore.JsonObject)
	if value, err := core.AsJsonObject(queryCondition[QueryConditionValue]); err == nil {
		conditions = value
	} else {
		return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("Key '%s' is not valid", QueryConditionValue)).WithData(gojsoncore.JsonObject{"QueryCondition": queryCondition}).WithNestedError(err))
	}
	conditionsResults := make([]bool, 0)

	for jsonPathKey, condition := range conditions {
		conditionJsonObject, err := core.AsJsonObject(condition)
		if err != nil {
			return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage("condition not JsonObject").WithNestedError(err))
		}

		conditionTrue, err := n.isFieldGroupConditionTrue(path.JSONPath(jsonPathKey), conditionJsonObject, currentValue, jsonPath)
		if err != nil {
			return n.returnErrorOrFalse(err)
		}

		if !conditionTrue {
			if logicalOperator == QuerySectionTypeLogicalOperatorAnd {
				if negate {
					return true, nil
				}
				return false, nil
			}
		}
		conditionsResults = append(conditionsResults, conditionTrue)
	}

	if logicalOperator == QuerySectionTypeLogicalOperatorOr {
		if slices.Contains(conditionsResults, true) {
			if negate {
				return false, nil
			}
			return true, nil
		} else {
			if negate {
				return true, nil
			}
			return false, nil
		}
	} else {
		if slices.Contains(conditionsResults, false) {
			if negate {
				return true, nil
			}
			return false, nil
		} else {
			if negate {
				return false, nil
			}
			return true, nil
		}
	}
}

func (n *DataFilter) isFieldGroupConditionTrue(jsonPathKey path.JSONPath, queryCondition gojsoncore.JsonObject, currentValue reflect.Value, jsonPath path.RecursiveDescentSegment) (bool, error) {
	const FunctionName = "isFieldGroupConditionTrue"

	if len(queryCondition) == 0 {
		return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage("Query condition is empty").WithData(gojsoncore.JsonObject{"JsonPathKey": jsonPathKey, "QueryCondition": queryCondition, "jsonPath": jsonPath}))
	}

	currentJsonPathKey := path.JSONPath(strings.Replace(string(jsonPathKey), string(n.rootJsonPathKey), path.JsonpathKeyRoot, 1))
	currentJsonPathToValue, err := core.NewJsonPathToValue().WithReplaceArrayPathPlaceholderWithActualIndexes(false).Get(currentJsonPathKey, nil)
	if err != nil {
		return n.returnErrorOrFalse(NewError().WithFunctionName(FunctionName).WithMessage("get current json path to value failed").WithData(gojsoncore.JsonObject{"CurrentJsonPathKey": currentJsonPathKey, "QueryCondition": queryCondition}).WithNestedError(err))
	}

	orConditionTrue := false
	var loopError error

	object.NewObject().WithSourceReflected(currentValue).ForEach(currentJsonPathToValue, func(jsonPath path.RecursiveDescentSegment, value reflect.Value) bool {
		//fmt.Println(jsonPath)
		//fmt.Println(queryCondition)
		andConditionTrue := true
		for filterConditionKey, filterConditionData := range queryCondition {
			filterConditionDataJsonObject, err := core.AsJsonObject(filterConditionData)
			if err != nil {
				if n.silenceAllErrors {
					continue
				}
				loopError = NewError().WithFunctionName(FunctionName).WithMessage("filterConditionData not JsonObject").WithNestedError(err).WithData(gojsoncore.JsonObject{"FilterConditionData": filterConditionData})
				return true
			}

			if filterProcessor, ok := n.defaultFilterProcessors[filterConditionKey]; ok {
				conditionTrue, err := filterProcessor(n, jsonPathKey, filterConditionKey, value, filterConditionDataJsonObject)
				if err != nil {
					if n.silenceAllErrors {
						continue
					}
					loopError = NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("filter processing for condition '%s' failed", filterConditionKey)).WithData(gojsoncore.JsonObject{"CurrentJsonPathKey": currentJsonPathKey, "QueryCondition": queryCondition}).WithNestedError(err)
					return true
				}
				if !conditionTrue {
					andConditionTrue = false
					break
				}
			} else {
				if n.silenceAllErrors {
					continue
				}
				loopError = NewError().WithFunctionName(FunctionName).WithMessage("filterConditionData not JsonObject").WithNestedError(err)
				return true
			}
		}
		if andConditionTrue {
			orConditionTrue = true
			return true
		}
		return false
	})
	if loopError != nil {
		return n.returnErrorOrFalse(loopError)
	}

	return orConditionTrue, nil
}

func (n *DataFilter) WithMetadataModel(value gojsoncore.JsonObject) *DataFilter {
	n.SetMetadataModel(value)
	return n
}

func (n *DataFilter) SetMetadataModel(value gojsoncore.JsonObject) {
	n.metadataModelObject = object.NewObject().WithSourceInterface(value)
}

func (n *DataFilter) WithSourceData(value *object.Object) *DataFilter {
	n.SetSourceData(value)
	return n
}

func (n *DataFilter) SetSourceData(value *object.Object) {
	n.sourceData = value
}

func (n *DataFilter) WithDefaultFilterProcessors(value FilterProcessors) *DataFilter {
	n.SetDefaultFilterProcessors(value)
	return n
}

func (n *DataFilter) SetDefaultFilterProcessors(value FilterProcessors) {
	n.defaultFilterProcessors = value
}

func (n *DataFilter) WithSilenceErrors(value bool) *DataFilter {
	n.SetSilenceErrors(value)
	return n
}

func (n *DataFilter) SetSilenceErrors(value bool) {
	n.silenceAllErrors = value
}

func (n *DataFilter) GetFieldGroupByJsonPathKey(jsonPath path.JSONPath) (gojsoncore.JsonObject, error) {
	const FunctionName = "GetFieldGroupByJsonPathKey"

	jsonPathToValue, err := core.NewJsonPathToValue().WithRemoveGroupFields(false).Get(jsonPath, nil)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get json path to value failed").WithNestedError(err)
	}

	noOfResults, err := n.metadataModelObject.Get(jsonPathToValue)
	if noOfResults == 0 {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get field/group failed").WithNestedError(err)
	}

	if jsonObject, err := core.AsJsonObject(n.metadataModelObject.GetValueFoundInterface()); err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("value found not JsonObject").WithNestedError(err)
	} else {
		return jsonObject, nil
	}
}

func (n *DataFilter) SilenceErrors() bool {
	return n.silenceAllErrors
}

func (n *DataFilter) returnErrorOrFalse(err error) (bool, error) {
	if n.silenceAllErrors {
		return false, nil
	}
	return false, err
}

/*
NewFilterData

Parameters:

  - sourceData - Refer to object.Object.
  - metadataModel - data model for sourceData.
*/
func NewFilterData(sourceData *object.Object, metadataModel gojsoncore.JsonObject) *DataFilter {
	n := new(DataFilter)
	n.SetSourceData(sourceData)
	n.SetMetadataModel(metadataModel)
	n.defaultFilterProcessors = DefaultFilterProcessors()
	return n
}

type DataFilter struct {
	// Use to loop through values using object.ForEach in source. Refer to object.
	sourceData *object.Object

	// Used by DataFilter.GetFieldGroupByJsonPathKey.
	metadataModelObject *object.Object

	// Set of functions to process filter conditions by unique filter key.
	defaultFilterProcessors FilterProcessors

	// Set the root source data within sourceData for filtering against sub-set of sourceData.
	//
	// Example: `$.GroupFields[*].Address`
	rootJsonPathKey path.JSONPath
	// Must match rootJsonPathKey format e.g., `$[2].Address`
	//
	// Combined with rootJsonPathKey, this means filter conditions will be executed against array/slice found at path in sourceData.
	rootJsonPathToValue path.JSONPath

	// if set to `true`, errors encountered default to the current context condition being `false`.
	silenceAllErrors bool
}
