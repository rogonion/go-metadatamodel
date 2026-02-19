package testdata

import (
	"github.com/brunoga/deep"
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
)

// UserMetadataModel returns a metadata model for the User entity.
// It defines fields: ID, Name, Email.
func UserMetadataModel(rootProperties gojsoncore.JsonObject) gojsoncore.JsonObject {
	const DefaultName = "User"

	if rootProperties == nil {
		rootProperties = make(gojsoncore.JsonObject)
	}

	if _, ok := rootProperties[core.FieldGroupJsonPathKey].(string); !ok {
		rootProperties[core.FieldGroupJsonPathKey] = path.JsonpathKeyRoot
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionName].(string); !ok {
		rootProperties[core.DatabaseTableCollectionName] = DefaultName
	}

	if _, ok := rootProperties[core.DatabaseJoinDepth].(float64); !ok {
		rootProperties[core.DatabaseJoinDepth] = float64(0)
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionUid].(string); !ok {
		rootProperties[core.DatabaseTableCollectionUid] = DefaultName
	}

	return deep.MustCopy(gojsoncore.JsonObject{
		core.FieldGroupJsonPathKey:       rootProperties[core.FieldGroupJsonPathKey],
		core.FieldGroupName:              DefaultName,
		core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
		core.DatabaseJoinDepth:           rootProperties[core.DatabaseJoinDepth],
		core.DatabaseTableCollectionName: rootProperties[core.DatabaseTableCollectionName],
		core.GroupFields: func() gojsoncore.JsonArray {
			FieldGroupJSONPathPrefixDepth0 := rootProperties[core.FieldGroupJsonPathKey].(string) + core.GroupJsonPathPrefix
			DatabaseJoinDepth := rootProperties[core.DatabaseJoinDepth].(float64)
			DatabaseTableCollectionName := rootProperties[core.DatabaseTableCollectionName].(string)
			return gojsoncore.JsonArray{
				gojsoncore.JsonObject{
					"ID": func() gojsoncore.JsonObject {
						const DefaultName = "ID"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeNumber,
							core.FieldUI:                     core.FieldUiNumber,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
							core.FieldGroupIsPrimaryKey:      true,
						}
					}(),
					"Name": func() gojsoncore.JsonObject {
						const DefaultName = "Name"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeText,
							core.FieldUI:                     core.FieldTypeText,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
						}
					}(),
					"Email": func() gojsoncore.JsonObject {
						const DefaultName = "Email"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeText,
							core.FieldUI:                     core.FieldTypeText,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
						}
					}(),
				},
			}
		}(),
		core.GroupReadOrderOfFields: gojsoncore.JsonArray{"ID", "Name", "Email"},
	})
}

// ProductMetadataModel returns a metadata model for the Product entity.
// It defines fields: ID, Name, Price.
func ProductMetadataModel(rootProperties gojsoncore.JsonObject) gojsoncore.JsonObject {
	const DefaultName = "Product"

	if rootProperties == nil {
		rootProperties = make(gojsoncore.JsonObject)
	}

	if _, ok := rootProperties[core.FieldGroupJsonPathKey].(string); !ok {
		rootProperties[core.FieldGroupJsonPathKey] = path.JsonpathKeyRoot
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionName].(string); !ok {
		rootProperties[core.DatabaseTableCollectionName] = DefaultName
	}

	if _, ok := rootProperties[core.DatabaseJoinDepth].(float64); !ok {
		rootProperties[core.DatabaseJoinDepth] = float64(0)
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionUid].(string); !ok {
		rootProperties[core.DatabaseTableCollectionUid] = DefaultName
	}

	return deep.MustCopy(gojsoncore.JsonObject{
		core.FieldGroupJsonPathKey:       rootProperties[core.FieldGroupJsonPathKey],
		core.FieldGroupName:              DefaultName,
		core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
		core.DatabaseJoinDepth:           rootProperties[core.DatabaseJoinDepth],
		core.DatabaseTableCollectionName: rootProperties[core.DatabaseTableCollectionName],
		core.GroupFields: func() gojsoncore.JsonArray {
			FieldGroupJSONPathPrefixDepth0 := rootProperties[core.FieldGroupJsonPathKey].(string) + core.GroupJsonPathPrefix
			DatabaseJoinDepth := rootProperties[core.DatabaseJoinDepth].(float64)
			DatabaseTableCollectionName := rootProperties[core.DatabaseTableCollectionName].(string)
			return gojsoncore.JsonArray{
				gojsoncore.JsonObject{
					"ID": func() gojsoncore.JsonObject {
						const DefaultName = "ID"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeNumber,
							core.FieldUI:                     core.FieldUiNumber,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
							core.FieldGroupIsPrimaryKey:      true,
						}
					}(),
					"Name": func() gojsoncore.JsonObject {
						const DefaultName = "Name"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeText,
							core.FieldUI:                     core.FieldTypeText,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
						}
					}(),
					"Price": func() gojsoncore.JsonObject {
						const DefaultName = "Price"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeNumber,
							core.FieldUI:                     core.FieldUiNumber,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
						}
					}(),
				},
			}
		}(),
		core.GroupReadOrderOfFields: gojsoncore.JsonArray{"ID", "Name", "Price"},
	})
}

