package filter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/internal"
	"github.com/rogonion/go-metadatamodel/testdata"
)

func TestFilter_FilterData(t *testing.T) {
	for testData := range filterDataTestData {
		fd := NewFilterData(testData.Object, testData.MetadataModel)

		res, err := fd.Filter(testData.QueryCondition, testData.RootJsonPathKey, testData.RootJsonPathToValue)

		if !reflect.DeepEqual(res, testData.FilterExcludeIndexes) {
			t.Error(
				testData.TestTitle, "\n",
				"expected res to be equal to testData.FilterExcludeIndexes\n",
				"filterExcludeIndexes=", gojsoncore.JsonStringifyMust(testData.FilterExcludeIndexes), "\n",
				"res=", gojsoncore.JsonStringifyMust(res),
			)
		}

		if err != nil && testData.LogErrorsIfExpectedNotOk {
			var filterError *core.Error
			if errors.As(err, &filterError) {
				t.Log(
					testData.TestTitle, "\n",
					"-----Error Details-----", "\n",
					filterError.String(), "\n",
					"-----------------------",
				)
			}
		}
	}
}

type filterData struct {
	internal.TestData
	Object               *object.Object
	MetadataModel        gojsoncore.JsonObject
	QueryCondition       gojsoncore.JsonObject
	RootJsonPathKey      path.JSONPath
	RootJsonPathToValue  path.JSONPath
	FilterExcludeIndexes []int
}

