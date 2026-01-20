package flattener

import (
	"reflect"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/fieldcolumns"
	"github.com/rogonion/go-metadatamodel/iter"
	"github.com/rogonion/go-metadatamodel/testdata"
)

func TestFlattener_Flatten_And_WriteTo(t *testing.T) {
	for data := range flattenerTestData {
		// 1. Prepare Destination (A 2D Slice holder)
		// We initialize it as empty [][]any so the Flattener writes rows/cols into it.
		destinationSource := make([][]any, 0)
		destination := object.NewObject().WithSourceInterface(destinationSource)

		// 3. Initialize Flattener
		f := NewFlattener(data.MetadataModel).
			WithColumnFields(data.ColumnFields) // Optional: For reordering/skipping

		// 4. Run Flatten (Process)
		err := f.Flatten(data.SourceObject)
		if err != nil {
			t.Errorf("%s: Flatten() unexpected error: %v", data.TestTitle, err)
			continue
		}

		// 5. Run WriteToDestination (Output)
		err = f.WriteToDestination(destination)
		if err != nil {
			t.Errorf("%s: WriteToDestination() unexpected error: %v", data.TestTitle, err)
			continue
		}

		// 6. Validate Result
		// We grab the actual data sitting in the destination object
		actualResult := destination.GetSourceInterface()

		if !reflect.DeepEqual(actualResult, data.ExpectedResult) {
			t.Errorf("%s: Result mismatch.\nExpected:\n%#v\nGot:\n%#v",
				data.TestTitle,
				data.ExpectedResult,
				actualResult,
			)
		}
	}
}

// --- Test Data Structures ---

type flattenTestData struct {
	TestTitle      string
	SourceObject   *object.Object
	MetadataModel  gojsoncore.JsonObject
	ColumnFields   *fieldcolumns.ColumnFields
	ExpectedResult [][]any
}

