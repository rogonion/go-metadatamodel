package fieldcolumns

import (
	"reflect"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/internal"
	"github.com/rogonion/go-metadatamodel/iter"
	"github.com/rogonion/go-metadatamodel/testdata"
)

func TestGroupsColumnsIndexesRetrieval_Get(t *testing.T) {
	for testData := range groupsColumnsIndexesRetrievalTestData {
		// 1. Extract ColumnFields first (dependency)
		fcExtraction := NewColumnFieldsExtraction(testData.MetadataModel)
		columnFields, err := fcExtraction.Extract()
		if err != nil {
			t.Fatalf("%s: Setup failed - ColumnFields extraction error: %v", testData.TestTitle, err)
		}
		columnFields.Reposition()
		columnFields.Skip(nil, nil)

		// 2. Run Retrieval
		retriever := NewGroupsColumnsIndexesRetrieval(columnFields)
		indexes, err := retriever.Get(testData.MetadataModel)

		// 3. Validate Error
		if testData.ExpectedOk && err != nil {
			t.Errorf("%s: Expected success, got error: %v", testData.TestTitle, err)
		} else if !testData.ExpectedOk && err == nil {
			t.Errorf("%s: Expected error, got success", testData.TestTitle)
		}

		// 4. Validate Result
		if testData.ExpectedOk {
			if !reflect.DeepEqual(indexes, testData.ExpectedIndexes) {
				t.Errorf("%s: Result mismatch.\nExpected:\n%#v\nGot:\n%#v",
					testData.TestTitle,
					testData.ExpectedIndexes,
					indexes,
				)
			}
		}
	}
}

type groupsColumnsIndexesRetrievalData struct {
	internal.TestData
	MetadataModel   gojsoncore.JsonObject
	ExpectedOk      bool
	ExpectedIndexes *GroupColumnIndexes
}

func groupsColumnsIndexesRetrievalTestData(yield func(data *groupsColumnsIndexesRetrievalData) bool) {
	// Case 1: Product metadata model
	if !yield(&groupsColumnsIndexesRetrievalData{
		TestData: internal.TestData{
			TestTitle: "Product",
		},
		MetadataModel: testdata.ProductMetadataModel(nil),
		ExpectedOk:    true,
		ExpectedIndexes: &GroupColumnIndexes{
			Primary: []int{0},
			All:     []int{0, 1, 2},
		},
	}) {
		return
	}

	// Case 2: Company metadata model
	if !yield(&groupsColumnsIndexesRetrievalData{
		TestData: internal.TestData{
			TestTitle: "Company",
		},
		MetadataModel: testdata.CompanyMetadataModel(nil),
		ExpectedOk:    true,
		ExpectedIndexes: &GroupColumnIndexes{
			Primary: []int{0},
			All:     []int{0},
		},
	}) {
		return
	}

	employeeMetadataModel := iter.Map(testdata.EmployeeMetadataModel(nil), func(node gojsoncore.JsonObject) (any, bool) {
		if name, ok := node[core.FieldGroupName].(string); ok && name == "UserProfile" {
			node[core.FieldGroupIsPrimaryKey] = true
		}
		return node, false
	}).(gojsoncore.JsonObject)

	// Case 3: Employee metadata model
	if !yield(&groupsColumnsIndexesRetrievalData{
		TestData: internal.TestData{
			TestTitle: "Employee with groups as primary keys",
		},
		MetadataModel: employeeMetadataModel,
		ExpectedOk:    true,
		ExpectedIndexes: &GroupColumnIndexes{
			Primary: []int{0, 1},
			All:     []int{0, 6},
		},
	}) {
		return
	}
}