// CompanyMetadataModel returns a metadata model for the Company entity.
// It defines fields: Name, Employees (nested UserMetadataModel).
func CompanyMetadataModel(rootProperties gojsoncore.JsonObject) gojsoncore.JsonObject {
	const DefaultName = "Company"

	if rootProperties == nil {
		rootProperties = make(gojsoncore.JsonObject)
	}

	if _, ok := rootProperties[core.FieldGroupJsonPathKey].(string); !ok {
		rootProperties[core.FieldGroupJsonPathKey] = path.JsonpathKeyRoot
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionName].(string); !ok {
		rootProperties[core.DatabaseTableCollectionName] = DefaultName
	}

	if _, ok := rootProperties[core.DatabaseJoinDepth].(float64); !ok {
		rootProperties[core.DatabaseJoinDepth] = float64(0)
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionUid].(string); !ok {
		rootProperties[core.DatabaseTableCollectionUid] = DefaultName
	}

	return deep.MustCopy(gojsoncore.JsonObject{
		core.FieldGroupJsonPathKey:       rootProperties[core.FieldGroupJsonPathKey],
		core.FieldGroupName:              DefaultName,
		core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
		core.DatabaseJoinDepth:           rootProperties[core.DatabaseJoinDepth],
		core.DatabaseTableCollectionName: rootProperties[core.DatabaseTableCollectionName],
		core.GroupFields: func() gojsoncore.JsonArray {
			FieldGroupJSONPathPrefixDepth0 := rootProperties[core.FieldGroupJsonPathKey].(string) + core.GroupJsonPathPrefix
			DatabaseJoinDepth := rootProperties[core.DatabaseJoinDepth].(float64)
			DatabaseTableCollectionName := rootProperties[core.DatabaseTableCollectionName].(string)
			return gojsoncore.JsonArray{
				gojsoncore.JsonObject{
					"Name": func() gojsoncore.JsonObject {
						const DefaultName = "Name"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeText,
							core.FieldUI:                     core.FieldTypeText,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
							core.FieldGroupIsPrimaryKey:      true,
						}
					}(),
					"Employees": func() gojsoncore.JsonObject {
						const DefaultName = "Employees"
						return UserMetadataModel(gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.DatabaseJoinDepth:           DatabaseJoinDepth + 1,
							core.DatabaseTableCollectionUid:  DefaultName,
							core.DatabaseTableCollectionName: DefaultName,
						})
					}(),
				},
			}
		}(),
		core.GroupReadOrderOfFields: gojsoncore.JsonArray{"Name", "Employees"},
	})
}

