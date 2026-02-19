package core

import (
	"reflect"
	"testing"

	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/internal"
)

func TestCore_JsonPathToValue(t *testing.T) {

	for testData := range JsonPathToValueTestData {
		jptv := NewJsonPathToValue().
			WithRemoveGroupFields(testData.RemoveGroupFields).
			WithSourceOfValueIsAnArray(testData.SourceOfValueIsAnArray).
			WithReplaceArrayPathPlaceholderWithActualIndexes(testData.ReplaceArrayPathPlaceholderWithActualIndexes)

		res, err := jptv.Get(testData.Path, testData.ArrayIndexes)

		if err != nil && testData.LogErrorsIfExpectedNotOk {
			t.Error(testData.TestTitle, "\n", "Error=", err)
		}

		if reflect.DeepEqual(res, testData.Expected) != testData.ExpectedOk {
			t.Error(
				testData.TestTitle, "\n",
				"expected DeepEqual result to be equal to ExpectedEqual", "\n",
				"testData.Path=", testData.Path, "\n",
				"testData.Expected=", testData.Expected, "\n",
				"res=", res,
			)
		}
	}
}

type JsonPathToValueData struct {
	internal.TestData
	Path                                         path.JSONPath
	RemoveGroupFields                            bool
	SourceOfValueIsAnArray                       bool
	ReplaceArrayPathPlaceholderWithActualIndexes bool
	ArrayIndexes                                 []int
	ExpectedOk                                   bool
	Expected                                     path.JSONPath
}

func JsonPathToValueTestData(yield func(data *JsonPathToValueData) bool) {
	testCaseIndex := 1

	// Case 1: Remove GroupFields, Source is Array
	// Scenario: The source is an array (e.g. `[...]`), so the first `GroupFields[*]` is removed entirely.
	// Remaining `[*]` placeholders are replaced by default index `0` because `ReplaceArrayPathPlaceholderWithActualIndexes` is true but no indexes are provided.
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: "Remove GroupFields, Source is Array, Default Index (0)",
			},
			Path:                   "$.GroupFields[*].Group1.GroupFields[*].Group1Field",
			RemoveGroupFields:      true,
			SourceOfValueIsAnArray: true,
			ReplaceArrayPathPlaceholderWithActualIndexes: true,
			ExpectedOk: true,
			Expected:   "$[0].Group1[0].Group1Field",
		},
	) {
		return
	}

	testCaseIndex++
	// Case 2: Remove GroupFields, Source is Object
	// Scenario: The source is an object (default), so the first `GroupFields[*]` is stripped of the `GroupFields` key but keeps the structure relative to root.
	// Remaining `[*]` placeholders are replaced by default index `0`.
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: "Remove GroupFields, Source is Object, Default Index (0)",
			},
			Path:              "$.GroupFields[*].Group1.GroupFields[*].Group1Field",
			RemoveGroupFields: true,
			ExpectedOk:        true,
			ReplaceArrayPathPlaceholderWithActualIndexes: true,
			Expected: "$.Group1[0].Group1Field",
		},
	) {
		return
	}

	testCaseIndex++
	// Case 3: Keep GroupFields, Specific Indexes
	// Scenario: `RemoveGroupFields` is false, so the path structure is preserved.
	// `[*]` placeholders are replaced by the provided `ArrayIndexes` [1, 2].
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: "Keep GroupFields, Specific Indexes [1, 2]",
			},
			Path:         "$.GroupFields[*].Group1.GroupFields[*].Group1Field",
			ExpectedOk:   true,
			ArrayIndexes: []int{1, 2},
			Expected:     "$.GroupFields[1].Group1.GroupFields[2].Group1Field",
		},
	) {
		return
	}

	testCaseIndex++
	// Case 4: Remove GroupFields, Source is Array, Specific Indexes
	// Scenario: Source is array (strips first group), remaining `[*]` replaced by `ArrayIndexes`.
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: "Remove GroupFields, Source is Array, Specific Indexes [1, 2]",
			},
			Path:                   "$.GroupFields[*].Group1.GroupFields[*].Group1Field",
			RemoveGroupFields:      true,
			SourceOfValueIsAnArray: true,
			ExpectedOk:             true,
			ArrayIndexes:           []int{1, 2},
			Expected:               "$[1].Group1[2].Group1Field",
		},
	) {
		return
	}

	testCaseIndex++
	// Case 5: Remove GroupFields, Source is Object, Single Index
	// Scenario: Source is object. `ArrayIndexes` has one value [2], which replaces the remaining `[*]` after group removal.
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: "Remove GroupFields, Source is Object, Single Index [2]",
			},
			Path:              "$.GroupFields[*].Group1.GroupFields[*].Group1Field",
			RemoveGroupFields: true,
			ExpectedOk:        true,
			ArrayIndexes:      []int{2},
			Expected:          "$.Group1[2].Group1Field",
		},
	) {
		return
	}
}