func filterDataTestData(yield func(data *filterData) bool) {
	obj := object.NewObject().WithSourceInterface([]*testdata.Product{
		{
			ID:    []int{0},
			Name:  []string{"Product 0"},
			Price: []float64{11.0},
		},
		{
			ID:   []int{1},
			Name: []string{"Product 1"},
		},
		{
			ID:   []int{2},
			Name: []string{"Product 2"},
		},
		{
			ID:   []int{3},
			Name: []string{"Product 3"},
		},
	})
	metadataModel := testdata.ProductMetadataModel(nil)
	testCaseIndex := 1
	if !yield(
		&filterData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d: Product Metadata Model", testCaseIndex),
			},
			Object:        obj,
			MetadataModel: metadataModel,
			QueryCondition: gojsoncore.JsonObject{
				QueryConditionType:              QuerySectionTypeLogicalOperator,
				QuerySectionTypeLogicalOperator: QuerySectionTypeLogicalOperatorAnd,
				QueryConditionValue: gojsoncore.JsonArray{
					gojsoncore.JsonObject{
						QueryConditionType: QuerySectionTypeFieldGroup,
						QueryConditionValue: gojsoncore.JsonObject{
							path.JsonpathKeyRoot + path.JsonpathDotNotation + core.GroupFields + core.ArrayPathPlaceholder + path.JsonpathDotNotation + "ID": gojsoncore.JsonObject{
								FilterConditionGreaterThan: gojsoncore.JsonObject{
									FilterConditionAssumedFieldType: core.FieldTypeNumber,
									FilterConditionValue:            0,
								},
								FilterConditionLessThan: gojsoncore.JsonObject{
									FilterConditionAssumedFieldType: core.FieldTypeNumber,
									FilterConditionValue:            3,
								},
							},
						},
					},
					gojsoncore.JsonObject{
						QueryConditionType: QuerySectionTypeFieldGroup,
						QueryConditionValue: gojsoncore.JsonObject{
							path.JsonpathKeyRoot + path.JsonpathDotNotation + core.GroupFields + core.ArrayPathPlaceholder + path.JsonpathDotNotation + "Name": gojsoncore.JsonObject{
								FilterConditionEndsWith: gojsoncore.JsonObject{
									FilterConditionAssumedFieldType: core.FieldTypeText,
									FilterConditionValue:            "2",
								},
							},
						},
					},
				},
			},
			FilterExcludeIndexes: []int{0, 1, 3},
		},
	) {
		return
	}

	obj = object.NewObject().WithSourceInterface([]*testdata.UserProfile{
		{
			Name: []string{"User 0"},
			Age:  []int{10},
			Address: []testdata.Address{
				{
					Street: []string{"Street 1"},
					City:   []string{"City 1"},
				},
			},
		},
		{
			Name: []string{"User 2"},
			Age:  []int{20},
			Address: []testdata.Address{
				{
					Street: []string{"Street 2"},
					City:   []string{"City 2"},
				},
			},
		},
		{
			Name: []string{"User 3"},
			Age:  []int{30},
			Address: []testdata.Address{
				{
					Street: []string{"Street 3"},
					City:   []string{"City 3"},
				},
				{
					Street: []string{"Street 4"},
					City:   []string{"City 4"},
				},
			},
		},
	})
	metadataModel = testdata.UserProfileMetadataModel(nil)

	testCaseIndex++
	if !yield(
		&filterData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d: User Profile Metadata Model", testCaseIndex),
			},
			Object:        obj,
			MetadataModel: metadataModel,
			QueryCondition: gojsoncore.JsonObject{
				QueryConditionType:              QuerySectionTypeLogicalOperator,
				QuerySectionTypeLogicalOperator: QuerySectionTypeLogicalOperatorOr,
				QueryConditionValue: gojsoncore.JsonArray{
					gojsoncore.JsonObject{
						QueryConditionType: QuerySectionTypeFieldGroup,
						QueryConditionValue: gojsoncore.JsonObject{
							path.JsonpathKeyRoot + path.JsonpathDotNotation + core.GroupFields + core.ArrayPathPlaceholder + path.JsonpathDotNotation + "Address": gojsoncore.JsonObject{
								FilterConditionNoOfEntriesGreaterThan: gojsoncore.JsonObject{
									FilterConditionValue: 1,
								},
							},
						},
					},
					gojsoncore.JsonObject{
						QueryConditionType: QuerySectionTypeFieldGroup,
						QueryConditionValue: gojsoncore.JsonObject{
							path.JsonpathKeyRoot + path.JsonpathDotNotation + core.GroupFields + core.ArrayPathPlaceholder + path.JsonpathDotNotation + "Address" + path.JsonpathDotNotation + core.GroupFields + core.ArrayPathPlaceholder + path.JsonpathDotNotation + "City": gojsoncore.JsonObject{
								FilterConditionBeginsWith: gojsoncore.JsonObject{
									FilterConditionAssumedFieldType: core.FieldTypeText,
									FilterConditionValue:            "Streete",
								},
							},
						},
					},
				},
			},
			FilterExcludeIndexes: []int{0, 1, 2},
		},
	) {
		return
	}

	testCaseIndex++
	if !yield(
		&filterData{
			TestData: internal.TestData{
				TestTitle: fmt.Sprintf("Test Case %d: User Profile Metadata Model with focus on Address of Profile at index 2", testCaseIndex),
			},
			Object:        obj,
			MetadataModel: metadataModel,
			QueryCondition: gojsoncore.JsonObject{
				QueryConditionType:   QuerySectionTypeFieldGroup,
				QueryConditionNegate: true,
				QueryConditionValue: gojsoncore.JsonObject{
					path.JsonpathKeyRoot + path.JsonpathDotNotation + core.GroupFields + core.ArrayPathPlaceholder + path.JsonpathDotNotation + "Address" + path.JsonpathDotNotation + core.GroupFields + core.ArrayPathPlaceholder + path.JsonpathDotNotation + "City": gojsoncore.JsonObject{
						FilterConditionBeginsWith: gojsoncore.JsonObject{
							FilterConditionAssumedFieldType: core.FieldTypeText,
							FilterConditionValue:            "City 4",
						},
					},
				},
			},
			RootJsonPathKey:      path.JSONPath(path.JsonpathKeyRoot + path.JsonpathDotNotation + core.GroupFields + core.ArrayPathPlaceholder + path.JsonpathDotNotation + "Address"),
			RootJsonPathToValue:  path.JSONPath(path.JsonpathKeyRoot + path.JsonpathLeftBracket + "2" + path.JsonpathRightBracket + path.JsonpathDotNotation + "Address"),
			FilterExcludeIndexes: []int{1},
		},
	) {
		return
	}

}
