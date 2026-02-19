package testdata

import (
	"reflect"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/schema"
)

// UserSchema returns the schema definition for the User struct.
func UserSchema() *schema.DynamicSchemaNode {
	return &schema.DynamicSchemaNode{
		Kind: reflect.Struct,
		Type: reflect.TypeOf(User{}),
		ChildNodes: schema.ChildNodes{
			"ID": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]int{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.Int,
					Type: reflect.TypeOf(int(0)),
				},
			},
			"Name": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]string{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.String,
					Type: reflect.TypeOf(""),
				},
			},
			"Email": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]string{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.String,
					Type: reflect.TypeOf(""),
				},
			},
		},
	}
}

// ProductSchema returns the schema definition for the Product struct.
func ProductSchema() *schema.DynamicSchemaNode {
	return &schema.DynamicSchemaNode{
		Kind: reflect.Struct,
		Type: reflect.TypeOf(Product{}),
		ChildNodes: schema.ChildNodes{
			"ID": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]int{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.Int,
					Type: reflect.TypeOf(int(0)),
				},
			},
			"Name": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]string{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.String,
					Type: reflect.TypeOf(""),
				},
			},
			"Price": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]float64{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.Float64,
					Type: reflect.TypeOf(float64(0)),
				},
			},
		},
	}
}

// CompanySchema returns the schema definition for the Company struct.
func CompanySchema() *schema.DynamicSchemaNode {
	return &schema.DynamicSchemaNode{
		Kind: reflect.Struct,
		Type: reflect.TypeOf(Company{}),
		ChildNodes: schema.ChildNodes{
			"Name": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]string{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.String,
					Type: reflect.TypeOf(""),
				},
			},
			"Employees": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]*User{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind:                    reflect.Pointer,
					Type:                    reflect.TypeOf(gojsoncore.Ptr(User{})),
					ChildNodesPointerSchema: UserSchema(),
				},
			},
		},
	}
}

// AddressSchema returns the schema definition for the Address struct.
func AddressSchema() *schema.DynamicSchemaNode {
	return &schema.DynamicSchemaNode{
		Kind: reflect.Struct,
		Type: reflect.TypeOf(Address{}),
		ChildNodes: schema.ChildNodes{
			"Street": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]string{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.String,
					Type: reflect.TypeOf(""),
				},
			},
			"City": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]string{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.String,
					Type: reflect.TypeOf(""),
				},
			},
			"ZipCode": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]*string{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.Pointer,
					Type: reflect.TypeOf(gojsoncore.Ptr("")),
					ChildNodesPointerSchema: &schema.DynamicSchemaNode{
						Kind: reflect.String,
						Type: reflect.TypeOf(""),
					},
				},
			},
		},
	}
}

// UserProfileSchema returns the schema definition for the UserProfile struct.
func UserProfileSchema() *schema.DynamicSchemaNode {
	return &schema.DynamicSchemaNode{
		Kind: reflect.Struct,
		Type: reflect.TypeOf(UserProfile{}),
		ChildNodes: schema.ChildNodes{
			"Name": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]string{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.String,
					Type: reflect.TypeOf(""),
				},
			},
			"Age": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]int{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.Int,
					Type: reflect.TypeOf(0),
				},
			},
			"Address": &schema.DynamicSchemaNode{
				Kind:                                     reflect.Slice,
				Type:                                     reflect.TypeOf([]Address{}),
				ChildNodesLinearCollectionElementsSchema: AddressSchema(),
			},
		},
	}
}

// EmployeeSchema returns the schema definition for the Employee struct.
func EmployeeSchema() *schema.DynamicSchemaNode {
	return &schema.DynamicSchemaNode{
		Kind: reflect.Struct,
		Type: reflect.TypeOf(Employee{}),
		ChildNodes: schema.ChildNodes{
			"ID": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]int{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.Int,
					Type: reflect.TypeOf(0),
				},
			},
			"Profile": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]*UserProfile{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind:                    reflect.Pointer,
					Type:                    reflect.TypeOf(gojsoncore.Ptr(UserProfile{})),
					ChildNodesPointerSchema: UserProfileSchema(),
				},
			},
			"Skills": &schema.DynamicSchemaNode{
				Kind: reflect.Slice,
				Type: reflect.TypeOf([]string{}),
				ChildNodesLinearCollectionElementsSchema: &schema.DynamicSchemaNode{
					Kind: reflect.String,
					Type: reflect.TypeOf(""),
				},
			},
		},
	}
}
