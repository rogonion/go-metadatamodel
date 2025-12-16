package core

/*
Metadata Model Field/Group Properties.
*/
const (
	FieldGroupJsonPathKey         string = "FieldGroupJsonPathKey"
	FieldGroupName                string = "FieldGroupName"
	FieldGroupDescription         string = "FieldGroupDescription"
	FieldGroupViewTableLockColumn string = "FieldGroupViewTableLockColumn"
	FieldGroupIsPrimaryKey        string = "FieldGroupIsPrimaryKey"

	// FieldGroupViewValuesInSeparateColumns For when you want to view the individual values in an array/slice of a field or non-nested group in separate columns in a flat/table.
	//
	// By default:
	//	- Fields with multiple values will be joined as a string separated by a comma or FieldMultipleValuesJoinSymbol.
	//  - Groups will appear as separate rows in a matrix.
	FieldGroupViewValuesInSeparateColumns        string = "FieldGroupViewValuesInSeparateColumns"
	FieldGroupViewMaxNoOfValuesInSeparateColumns string = "FieldGroupViewMaxNoOfValuesInSeparateColumns"
	FieldViewValuesInSeparateColumnsHeaderFormat string = "FieldViewValuesInSeparateColumnsHeaderFormat"
	FieldViewValuesInSeparateColumnsHeaderIndex  string = "FieldViewValuesInSeparateColumnsHeaderIndex"
	FieldMultipleValuesJoinSymbol                string = "FieldMultipleValuesJoinSymbol"

	FieldGroupInputDisable               string = "FieldGroupInputDisable"
	FieldGroupDisablePropertiesEdit      string = "FieldGroupDisablePropertiesEdit"
	FieldGroupViewDisable                string = "FieldGroupViewDisable"
	FieldGroupQueryConditionsEditDisable string = "FieldGroupQueryConditionsEditDisable"
	FieldGroupMaxEntries                 string = "FieldGroupMaxEntries"

	FieldDataType string = "FieldDataType"
	FieldUI       string = "FieldUi"

	// FieldColumnPosition for setting a custom position for field/group when working with data in a flat2D
	FieldColumnPosition      string = "FieldColumnPosition"
	FieldGroupPositionBefore string = "FieldGroupPositionBefore"

	FieldGroupTypeAny               string = "FieldGroupTypeAny"
	FieldCheckboxValueIfTrue        string = "FieldCheckboxValueIfTrue"
	FieldCheckboxValueIfFalse       string = "FieldCheckboxValueIfFalse"
	FieldCheckboxValuesUseInView    string = "FieldCheckboxValuesUseInView"
	FieldCheckboxValuesUseInStorage string = "FieldCheckboxValuesUseInStorage"
	FieldInputPlaceholder           string = "FieldInputPlaceholder"
	FieldDatetimeFormat             string = "FieldDatetimeFormat"
	FieldSelectOptions              string = "FieldSelectOptions"
	FieldPlaceholder                string = "FieldPlaceholder"
	FieldDefaultValue               string = "FieldDefaultValue"

	GroupViewTableIn2D             string = "GroupViewTableIn2D"
	GroupQueryAddFullTextSearchBox string = "GroupQueryAddFullTextSearchBox"
	GroupExtractAsSingleField      string = "GroupExtractAsSingleField"
	GroupReadOrderOfFields         string = "GroupReadOrderOfFields"
	GroupFields                    string = "GroupFields"

	DatabaseFieldAddDataToFullTextSearchIndex string = "DatabaseFieldAddDataToFullTextSearchIndex"
	DatabaseSkipDataExtraction                string = "DatabaseSkipDataExtraction"
	// DatabaseTableCollectionUid Unique id for a set of fields that belong to the same DatabaseTableCollectionName especially in a join query.
	//
	// Useful especially if metadatamodel/query contains more than one instance of DatabaseTableCollectionName group.
	DatabaseTableCollectionUid string = "DatabaseTableCollectionUid"
	// DatabaseTableCollectionName Name of table/collection that field/group belongs to.
	DatabaseTableCollectionName string = "DatabaseTableCollectionName"
	// DatabaseFieldColumnName Name of field in database.
	//
	// DatabaseFieldColumnName Ensure it is unique in a set of fields for a particular DatabaseTableCollectionName
	DatabaseFieldColumnName string = "DatabaseFieldColumnName"
	// DatabaseJoinDepth represents the join level.
	//
	// `0` represents no join while any number below `0` represents infinite join while any number above `0` represents max join.
	//
	// When combined with DatabaseTableCollectionName, it can form an alternative for a unique DatabaseTableCollectionUid.
	DatabaseJoinDepth string = "DatabaseJoinDepth"
	// DatabaseDistinct return unique DatabaseFieldColumnName results if `true`.
	DatabaseDistinct string = "DatabaseDistinct"
	// DatabaseSortByAsc Sort DatabaseFieldColumnName in ascending order if `true`.
	DatabaseSortByAsc string = "DatabaseSortByAsc"
	// DatabaseLimit Maximum number of results to return from a database query
	DatabaseLimit string = "DatabaseLimit"
	// DatabaseOffset Number of results to skip in database query.
	DatabaseOffset string = "DatabaseOffset"

	DatumInputView string = "DatumInputView"
)

/*
Properties for FieldUiSelect, Query condition value, etc
*/
const (
	Label string = "Label"
	Type  string = "Type"
	Value string = "Value"
)

/*
Filter Condition additional props
*/
const (
	FilterConditionNegate string = "FilterConditionNegate"
	FilterCondition       string = "FilterCondition"
	FieldGroup            string = "FieldGroup"
)

/*
Field of type FieldTypeAny properties
*/
const (
	FieldAnyMetadataModelActionID              string = "MetadataModelActionID"
	FieldAnyPickMetadataModelMessagePrompt     string = "PickMetadataModelMessagePrompt"
	FieldAnyGetMetadataModelPathToDataArgument string = "GetMetadataModelPathToDataArgument"
)

/*
Field type properties
*/
const (
	FieldTypeText      string = "Text"
	FieldTypeNumber    string = "Number"
	FieldTypeBoolean   string = "Boolean"
	FieldTypeTimestamp string = "Timestamp"
	FieldTypeAny       string = "Any"
)

/*
Field UI properties
*/
const (
	FieldUiText     string = "Text"
	FieldUiTextArea string = "TextArea"
	FieldUiNumber   string = "Number"
	FieldUiCheckbox string = "Checkbox"
	FieldUiSelect   string = "Select"
	FieldUiDatetime string = "DateTime"
)

/*
Date Time Formats for field type FieldTypeTimestamp
*/
const (
	FieldDatetimeFormatYYYYMMDDHHMM string = "yyyy-mm-dd hh:mm"
	FieldDatetimeFormatYYYYMMDD     string = "yyyy-mm-dd"
	FieldDatetimeFormatYYYYMM       string = "yyyy-mm"
	FieldDatetimeFormatYYYY         string = "yyyy"
	FieldDatetimeFormatMM           string = "mm"
	FieldDatetimeFormatHHMM         string = "hh:mm"
)
