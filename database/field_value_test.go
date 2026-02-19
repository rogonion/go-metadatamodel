package database

import (
	"reflect"
	"testing"

	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-metadatamodel/testdata"
)

// TestDatabase_FieldValueOnEmployee tests setting a value on a nested field (Profile.Age) in the Employee struct.
func TestDatabase_FieldValueOnEmployee(t *testing.T) {
	employee := &testdata.Employee{
		ID: []int{1},
	}

	obj := object.NewObject().WithSourceInterface(employee).WithSchema(testdata.EmployeeSchema())
	metadataModel := testdata.EmployeeMetadataModel(nil)
	userProfileColumnFields, err := NewGetColumnFields().WithJoinDepth(1).WithTableCollectionName("Profile").Get(metadataModel)
	if err != nil {
		t.Fatal("Get Column Fields failed:", err)
	}

	fieldValue := NewFieldValue(obj, userProfileColumnFields)
	noOfModifications, err := fieldValue.Set("Age", 16, "", nil)
	if noOfModifications != 1 {
		t.Fatal(
			"Set Age on User Profile does not have expected number of modifications", noOfModifications, "\n",
			"err=", err,
		)
	}

	if !reflect.DeepEqual(employee.Profile[0].Age, []int{16}) {
		t.Fatal(
			"Value inserted using Set not equal to res\n",
			"res=", employee, "\n",
			"expected=", []int{16},
		)
	}
}

// TestDatabase_FieldValueOnProduct tests Get, Set, and Delete operations on the Product struct.
func TestDatabase_FieldValueOnProduct(t *testing.T) {
	product := &testdata.Product{
		ID: []int{1},
	}

	obj := object.NewObject().WithSourceInterface(product).WithSchema(testdata.ProductSchema())
	metadataModel := testdata.ProductMetadataModel(nil)
	columnFields, err := NewGetColumnFields().WithJoinDepth(0).WithTableCollectionName("Product").Get(metadataModel)
	if err != nil {
		t.Fatal("Get Column Fields failed:", err)
	}

	fieldValue := NewFieldValue(obj, columnFields)
	noOfResults, err := fieldValue.Get("ID", "", nil)
	if noOfResults != 1 {
		t.Fatal(
			"Get Field Value failed:", err, "\n",
			"metadata model=", metadataModel, "\n",
			"column fields=", columnFields,
		)
	}

	valueFound := fieldValue.GetValueFoundInterface()
	if !reflect.DeepEqual(valueFound, product.ID) {
		t.Fatal(
			"Value retrieved using Get not equal to res\n",
			"res=", valueFound, "\n",
			"expected=", []int{1},
		)
	}

	noOfModifications, err := fieldValue.Set("Name", "Twinkies", "", nil)
	if noOfModifications != 1 {
		t.Fatal(
			"Update Name on Product does not have expected number of modifications", noOfModifications, "\n",
			"err=", err,
		)
	}

	if !reflect.DeepEqual(product.Name, []string{"Twinkies"}) {
		t.Fatal(
			"Value inserted using Set not equal to res\n",
			"res=", product.Name, "\n",
			"expected=", []string{"Twinkies"},
		)
	}

	product.Price = []float64{0.0}

	noOfModifications, err = fieldValue.Delete("Price", "", nil)
	if noOfModifications != 1 {
		t.Fatal(
			"Delete Price on Product does not have expected number of modifications", noOfModifications, "\n",
			"err=", err,
		)
	}

	if len(product.Price) != 0 {
		t.Fatal("Delete Price on Product does not have expected number of modifications")
	}
}
