package unflattener

import (
	"reflect"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/fieldcolumns"
	"github.com/rogonion/go-metadatamodel/flattener"
	"github.com/rogonion/go-metadatamodel/iter"
	"github.com/rogonion/go-metadatamodel/testdata"
)

func TestUnflattener_Unflatten(t *testing.T) {
	for data := range unflattenerTestData {
		t.Run(data.TestTitle, func(t *testing.T) {
			// 1. Prepare Source (FlattenedTable)
			sourceTable := toFlattenedTable(data.SourceTable)

			// 2. Prepare Destination
			destination := object.NewObject()
			if data.Schema == nil {
				destination = destination.WithSourceInterface(reflect.New(reflect.TypeOf(data.ExpectedResult)).Elem().Interface())
			} else {
				destination = destination.WithSchema(data.Schema)
			}

			// 3. Initialize Unflattener
			u := NewUnflattener(data.MetadataModel, NewSignature())

			if data.ColumnFields != nil {
				u.WithColumnFields(data.ColumnFields)
			}

			u.WithDestination(destination)

			// 4. Run Unflatten
			err := u.Unflatten(sourceTable)
			if err != nil {
				t.Fatalf("Unflatten() unexpected error: %v", err)
			}

			// 5. Validate Result
			actualResult := destination.GetSourceInterface()
			if !reflect.DeepEqual(actualResult, data.ExpectedResult) {
				t.Errorf("Result mismatch.\nExpected:\n%#v\nGot:\n%#v",
					data.ExpectedResult,
					actualResult,
				)
			}
		})
	}
}

func toFlattenedTable(data [][]any) flattener.FlattenedTable {
	table := make(flattener.FlattenedTable, len(data))
	for i, row := range data {
		fRow := make(flattener.FlattenedRow, len(row))
		for j, val := range row {
			if val == nil {
				// Represent nil as invalid value
				fRow[j] = reflect.Value{}
			} else {
				fRow[j] = reflect.ValueOf(val)
			}
		}
		table[i] = fRow
	}
	return table
}

type unflattenTestData struct {
	TestTitle      string
	SourceTable    [][]any
	Schema         schema.Schema
	MetadataModel  gojsoncore.JsonObject
	ColumnFields   *fieldcolumns.ColumnFields
	ExpectedResult any
}