// AddressMetadataModel returns a metadata model for the Address entity.
// It defines fields: Street, City, ZipCode.
func AddressMetadataModel(rootProperties gojsoncore.JsonObject) gojsoncore.JsonObject {
	const DefaultName = "Address"

	if rootProperties == nil {
		rootProperties = make(gojsoncore.JsonObject)
	}

	if _, ok := rootProperties[core.FieldGroupJsonPathKey].(string); !ok {
		rootProperties[core.FieldGroupJsonPathKey] = path.JsonpathKeyRoot
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionName].(string); !ok {
		rootProperties[core.DatabaseTableCollectionName] = DefaultName
	}

	if _, ok := rootProperties[core.DatabaseJoinDepth].(float64); !ok {
		rootProperties[core.DatabaseJoinDepth] = float64(0)
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionUid].(string); !ok {
		rootProperties[core.DatabaseTableCollectionUid] = DefaultName
	}

	return deep.MustCopy(gojsoncore.JsonObject{
		core.FieldGroupJsonPathKey:       rootProperties[core.FieldGroupJsonPathKey],
		core.FieldGroupName:              DefaultName,
		core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
		core.DatabaseJoinDepth:           rootProperties[core.DatabaseJoinDepth],
		core.DatabaseTableCollectionName: rootProperties[core.DatabaseTableCollectionName],
		core.GroupFields: func() gojsoncore.JsonArray {
			FieldGroupJSONPathPrefixDepth0 := rootProperties[core.FieldGroupJsonPathKey].(string) + core.GroupJsonPathPrefix
			DatabaseJoinDepth := rootProperties[core.DatabaseJoinDepth].(float64)
			DatabaseTableCollectionName := rootProperties[core.DatabaseTableCollectionName].(string)
			return gojsoncore.JsonArray{
				gojsoncore.JsonObject{
					"Street": func() gojsoncore.JsonObject {
						const DefaultName = "Street"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeText,
							core.FieldUI:                     core.FieldTypeText,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
						}
					}(),
					"City": func() gojsoncore.JsonObject {
						const DefaultName = "City"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeText,
							core.FieldUI:                     core.FieldTypeText,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
						}
					}(),
					"ZipCode": func() gojsoncore.JsonObject {
						const DefaultName = "ZipCode"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeText,
							core.FieldUI:                     core.FieldTypeText,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
						}
					}(),
				},
			}
		}(),
		core.GroupReadOrderOfFields: gojsoncore.JsonArray{"Street", "City", "ZipCode"},
	})
}

// UserProfileMetadataModel returns a metadata model for the UserProfile entity.
// It defines fields: Name, Age, Address (nested AddressMetadataModel).
func UserProfileMetadataModel(rootProperties gojsoncore.JsonObject) gojsoncore.JsonObject {
	const DefaultName = "UserProfile"

	if rootProperties == nil {
		rootProperties = make(gojsoncore.JsonObject)
	}

	if _, ok := rootProperties[core.FieldGroupJsonPathKey].(string); !ok {
		rootProperties[core.FieldGroupJsonPathKey] = path.JsonpathKeyRoot
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionName].(string); !ok {
		rootProperties[core.DatabaseTableCollectionName] = DefaultName
	}

	if _, ok := rootProperties[core.DatabaseJoinDepth].(float64); !ok {
		rootProperties[core.DatabaseJoinDepth] = float64(0)
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionUid].(string); !ok {
		rootProperties[core.DatabaseTableCollectionUid] = DefaultName
	}

	return deep.MustCopy(gojsoncore.JsonObject{
		core.FieldGroupJsonPathKey:       rootProperties[core.FieldGroupJsonPathKey],
		core.FieldGroupName:              DefaultName,
		core.DatabaseJoinDepth:           rootProperties[core.DatabaseJoinDepth],
		core.DatabaseTableCollectionName: rootProperties[core.DatabaseTableCollectionName],
		core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
		core.GroupFields: func() gojsoncore.JsonArray {
			FieldGroupJSONPathPrefixDepth0 := rootProperties[core.FieldGroupJsonPathKey].(string) + core.GroupJsonPathPrefix
			DatabaseJoinDepth := rootProperties[core.DatabaseJoinDepth].(float64)
			DatabaseTableCollectionName := rootProperties[core.DatabaseTableCollectionName].(string)
			return gojsoncore.JsonArray{
				gojsoncore.JsonObject{
					"Name": func() gojsoncore.JsonObject {
						const DefaultName = "Name"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeText,
							core.FieldUI:                     core.FieldTypeText,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
							core.FieldGroupIsPrimaryKey:      true,
						}
					}(),
					"Age": func() gojsoncore.JsonObject {
						const DefaultName = "Age"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeNumber,
							core.FieldUI:                     core.FieldUiNumber,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
						}
					}(),
					"Address": func() gojsoncore.JsonObject {
						const AddressDefaultName = "Address"
						return AddressMetadataModel(gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + AddressDefaultName,
							core.FieldGroupName:              AddressDefaultName,
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionUid:  DefaultName,
							core.DatabaseTableCollectionName: DefaultName,
						})
					}(),
				},
			}
		}(),
		core.GroupReadOrderOfFields: gojsoncore.JsonArray{"Name", "Age", "Address"},
	})
}

