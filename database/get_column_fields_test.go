package database

import (
	"errors"
	"reflect"
	"testing"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/internal"
	"github.com/rogonion/go-metadatamodel/iter"
	"github.com/rogonion/go-metadatamodel/testdata"
)

func TestDatabase_GetColumnFields(t *testing.T) {
	for testData := range getColumnFieldsTestData {
		gcf := NewGetColumnFields().WithSkip(testData.Skip).WithAdd(testData.Add)
		if testData.TableCollectionUid != nil {
			gcf.SetTableCollectionUID(*testData.TableCollectionUid)
		}
		if testData.JoinDepth != nil {
			gcf.SetJoinDepth(*testData.JoinDepth)
		}
		if testData.TableCollectionName != nil {
			gcf.SetTableCollectionName(*testData.TableCollectionName)
		}

		res, err := gcf.Get(testData.MetadataModel)
		if testData.ExpectedOk && err != nil {
			t.Error(
				"expected ok=", testData.ExpectedOk, "got error=", err, "\n",
				"MetadataModel=", testData.MetadataModel, "\n",
			)
			var databaseError *core.Error
			if errors.As(err, &databaseError) {
				t.Error("Test Title:", testData.TestTitle, "\n",
					"-----Error Details-----", "\n",
					databaseError.String(), "\n",
					"-----------------------",
				)
			}
		} else {
			if !reflect.DeepEqual(res, testData.Expected) {
				t.Error(
					"expected res to be equal to testData.Expected\n",
					"Test Title:", testData.TestTitle, "\n",
					"MetadataModel=", testData.MetadataModel, "\n",
					"res=", gojsoncore.JsonStringifyMust(res), "\n",
					"testData.Expected", gojsoncore.JsonStringifyMust(testData.Expected),
				)
			}
		}

		if err != nil && testData.LogErrorsIfExpectedNotOk {
			var databaseError *core.Error
			if errors.As(err, &databaseError) {
				t.Log(
					"-----Error Details-----", "\n",
					"Test Tile:", testData.TestTitle, "\n",
					databaseError.String(), "\n",
					"-----------------------",
				)
			}
		}
	}
}

type getColumnFieldsData struct {
	internal.TestData
	MetadataModel       gojsoncore.JsonObject
	JoinDepth           *int64
	TableCollectionName *string
	TableCollectionUid  *string
	Skip                core.FieldGroupPropertiesMatch
	Add                 core.FieldGroupPropertiesMatch
	ExpectedOk          bool
	Expected            *ColumnFields
}

