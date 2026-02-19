package filter

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
)

/*
FilterProcessors represents a map of filter condition processors.

The key being a unique FilterCondition like FilterConditionNoOfEntriesGreaterThan and the value being a function that returns whether condition is true.
*/
type FilterProcessors map[string]ConditionTrue

/*
ConditionTrue

Parameters:
  - ctx - Current module processing the filter query.
  - fieldGroupJsonPathKey - Current core.FieldGroupJsonPathKey for field/group in metadata model. Can be used to fetch the field/group properties using ctx.
  - filterCondition - the filter condition key that caused the function to be called.
  - valueFound - the value to execute the filter against.
  - filterValue - the filter condition value.
*/
type ConditionTrue func(ctx FilterContext, fieldGroupJsonPathKey path.JSONPath, filterCondition string, valueFound reflect.Value, filterValue gojsoncore.JsonObject) (bool, error)

// FilterContext represents the context in which a filter is executed.
/*
FilterContext represents the current module processing the query.
*/
type FilterContext interface {
	// GetFieldGroupByJsonPathKey retrieves field/group properties at path.
	GetFieldGroupByJsonPathKey(path path.JSONPath) (gojsoncore.JsonObject, error)

	// SilenceErrors if set to `true`, errors encountered default to the current context condition being `false`.
	SilenceErrors() bool
}

// Constants for filter condition properties.
const (
	FilterConditionValue            string = "Value"
	FilterConditionValues           string = "Values"
	FilterConditionDateTimeFormat   string = "DateTimeFormat"
	FilterConditionCaseInsensitive  string = "CaseInsensitive"
	FilterConditionAssumedFieldType string = "AssumedFieldType"
)

// Constants for query condition properties.
const (
	QueryConditionType   string = "Type"
	QueryConditionNegate string = "Negate"
	QueryConditionValue  string = "Value"
)

/*
The "Type" property in the query module.
*/
const (
	QuerySectionTypeLogicalOperator string = "LogicalOperator"
	QuerySectionTypeFieldGroup      string = "FieldGroup"
)

// Constants for logical operators.
const (
	QuerySectionTypeLogicalOperatorAnd string = "And"
	QuerySectionTypeLogicalOperatorOr  string = "Or"
)

// LogicalOperators returns a list of supported logical operators.
func LogicalOperators() []string {
	return []string{QuerySectionTypeLogicalOperatorAnd, QuerySectionTypeLogicalOperatorOr}
}

// GetQuerySectionTypeLogicalOperator retrieves the logical operator from the query condition.
func GetQuerySectionTypeLogicalOperator(queryCondition gojsoncore.JsonObject) (string, error) {
	logicalOperator := QuerySectionTypeLogicalOperatorAnd
	if value, ok := queryCondition[QuerySectionTypeLogicalOperator].(string); ok {
		logicalOperator = value
	} else {
		return logicalOperator, nil
	}

	if !slices.Contains(LogicalOperators(), logicalOperator) {
		return "", errors.New(fmt.Sprintf("invalid logical operator '%s'", logicalOperator))
	}
	return logicalOperator, nil
}

// DefaultFilterProcessors returns a set of filter processors built on assumption of json-like data.
func DefaultFilterProcessors() FilterProcessors {
	return FilterProcessors{
		FilterConditionNoOfEntriesGreaterThan: IsConditionTrue,
		FilterConditionNoOfEntriesLessThan:    IsConditionTrue,
		FilterConditionNoOfEntriesEqualTo:     IsConditionTrue,
		FilterConditionGreaterThan:            IsConditionTrue,
		FilterConditionLessThan:               IsConditionTrue,
		FilterConditionEqualTo:                IsConditionTrue,
		FilterConditionBeginsWith:             IsConditionTrue,
		FilterConditionEndsWith:               IsConditionTrue,
		FilterConditionContains:               IsConditionTrue,
	}
}

/*
Filter query conditions
*/
const (
	FilterConditionNoOfEntriesGreaterThan string = "NoOfEntriesGreaterThan"
	FilterConditionNoOfEntriesLessThan    string = "NoOfEntriesLessThan"
	FilterConditionNoOfEntriesEqualTo     string = "NoOfEntriesEqualTo"
	FilterConditionGreaterThan            string = "GreaterThan"
	FilterConditionLessThan               string = "LessThan"
	FilterConditionEqualTo                string = "EqualTo"
	FilterConditionBeginsWith             string = "BeginsWith"
	FilterConditionEndsWith               string = "EndsWith"
	FilterConditionContains               string = "Contains"
)

/*
Filterable represents a unique filter condition.

Should serve as the base for building custom filter processing tools.
*/
type Filterable interface {
	// UniqueKey represents a unique filter condition key in a system.
	UniqueKey() string
}

var (
	// ErrFilterError default error for filter functions.
	ErrFilterError = errors.New("filter execution failed")

	ErrUnsupportedFilterConditionType = errors.New("unsupported filter condition type")

	ErrFilterConditionPropertyNotFound = errors.New("filter condition property not found")
)

// NewError creates a new core.Error with the default filter error base.
func NewError() *core.Error {
	n := core.NewError().WithDefaultBaseError(ErrFilterError)
	return n
}
