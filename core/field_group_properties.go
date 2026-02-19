package core

/*
Metadata Model Field/Group Properties.
*/
const (
	// FieldGroupJsonPathKey is the key used to store the JSON Path to the value in the source object.
	FieldGroupJsonPathKey string = "FieldGroupJsonPathKey"
	// FieldGroupName is the key for the display name of the field or group.
	FieldGroupName string = "FieldGroupName"
	// FieldGroupDescription is the key for the description of the field or group.
	FieldGroupDescription string = "FieldGroupDescription"
	// FieldGroupViewTableLockColumn is the key to indicate if the column should be locked/frozen in a table view.
	FieldGroupViewTableLockColumn string = "FieldGroupViewTableLockColumn"
	// FieldGroupIsPrimaryKey is the key to indicate if the field is a primary key.
	FieldGroupIsPrimaryKey string = "FieldGroupIsPrimaryKey"

	// FieldGroupViewValuesInSeparateColumns For when you want to view the individual values in an array/slice of a field or non-nested group in separate columns in a flat/table.
	//
	// By default:
	//	- Fields with multiple values will be joined as a string separated by a comma or FieldMultipleValuesJoinSymbol.
	//  - Groups will appear as separate rows in a matrix.
	FieldGroupViewValuesInSeparateColumns string = "FieldGroupViewValuesInSeparateColumns"
	// FieldGroupViewMaxNoOfValuesInSeparateColumns is the key for the maximum number of columns to generate when pivoting array values.
	FieldGroupViewMaxNoOfValuesInSeparateColumns string = "FieldGroupViewMaxNoOfValuesInSeparateColumns"
	// FieldViewValuesInSeparateColumnsHeaderFormat is the key for the format string used to generate headers for pivoted columns (e.g., "Address %d").
	FieldViewValuesInSeparateColumnsHeaderFormat string = "FieldViewValuesInSeparateColumnsHeaderFormat"
	// FieldViewValuesInSeparateColumnsHeaderIndex is the key for the index of the specific pivoted column (0-based).
	FieldViewValuesInSeparateColumnsHeaderIndex string = "FieldViewValuesInSeparateColumnsHeaderIndex"
	// FieldMultipleValuesJoinSymbol is the key for the symbol used to join multiple values into a single string if not pivoted.
	FieldMultipleValuesJoinSymbol string = "FieldMultipleValuesJoinSymbol"

	// FieldGroupInputDisable is the key to disable input for this field.
	FieldGroupInputDisable string = "FieldGroupInputDisable"
	// FieldGroupDisablePropertiesEdit is the key to disable editing of the field's properties.
	FieldGroupDisablePropertiesEdit string = "FieldGroupDisablePropertiesEdit"
	// FieldGroupViewDisable is the key to hide the field from the view.
	FieldGroupViewDisable string = "FieldGroupViewDisable"
	// FieldGroupQueryConditionsEditDisable is the key to disable editing of query conditions for this field.
	FieldGroupQueryConditionsEditDisable string = "FieldGroupQueryConditionsEditDisable"
	// FieldGroupMaxEntries is the key for the maximum number of entries allowed.
	FieldGroupMaxEntries string = "FieldGroupMaxEntries"

	// FieldDataType is the key for the data type of the field (e.g., Text, Number).
	FieldDataType string = "FieldDataType"
	// FieldUI is the key for the UI component type (e.g., Text, Select).
	FieldUI string = "FieldUi"

	// FieldColumnPosition for setting a custom position for field/group when working with data in a flat2D
	FieldColumnPosition string = "FieldColumnPosition"
	// FieldGroupPositionBefore is the key to indicate if the field should be positioned before another specific field.
	FieldGroupPositionBefore string = "FieldGroupPositionBefore"

	// FieldGroupTypeAny is the key for properties specific to 'Any' type fields.
	FieldGroupTypeAny string = "FieldGroupTypeAny"
	// FieldCheckboxValueIfTrue is the key for the value to store when a checkbox is checked.
	FieldCheckboxValueIfTrue string = "FieldCheckboxValueIfTrue"
	// FieldCheckboxValueIfFalse is the key for the value to store when a checkbox is unchecked.
	FieldCheckboxValueIfFalse string = "FieldCheckboxValueIfFalse"
	// FieldCheckboxValuesUseInView is the key to indicate if specific values should be used in the view for checkboxes.
	FieldCheckboxValuesUseInView string = "FieldCheckboxValuesUseInView"
	// FieldCheckboxValuesUseInStorage is the key to indicate if specific values should be used in storage for checkboxes.
	FieldCheckboxValuesUseInStorage string = "FieldCheckboxValuesUseInStorage"
	// FieldInputPlaceholder is the key for the input placeholder text.
	FieldInputPlaceholder string = "FieldInputPlaceholder"
	// FieldDatetimeFormat is the key for the date/time format string.
	FieldDatetimeFormat string = "FieldDatetimeFormat"
	// FieldSelectOptions is the key for the list of options in a select/dropdown.
	FieldSelectOptions string = "FieldSelectOptions"
	// FieldPlaceholder is the key for a placeholder value.
	FieldPlaceholder string = "FieldPlaceholder"
	// FieldDefaultValue is the key for the default value of the field.
	FieldDefaultValue string = "FieldDefaultValue"

	// GroupViewTableIn2D is the key to indicate if the group should be viewed as a 2D table.
	GroupViewTableIn2D string = "GroupViewTableIn2D"
	// GroupQueryAddFullTextSearchBox is the key to add a full-text search box for the group.
	GroupQueryAddFullTextSearchBox string = "GroupQueryAddFullTextSearchBox"
	// GroupExtractAsSingleField is the key to treat a group as a single field during extraction (no recursion).
	GroupExtractAsSingleField string = "GroupExtractAsSingleField"
	// GroupReadOrderOfFields is the key for the array defining the order of fields in the group.
	GroupReadOrderOfFields string = "GroupReadOrderOfFields"
	// GroupFields is the key for the map containing the child fields of the group.
	GroupFields string = "GroupFields"

	// DatabaseFieldAddDataToFullTextSearchIndex is the key to indicate if the field should be added to the full-text search index.
	DatabaseFieldAddDataToFullTextSearchIndex string = "DatabaseFieldAddDataToFullTextSearchIndex"
	// DatabaseSkipDataExtraction is the key to skip this field during database extraction.
	DatabaseSkipDataExtraction string = "DatabaseSkipDataExtraction"
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

	// DatumInputView is the key for the datum input view.
	DatumInputView string = "DatumInputView"
)