func getColumnFieldsTestData(yield func(data *getColumnFieldsData) bool) {
	currentMetadataModel := testdata.ProductMetadataModel(nil)
	// Case 1: Product Metadata Model
	// Scenario: Extract fields for "Product" table collection using UID.
	if !yield(&getColumnFieldsData{
		TestData: internal.TestData{
			TestTitle: "Product MetadataModel (By UID)",
		},
		MetadataModel:      currentMetadataModel,
		TableCollectionUid: gojsoncore.Ptr(currentMetadataModel[core.FieldGroupName].(string)),
		ExpectedOk:         true,
		Expected: func() *ColumnFields {
			columnFields := NewColumnFields()

			iter.ForEach(currentMetadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
				if tableCollectionUid, ok := fieldGroup[core.DatabaseTableCollectionUid].(string); ok && tableCollectionUid == currentMetadataModel[core.FieldGroupName] {
					if !core.IsFieldAGroup(fieldGroup) {
						fieldColumnName := fieldGroup[core.DatabaseFieldColumnName].(string)
						columnFields.ColumnFieldsReadOrder = append(columnFields.ColumnFieldsReadOrder, fieldColumnName)
						columnFields.Fields[fieldColumnName] = fieldGroup
					}
					return false, false
				}
				return false, true
			})

			return columnFields
		}(),
	}) {
		return
	}

	currentMetadataModel = testdata.CompanyMetadataModel(nil)
	// Case 2: Company Metadata Model
	// Scenario: Extract fields for "Company" table collection using Name and JoinDepth.
	if !yield(&getColumnFieldsData{
		TestData: internal.TestData{
			TestTitle: "Company MetadataModel (By Name & JoinDepth)",
		},
		MetadataModel:       currentMetadataModel,
		JoinDepth:           gojsoncore.Ptr(int64(0)),
		TableCollectionName: gojsoncore.Ptr(currentMetadataModel[core.FieldGroupName].(string)),
		ExpectedOk:          true,
		Expected: func() *ColumnFields {
			columnFields := NewColumnFields()

			iter.ForEach(currentMetadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
				if joinDepth, ok := fieldGroup[core.DatabaseJoinDepth].(float64); ok && joinDepth == float64(0) {
					if tableCollectionName, ok := fieldGroup[core.DatabaseTableCollectionName].(string); ok && tableCollectionName == currentMetadataModel[core.FieldGroupName] {
						if !core.IsFieldAGroup(fieldGroup) {
							fieldColumnName := fieldGroup[core.DatabaseFieldColumnName].(string)
							columnFields.ColumnFieldsReadOrder = append(columnFields.ColumnFieldsReadOrder, fieldColumnName)
							columnFields.Fields[fieldColumnName] = fieldGroup
						}
						return false, false
					}
				}
				return false, true
			})

			return columnFields
		}(),
	}) {
		return
	}

	currentMetadataModel = testdata.UserProfileMetadataModel(nil)
	// Case 3: UserProfile Metadata Model
	// Scenario: Extract fields for "UserProfile" table collection using Name and JoinDepth.
	if !yield(&getColumnFieldsData{
		TestData: internal.TestData{
			TestTitle: "UserProfile MetadataModel (By Name & JoinDepth)",
		},
		MetadataModel:       currentMetadataModel,
		JoinDepth:           gojsoncore.Ptr(int64(0)),
		TableCollectionName: gojsoncore.Ptr(currentMetadataModel[core.FieldGroupName].(string)),
		ExpectedOk:          true,
		Expected: func() *ColumnFields {
			columnFields := NewColumnFields()

			iter.ForEach(currentMetadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
				if joinDepth, ok := fieldGroup[core.DatabaseJoinDepth].(float64); ok && joinDepth == float64(0) {
					if tableCollectionName, ok := fieldGroup[core.DatabaseTableCollectionName].(string); ok && tableCollectionName == currentMetadataModel[core.FieldGroupName] {
						if !core.IsFieldAGroup(fieldGroup) {
							fieldColumnName := fieldGroup[core.DatabaseFieldColumnName].(string)
							columnFields.ColumnFieldsReadOrder = append(columnFields.ColumnFieldsReadOrder, fieldColumnName)
							columnFields.Fields[fieldColumnName] = fieldGroup
						}
						return false, false
					}
				}
				return false, true
			})

			return columnFields
		}(),
	}) {
		return
	}

	currentMetadataModel = testdata.EmployeeMetadataModel(nil)
	// Case 4: Employee Metadata Model (Profile)
	// Scenario: Extract fields for "Profile" table collection within Employee model using UID.
	if !yield(&getColumnFieldsData{
		TestData: internal.TestData{
			TestTitle: "Employee MetadataModel - Profile Collection (By UID)",
		},
		MetadataModel:      currentMetadataModel,
		TableCollectionUid: gojsoncore.Ptr("Profile"),
		ExpectedOk:         true,
		Expected: func() *ColumnFields {
			columnFields := NewColumnFields()

			iter.ForEach(currentMetadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
				if tableCollectionUid, ok := fieldGroup[core.DatabaseTableCollectionUid].(string); ok && tableCollectionUid == "Profile" {
					if !core.IsFieldAGroup(fieldGroup) {
						fieldColumnName := fieldGroup[core.DatabaseFieldColumnName].(string)
						columnFields.ColumnFieldsReadOrder = append(columnFields.ColumnFieldsReadOrder, fieldColumnName)
						columnFields.Fields[fieldColumnName] = fieldGroup
					}
					return false, false
				}
				return false, true
			})

			return columnFields
		}(),
	}) {
		return
	}

	// Case 5: Employee Metadata Model (Profile) with Add Filter
	// Scenario: Extract fields for "Profile" table collection but ONLY include "Age" field.
	if !yield(&getColumnFieldsData{
		TestData: internal.TestData{
			TestTitle: "Employee MetadataModel - Profile Collection (By UID) with Add Filter (Age)",
		},
		MetadataModel:      currentMetadataModel,
		TableCollectionUid: gojsoncore.Ptr("Profile"),
		ExpectedOk:         true,
		Add: core.FieldGroupPropertiesMatch{
			core.DatabaseFieldColumnName: "Age",
		},
		Expected: func() *ColumnFields {
			columnFields := NewColumnFields()

			iter.ForEach(currentMetadataModel, func(fieldGroup gojsoncore.JsonObject) (bool, bool) {
				if tableCollectionUid, ok := fieldGroup[core.DatabaseTableCollectionUid].(string); ok && tableCollectionUid == "Profile" {
					if !core.IsFieldAGroup(fieldGroup) {
						if fieldGroup[core.DatabaseFieldColumnName] == "Age" {
							fieldColumnName := fieldGroup[core.DatabaseFieldColumnName].(string)
							columnFields.ColumnFieldsReadOrder = append(columnFields.ColumnFieldsReadOrder, fieldColumnName)
							columnFields.Fields[fieldColumnName] = fieldGroup
						}
					}
					return false, false
				}
				return false, true
			})

			return columnFields
		}(),
	}) {
		return
	}
}