func unflattenerTestData(yield func(data *unflattenTestData) bool) {
	// -------------------------------------------------------------------------
	// Case 1: Simple Flat Object (User)
	// -------------------------------------------------------------------------
	userMeta := testdata.UserMetadataModel(nil)
	// Unflattener always produces a collection (slice) of objects
	expectedUser := []*testdata.User{{
		ID:    []int{101},
		Name:  []string{"Alice"},
		Email: []string{"alice@example.com"},
	}}
	// Flattened: ID, Name, Email
	userTable := [][]any{
		{
			[]int{101},
			[]string{"Alice"},
			[]string{"alice@example.com"},
		},
	}

	if !yield(
		&unflattenTestData{
			TestTitle:      "Simple User Unflattening",
			SourceTable:    userTable,
			MetadataModel:  userMeta,
			ExpectedResult: expectedUser,
		},
	) {
		return
	}

	// -------------------------------------------------------------------------
	// Case 2: Deep Nested Object (Employee -> Profile -> Address)
	// -------------------------------------------------------------------------
	empMeta := testdata.EmployeeMetadataModel(nil)
	expectedEmp := []*testdata.Employee{
		{
			ID:     []int{500},
			Skills: []string{"Go", "Rust"},
			Profile: []*testdata.UserProfile{
				{
					Name: []string{"Bob"},
					Age:  []int{30},
					Address: []testdata.Address{
						{
							Street:  []string{"123 Tech Ln"},
							City:    []string{"Silicon Valley"},
							ZipCode: []*string{gojsoncore.Ptr("94000")},
						},
						{
							Street:  []string{"456 Tech Ln"},
							City:    []string{"Silicon Valley"},
							ZipCode: []*string{gojsoncore.Ptr("94000")},
						},
					},
				},
				{
					Name:    []string{"Alice"},
					Age:     []int{100},
					Address: nil,
				},
			},
		},
		{
			ID:     []int{600},
			Skills: []string{"HTML", "CSS"},
			Profile: []*testdata.UserProfile{
				{
					Name: []string{"Doe"},
					Age:  []int{35},
					Address: []testdata.Address{
						{
							Street:  []string{"123 Tech Ln"},
							City:    []string{"Silicon Valley"},
							ZipCode: []*string{gojsoncore.Ptr("94000")},
						},
					},
				},
			},
		},
	}
	// Flattened Columns: ID, Profile.Name, Profile.Age, Profile.Address.Street, Profile.Address.City, Profile.Address.ZipCode, Skills
	empTable := [][]any{
		{
			[]int{500},
			[]string{"Bob"},
			[]int{30},
			[]string{"123 Tech Ln"},
			[]string{"Silicon Valley"},
			[]string{"94000"},
			[]string{"Go", "Rust"},
		},
		{
			[]int{500},
			[]string{"Bob"},
			[]int{30},
			[]string{"456 Tech Ln"},
			[]string{"Silicon Valley"},
			[]string{"94000"},
			[]string{"Go", "Rust"},
		},
		{
			[]int{500},
			[]string{"Alice"},
			[]int{100},
			nil,
			nil,
			nil,
			nil,
		},
		{
			600,
			"Doe",
			"35",
			"123 Tech Ln",
			"Silicon Valley",
			"94000",
			[]string{"HTML", "CSS"},
		},
		{
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}

	if !yield(&unflattenTestData{
		TestTitle:   "Deep Nested Employee Unflattening",
		SourceTable: empTable,
		Schema: &schema.DynamicSchemaNode{
			Type: reflect.TypeOf([]*testdata.Employee{}),
			Kind: reflect.Slice,
			ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
				Type:                    reflect.TypeOf(&testdata.Employee{}),
				Kind:                    reflect.Pointer,
				ChildNodesPointerSchema: testdata.EmployeeSchema(),
			},
		},
		MetadataModel:  empMeta,
		ExpectedResult: expectedEmp,
	}) {
		return
	}

	// -------------------------------------------------------------------------
	// Case 3: Linear Collection (Array of Products)
	// -------------------------------------------------------------------------
	// We need to ensure ID is a Key for grouping to work on the array.
	expectedProd := []*testdata.Product{
		{ID: []int{1}, Name: []string{"Laptop"}, Price: []float64{999.99}},
		{ID: []int{2}, Name: []string{"Mouse"}, Price: []float64{25.50}},
	}
	prodTable := [][]any{
		{[]int{1}, []string{"Laptop"}, []float64{999.99}},
		{[]int{2}, []string{"Mouse"}, []float64{25.50}},
	}

	if !yield(&unflattenTestData{
		TestTitle:      "Array of Products Unflattening",
		SourceTable:    prodTable,
		MetadataModel:  testdata.ProductMetadataModel(nil),
		ExpectedResult: expectedProd,
	}) {
		return
	}

	// -------------------------------------------------------------------------
	// Case 4: Empty Rows / Partial Data
	// -------------------------------------------------------------------------
	// Product with missing price (nil)

	expectedPartialProd := []*testdata.Product{
		{ID: []int{1}, Name: []string{"Laptop"}, Price: []float64{999.99}},
		{ID: []int{2}, Name: []string{"Freebie"}, Price: nil}, // Price missing
	}
	// Note: Price is nil in the source table
	partialProdTable := [][]any{
		{[]int{1}, []string{"Laptop"}, []float64{999.99}},
		{[]int{2}, []string{"Freebie"}, nil},
	}

	if !yield(&unflattenTestData{
		TestTitle:      "Partial Data (Nil Values)",
		SourceTable:    partialProdTable,
		MetadataModel:  testdata.ProductMetadataModel(nil),
		ExpectedResult: expectedPartialProd,
	}) {
		return
	}

	// -------------------------------------------------------------------------
	// Case 5: Horizontal Expansion (Pivoting)
	// -------------------------------------------------------------------------
	pivotedRows := [][]any{
		{
			[]int{999},
			[]string{"Pivot Master"},
			[]int{40},

			// Address 1
			[]string{"Office 1"},
			[]string{"Nairobi"},
			[]*string{gojsoncore.Ptr("001")},

			// Address 2
			[]string{"Home 2"},
			[]string{"Mombasa"},
			[]*string{gojsoncore.Ptr("002")},

			// Skills (Max 3)
			[]string{"Go"},
			[]string{"Python"},
			nil, // Padding for 3rd column (Slice logic returns empty or nil slice)
		},
	}

	expectedPivotedEmpData := []*testdata.Employee{
		{
			ID:     []int{999},
			Skills: []string{"Go", "Python"},
			Profile: []*testdata.UserProfile{
				{
					Name: []string{"Pivot Master"},
					Age:  []int{40},
					Address: []testdata.Address{
						{
							Street:  []string{"Office 1"},
							City:    []string{"Nairobi"},
							ZipCode: []*string{gojsoncore.Ptr("001")},
						},
						{
							Street:  []string{"Home 2"},
							City:    []string{"Mombasa"},
							ZipCode: []*string{gojsoncore.Ptr("002")},
						},
					},
				},
			},
		},
	}
	// Helper to find and update nodes
	// Logic:
	// - Find "Skills" -> Set MaxCols = 3 (Even though data only has 2, ensures padding works)
	// - Find "Address" -> Set MaxCols = 2
	pivotedMeta := iter.Map(testdata.EmployeeMetadataModel(nil), func(node gojsoncore.JsonObject) (any, bool) {
		if name, ok := node[core.FieldGroupName].(string); ok {
			switch name {
			case "Skills":
				node[core.FieldGroupViewValuesInSeparateColumns] = true
				node[core.FieldGroupViewMaxNoOfValuesInSeparateColumns] = 3
			case "Address":
				node[core.FieldGroupViewValuesInSeparateColumns] = true
				node[core.FieldGroupViewMaxNoOfValuesInSeparateColumns] = 2
			}
		}
		return node, false
	}).(gojsoncore.JsonObject)

	if !yield(&unflattenTestData{
		TestTitle:      "Horizontal Expansion (Pivoting)",
		SourceTable:    pivotedRows,
		MetadataModel:  pivotedMeta,
		ExpectedResult: expectedPivotedEmpData,
	}) {
		return
	}
}