func flattenerTestData(yield func(data *flattenTestData) bool) {
	// -------------------------------------------------------------------------
	// Case 1: Simple Flat Object (User)
	// -------------------------------------------------------------------------
	userMeta := testdata.UserMetadataModel(nil)
	userData := testdata.User{
		ID:    []int{101},
		Name:  []string{"Alice"},
		Email: []string{"alice@example.com"},
	}
	// Default Read Order: ID, Name, Email
	expectedUser := [][]any{
		{
			[]int{101},
			[]string{"Alice"},
			[]string{"alice@example.com"},
		},
	}

	if !yield(&flattenTestData{
		TestTitle:      "Simple User Flattening",
		SourceObject:   object.NewObject().WithSourceInterface(userData),
		MetadataModel:  userMeta,
		ExpectedResult: expectedUser,
	}) {
		return
	}

	// -------------------------------------------------------------------------
	// Case 2: Deep Nested Object (Employee -> Profile -> Address)
	// -------------------------------------------------------------------------
	empMeta := testdata.EmployeeMetadataModel(nil)
	empData := testdata.Employee{
		ID:     []int{500},
		Skills: []string{"Go", "Rust"}, // Leaf collection
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
				},
			},
		},
	}

	// Expected Columns based on EmployeeMetadataModel default read order:
	// 1. ID
	// 2. Profile.Name
	// 3. Profile.Age
	// 4. Profile.Address.Street
	// 5. Profile.Address.City
	// 6. Profile.Address.ZipCode
	// 7. Skills
	expectedEmp := [][]any{
		{
			[]int{500},
			[]string{"Bob"},
			[]int{30},
			[]string{"123 Tech Ln"},
			[]string{"Silicon Valley"},
			[]*string{gojsoncore.Ptr("94000")},
			[]string{"Go", "Rust"}, // Note: Skills is a slice in a single cell
		},
	}

	if !yield(&flattenTestData{
		TestTitle:      "Deep Nested Employee Flattening",
		SourceObject:   object.NewObject().WithSourceInterface(empData),
		MetadataModel:  empMeta,
		ExpectedResult: expectedEmp,
	}) {
		return
	}

	// -------------------------------------------------------------------------
	// Case 3: Linear Collection (Array of Products)
	// -------------------------------------------------------------------------
	prodMeta := testdata.ProductMetadataModel(nil)
	prodData := []testdata.Product{
		{ID: []int{1}, Name: []string{"Laptop"}, Price: []float64{999.99}},
		{ID: []int{2}, Name: []string{"Mouse"}, Price: []float64{25.50}},
	}
	// Expected: 2 Rows
	expectedProd := [][]any{
		{[]int{1}, []string{"Laptop"}, []float64{999.99}},
		{[]int{2}, []string{"Mouse"}, []float64{25.50}},
	}

	if !yield(&flattenTestData{
		TestTitle:      "Array of Products Flattening",
		SourceObject:   object.NewObject().WithSourceInterface(prodData),
		MetadataModel:  prodMeta,
		ExpectedResult: expectedProd,
	}) {
		return
	}

	// -------------------------------------------------------------------------
	// Case 4: WriteTo with Reordering (Using ColumnFields)
	// -------------------------------------------------------------------------
	// We reuse the Product data but force a specific column order: Price, Name (ID skipped)

	// Manually construct ColumnFields to simulate an Extraction result
	// Original Order: ID (0), Name (1), Price (2)
	// Desired Order: Price (2), Name (1)
	reorderedCols := fieldcolumns.NewColumnFields()
	reorderedCols.CurrentIndexOfReadOrderOfColumnFields = []int{2, 1} // Index in Original Read Order

	// We must populate Fields map so WriteTo can validate
	// (In a real scenario, Extraction fills this, here we mock what WriteTo needs)
	// Actually, WriteTo relies on GetCurrentIndexOfReadOrderOfFields which uses Fields to check skips.
	// Let's rely on the fact that WriteTo uses n.currentReadOrderOfColumnFields directly if set
	// via columnFields.GetCurrentIndexOfReadOrderOfFields().

	// Since we can't easily mock the full internal structure of ColumnFields without Extraction,
	// we will perform a real extraction first using the user's library logic, then reposition.

	prodSchema := testdata.ProductSchema()
	extraction := fieldcolumns.NewColumnFieldsExtraction(prodMeta).WithSchema(prodSchema)
	extractedCols, _ := extraction.Extract()

	// Reposition: Swap ID(0) and Price(2), essentially.
	// Let's just manually override the Index Slice for simplicity in this test
	// We want Price (index 2 in original) then Name (index 1 in original). ID (0) is omitted.
	extractedCols.CurrentIndexOfReadOrderOfColumnFields = []int{2, 1}

	expectedReordered := [][]any{
		{[]float64{999.99}, []string{"Laptop"}},
		{[]float64{25.50}, []string{"Mouse"}},
	}

	if !yield(&flattenTestData{
		TestTitle:      "WriteTo with Reordering and Skipping",
		SourceObject:   object.NewObject().WithSourceInterface(prodData),
		MetadataModel:  prodMeta,
		ColumnFields:   extractedCols,
		ExpectedResult: expectedReordered,
	}) {
		return
	}

	// -------------------------------------------------------------------------
	// Case 5: Matrix Multiplication / Cartesian Product
	// Multiple Employees -> Multiple Profiles -> Multiple Addresses
	// -------------------------------------------------------------------------
	// Scenario:
	// Employee 1:
	//   - Profile A ("Dev"):
	//       - Address 1 ("Home")
	//       - Address 2 ("Work")  <-- Branching happens here (1 Prof * 2 Addr)
	//
	// Employee 2:
	//   - Profile B ("Admin"):
	//       - Address 3 ("Office")
	//   - Profile C ("Consultant"): <-- Iteration happens here (Prof B, then Prof C)
	//       - Address 4 ("Remote")
	//
	// Expected Total Rows: 4
	// Row 1: Emp1, Dev, Home
	// Row 2: Emp1, Dev, Work
	// Row 3: Emp2, Admin, Office
	// Row 4: Emp2, Consultant, Remote

	complexEmpData := []testdata.Employee{
		{
			ID:     []int{100},
			Skills: []string{"Go"},
			Profile: []*testdata.UserProfile{
				{
					Name: []string{"Dev"},
					Age:  []int{30},
					Address: []testdata.Address{
						{
							Street:  []string{"Home St"},
							City:    []string{"Nairobi"},
							ZipCode: []*string{gojsoncore.Ptr("00100")},
						},
						{
							Street:  []string{"Work Ave"},
							City:    []string{"Westlands"},
							ZipCode: []*string{gojsoncore.Ptr("00200")},
						},
					},
				},
			},
		},
		{
			ID:     []int{200},
			Skills: []string{"Management"},
			Profile: []*testdata.UserProfile{
				{
					Name: []string{"Admin"},
					Age:  []int{45},
					Address: []testdata.Address{
						{
							Street:  []string{"HQ Blvd"},
							City:    []string{"Mombasa"},
							ZipCode: []*string{gojsoncore.Ptr("80100")},
						},
					},
				},
				{
					Name: []string{"Consultant"},
					Age:  []int{50},
					Address: []testdata.Address{
						{
							Street:  []string{"Remote Ln"},
							City:    []string{"Kisumu"},
							ZipCode: []*string{gojsoncore.Ptr("40100")},
						},
					},
				},
			},
		},
	}

	// Based on EmployeeMetadataModel, the read order is:
	// 1. ID
	// 2. Profile.Name
	// 3. Profile.Age
	// 4. Profile.Address.Street
	// 5. Profile.Address.City
	// 6. Profile.Address.ZipCode
	// 7. Skills
	expectedComplexRows := [][]any{
		// -- Employee 1, Profile A, Address 1 --
		{
			[]int{100},
			[]string{"Dev"},
			[]int{30},
			[]string{"Home St"},
			[]string{"Nairobi"},
			[]*string{gojsoncore.Ptr("00100")},
			[]string{"Go"},
		},
		// -- Employee 1, Profile A, Address 2 (Note: Parent fields 'Dev'/'100' repeated) --
		{
			[]int{100},
			[]string{"Dev"},
			[]int{30},
			[]string{"Work Ave"},
			[]string{"Westlands"},
			[]*string{gojsoncore.Ptr("00200")},
			[]string{"Go"},
		},
		// -- Employee 2, Profile B, Address 3 --
		{
			[]int{200},
			[]string{"Admin"},
			[]int{45},
			[]string{"HQ Blvd"},
			[]string{"Mombasa"},
			[]*string{gojsoncore.Ptr("80100")},
			[]string{"Management"},
		},
		// -- Employee 2, Profile C, Address 4 --
		{
			[]int{200},
			[]string{"Consultant"},
			[]int{50},
			[]string{"Remote Ln"},
			[]string{"Kisumu"},
			[]*string{gojsoncore.Ptr("40100")},
			[]string{"Management"},
		},
	}

	if !yield(&flattenTestData{
		TestTitle:      "Matrix Multiplication: Multi-Employee/Profile/Address",
		SourceObject:   object.NewObject().WithSourceInterface(complexEmpData),
		MetadataModel:  empMeta, // Reusing metadata from Case 2
		ExpectedResult: expectedComplexRows,
	}) {
		return
	}

	// -------------------------------------------------------------------------
	// Case 6: Horizontal Expansion (ViewInSeparateColumns)
	// Scenario:
	// User wants Skills (Max 2) and Addresses (Max 2) to be columns, not rows.
	// -------------------------------------------------------------------------

	// 1. Prepare Data
	// One Employee, 2 Skills, 2 Addresses.
	// Normally this would produce 2x2 = 4 rows.
	// With Horizontal Expansion, this should produce 1 Row.
	pivotedEmpData := testdata.Employee{
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
	}

	// 2. Prepare Metadata with Horizontal Expansion
	// We iterate over the default model and inject the properties.
	pivotedMeta := testdata.EmployeeMetadataModel(nil)

	// Helper to find and update nodes
	// Logic:
	// - Find "Skills" -> Set MaxCols = 3 (Even though data only has 2, ensures padding works)
	// - Find "Address" -> Set MaxCols = 2
	pivotedMeta = iter.Map(pivotedMeta, func(node gojsoncore.JsonObject) (any, bool) {
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

	// 3. Define Expected Result (Single Wide Row)
	// The columns will follow the order:
	// ID, Profile.Name, Profile.Age,
	// Address[0].Street, Address[0].City, Address[0].Zip,
	// Address[1].Street, Address[1].City, Address[1].Zip,
	// Skills[0], Skills[1], Skills[2] (Empty/Padding)

	expectedPivotedRow := [][]any{
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

	if !yield(&flattenTestData{
		TestTitle:      "Horizontal Expansion (Pivoting)",
		SourceObject:   object.NewObject().WithSourceInterface(pivotedEmpData),
		MetadataModel:  pivotedMeta,
		ExpectedResult: expectedPivotedRow,
	}) {
		return
	}
}
