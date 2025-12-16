package iter

import (
	"fmt"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/internal"
)

func TestIter_ForEach(t *testing.T) {
	for testData := range forEachTestData {
		noOfIterations := 0
		ForEach(testData.MetadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
			noOfIterations++
			return false, false
		})

		if noOfIterations != testData.Expected {
			t.Error(
				testData.TestTitle, "\n",
				"expected noOfIterations to be equal to Expected",
				"testData.Expected=", testData.Expected, "\n",
				"noOfIterations=", noOfIterations,
				"testData.MetadataModel=", gojsoncore.JsonStringifyMust(testData.MetadataModel), "\n",
			)
		}
	}
}

type forEachData struct {
	internal.TestData
	MetadataModel gojsoncore.JsonObject
	Expected      int
}

func forEachTestData(yield func(data *forEachData) bool) {
	testCaseIndex := 1
	if !yield(&forEachData{
		TestData: internal.TestData{
			TestTitle: fmt.Sprintf("Test Case %d: Test Number of iterations from UserInformationMetadataModel", testCaseIndex),
		},
		MetadataModel: UserInformationMetadataModel(),
		Expected:      5,
	}) {
		return
	}

}
