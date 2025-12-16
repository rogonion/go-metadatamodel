package core

import (
	"fmt"
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
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d", testCaseIndex),
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
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d", testCaseIndex),
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
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d", testCaseIndex),
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
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d", testCaseIndex),
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
	if !yield(
		&JsonPathToValueData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d", testCaseIndex),
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