/*
Properties for FieldUiSelect, Query condition value, etc
*/
const (
	// Label is the key for the label property.
	Label string = "Label"
	// Type is the key for the type property.
	Type string = "Type"
	// Value is the key for the value property.
	Value string = "Value"
)

/*
Filter Condition additional props
*/
const (
	// FilterConditionNegate is the key to negate a filter condition.
	FilterConditionNegate string = "FilterConditionNegate"
	// FilterCondition is the key for the filter condition.
	FilterCondition string = "FilterCondition"
	// FieldGroup is the key for the field group reference.
	FieldGroup string = "FieldGroup"
)

/*
Field of type FieldTypeAny properties
*/
const (
	// FieldAnyMetadataModelActionID is the key for the action ID in 'Any' type fields.
	FieldAnyMetadataModelActionID string = "MetadataModelActionID"
	// FieldAnyPickMetadataModelMessagePrompt is the key for the message prompt in 'Any' type fields.
	FieldAnyPickMetadataModelMessagePrompt string = "PickMetadataModelMessagePrompt"
	// FieldAnyGetMetadataModelPathToDataArgument is the key for the path argument in 'Any' type fields.
	FieldAnyGetMetadataModelPathToDataArgument string = "GetMetadataModelPathToDataArgument"
)

/*
Field type properties
*/
const (
	// FieldTypeText represents the Text data type.
	FieldTypeText string = "Text"
	// FieldTypeNumber represents the Number data type.
	FieldTypeNumber string = "Number"
	// FieldTypeBoolean represents the Boolean data type.
	FieldTypeBoolean string = "Boolean"
	// FieldTypeTimestamp represents the Timestamp data type.
	FieldTypeTimestamp string = "Timestamp"
	// FieldTypeAny represents the Any data type.
	FieldTypeAny string = "Any"
)

/*
Field UI properties
*/
const (
	// FieldUiText represents the Text UI component.
	FieldUiText string = "Text"
	// FieldUiTextArea represents the TextArea UI component.
	FieldUiTextArea string = "TextArea"
	// FieldUiNumber represents the Number UI component.
	FieldUiNumber string = "Number"
	// FieldUiCheckbox represents the Checkbox UI component.
	FieldUiCheckbox string = "Checkbox"
	// FieldUiSelect represents the Select UI component.
	FieldUiSelect string = "Select"
	// FieldUiDatetime represents the DateTime UI component.
	FieldUiDatetime string = "DateTime"
)

/*
Date Time Formats for field type FieldTypeTimestamp
*/
const (
	// FieldDatetimeFormatYYYYMMDDHHMM represents the format "yyyy-mm-dd hh:mm".
	FieldDatetimeFormatYYYYMMDDHHMM string = "yyyy-mm-dd hh:mm"
	// FieldDatetimeFormatYYYYMMDD represents the format "yyyy-mm-dd".
	FieldDatetimeFormatYYYYMMDD string = "yyyy-mm-dd"
	// FieldDatetimeFormatYYYYMM represents the format "yyyy-mm".
	FieldDatetimeFormatYYYYMM string = "yyyy-mm"
	// FieldDatetimeFormatYYYY represents the format "yyyy".
	FieldDatetimeFormatYYYY string = "yyyy"
	// FieldDatetimeFormatMM represents the format "mm".
	FieldDatetimeFormatMM string = "mm"
	// FieldDatetimeFormatHHMM represents the format "hh:mm".
	FieldDatetimeFormatHHMM string = "hh:mm"
)