// EmployeeMetadataModel returns a metadata model for the Employee entity.
// It defines fields: ID, Profile (nested UserProfileMetadataModel), Skills.
func EmployeeMetadataModel(rootProperties gojsoncore.JsonObject) gojsoncore.JsonObject {
	const DefaultName = "Employee"

	if rootProperties == nil {
		rootProperties = make(gojsoncore.JsonObject)
	}

	if _, ok := rootProperties[core.FieldGroupJsonPathKey].(string); !ok {
		rootProperties[core.FieldGroupJsonPathKey] = path.JsonpathKeyRoot
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionName].(string); !ok {
		rootProperties[core.DatabaseTableCollectionName] = DefaultName
	}

	if _, ok := rootProperties[core.DatabaseJoinDepth].(float64); !ok {
		rootProperties[core.DatabaseJoinDepth] = float64(0)
	}

	if _, ok := rootProperties[core.DatabaseTableCollectionUid].(string); !ok {
		rootProperties[core.DatabaseTableCollectionUid] = DefaultName
	}

	return deep.MustCopy(gojsoncore.JsonObject{
		core.FieldGroupJsonPathKey:       rootProperties[core.FieldGroupJsonPathKey],
		core.FieldGroupName:              DefaultName,
		core.DatabaseJoinDepth:           rootProperties[core.DatabaseJoinDepth],
		core.DatabaseTableCollectionName: rootProperties[core.DatabaseTableCollectionName],
		core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
		core.GroupFields: func() gojsoncore.JsonArray {
			FieldGroupJSONPathPrefixDepth0 := rootProperties[core.FieldGroupJsonPathKey].(string) + core.GroupJsonPathPrefix
			DatabaseJoinDepth := rootProperties[core.DatabaseJoinDepth].(float64)
			DatabaseTableCollectionName := rootProperties[core.DatabaseTableCollectionName].(string)
			return gojsoncore.JsonArray{
				gojsoncore.JsonObject{
					"ID": func() gojsoncore.JsonObject {
						const DefaultName = "ID"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeNumber,
							core.FieldUI:                     core.FieldUiNumber,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
							core.FieldGroupIsPrimaryKey:      true,
						}
					}(),
					"Profile": func() gojsoncore.JsonObject {
						const DefaultName = "Profile"
						return UserProfileMetadataModel(gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.DatabaseJoinDepth:           DatabaseJoinDepth + 1,
							core.DatabaseTableCollectionUid:  DefaultName,
							core.DatabaseTableCollectionName: DefaultName,
						})
					}(),
					"Skills": func() gojsoncore.JsonObject {
						const DefaultName = "Skills"
						return gojsoncore.JsonObject{
							core.FieldGroupJsonPathKey:       FieldGroupJSONPathPrefixDepth0 + DefaultName,
							core.FieldGroupName:              DefaultName,
							core.FieldDataType:               core.FieldTypeText,
							core.FieldUI:                     core.FieldTypeText,
							core.DatabaseTableCollectionUid:  rootProperties[core.DatabaseTableCollectionUid],
							core.DatabaseJoinDepth:           DatabaseJoinDepth,
							core.DatabaseTableCollectionName: DatabaseTableCollectionName,
							core.DatabaseFieldColumnName:     DefaultName,
							core.FieldGroupMaxEntries:        0,
						}
					}(),
				},
			}
		}(),
		core.GroupReadOrderOfFields: gojsoncore.JsonArray{"ID", "Profile", "Skills"},
	})
}
