/*
Package filter is for filtering through source data whose structure is defined by a metadata model.

Designed to support both simple queries and deeply nested logical operator queries which are extensible and customizable.

Below is an example query condition that is used to filter through data:

	{
	  "Type": "LogicalOperator",
	  "Negate": false,
	  "LogicalOperator": "And",
	  "Value": [
		{
		  "Type": "LogicalOperator",
		  "LogicalOperator": "Or",
		  "Value": [
			{
			  "Type": "FieldGroup",
			  "Negate": false,
			  "LogicalOperator": "And",
			  "Value": {
				"$.GroupFields[*].Bio": {
				  "EqualTo": {
					"AssumedFieldType": "Any",
					"Values": [
					  true,
					  "Yes"
					]
				  }
				}
			  }
			},
			{
			  "Type": "FieldGroup",
			  "Negate": false,
			  "LogicalOperator": "And",
			  "Value": {
				"$.GroupFields[*].Bio": {
				  "EqualTo": {
					"AssumedFieldType": "Text",
					"Negate": true,
					"Value": "no"
				  }
				},
				"$.GroupFields[*].Occ": {
				  "EqualTo": {
					"AssumedFieldType": "Text",
					"Negate": true,
					"Value": "no"
				  }
				}
			  }
			}
		  ]
		},
		{
		  "Type": "FieldGroup",
		  "Value": {
			"$.GroupFields[*].SiteAndGeoreferencing.GroupFields[*].Country": {
			  "FullTextSearchQuery": {
				"AssumedFieldType": "Text",
				"Value": "Kenya",
				"ExactMatch": true
			  }
			}
		  }
		},
		{
		  "Type": "FieldGroup",
		  "Negate": false,
		  "LogicalOperator": "And",
		  "Value": {
			"$.GroupFields[*].SiteAndGeoreferencing.GroupFields[*].Sites.GroupFields[*].Coordinates.GroupFields[*].Latitude": {
			  "GreaterThan": {
				"AssumedFieldType": "Number",
				"Value": 20.00
			  },
			  "LessThan": {
				"AssumedFieldType": "Number",
				"Value": 21.00
			  }
			}
		  }
		}
	  ]
	}

Example filtering data usage:

	import (
		gojsoncore "github.com/rogonion/go-json/core"
		"github.com/rogonion/go-json/object"
		"github.com/rogonion/go-metadatamodel/filter"
	)

	// Set metadata model
	var metadataModel gojsoncore.JsonObject

	// Set source data
	sourceData := object.NewObject()

	// Set query condition
	var queryCondition gojsoncore.JsonObject

	// Set other properties using builder pattern 'With' or 'Set'. Refer to filter.DataFilter structure.
	filterData := filter.NewFilterData(sourceData, metadataModel)

	var filterExcludeIndexes []int
	var err error

	filterExcludeIndexes, err = filterData.Filter(queryCondition, "", "")
*/
package filter
