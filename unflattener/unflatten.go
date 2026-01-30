package unflattener

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/object"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-metadatamodel/core"
	"github.com/rogonion/go-metadatamodel/fieldcolumns"
	"github.com/rogonion/go-metadatamodel/flattener"
)

func (n *Unflattener) recursiveConvert(groupIndexTree *RecursiveGroupIndexTree, parentCollection *GroupCollection, linearCollectionIndexes []int) error {
	const FunctionName = "recursiveConvert"

	// 1. Identify Instance (Signature)
	var pkColumns []int
	if groupIndexTree.GroupColumnIndexes != nil {
		if len(groupIndexTree.GroupColumnIndexes.Primary) > 0 {
			pkColumns = groupIndexTree.GroupColumnIndexes.Primary
		} else {
			pkColumns = groupIndexTree.GroupColumnIndexes.All
		}
	}

	signature := n.signature.GenerateSignature(n.currentSourceRow, pkColumns)

	// 2. Get/Create Instance
	node, instanceIndex := parentCollection.GetOrCreateInstance(signature, groupIndexTree.FieldColumnPosition.FieldGroupJsonPathKey)

	currentPathIndexes := append(linearCollectionIndexes, instanceIndex)

	// 3. Write Fields
	if groupIndexTree.GroupColumnIndexes != nil {
		for _, colIndex := range groupIndexTree.GroupColumnIndexes.All {
			if colIndex >= len(n.currentSourceRow) {
				continue
			}
			val := n.currentSourceRow[colIndex]

			if !val.IsValid() {
				continue
			}

			// Check if the underlying value is nil (for nullable types)
			switch val.Kind() {
			case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func, reflect.UnsafePointer:
				if val.IsNil() {
					continue
				}
			}

			colField, ok := n.columnFields.GetColumnFieldByIndexInUnskippedReadOrder(colIndex)
			if !ok {
				continue
			}

			// Construct Target Path
			targetPath, err := core.NewJsonPathToValue().WithSourceOfValueIsAnArray(true).Get(colField.FieldColumnPosition.JSONPath(), currentPathIndexes)
			if err != nil {
				return NewError().WithFunctionName(FunctionName).WithMessage("resolve field path failed").WithNestedError(err)
			}
			targetPathSuffixIsLinearCollection := strings.HasSuffix(string(targetPath), path.JsonpathRightBracket)

			if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
				if !targetPathSuffixIsLinearCollection {
					newVal := reflect.MakeSlice(reflect.SliceOf(val.Type()), 1, 1)
					newVal.Index(0).Set(val)
					val = newVal
				}
			} else {
				// Target Path end is slice/array, unwrap
				if targetPathSuffixIsLinearCollection {
					val = val.Index(0)
				}
			}

			// Write
			if _, err := n.destination.SetReflect(targetPath, val); err != nil {
				return NewError().WithFunctionName(FunctionName).WithMessage("write field failed").WithNestedError(err)
			}
		}
	}

	// 4. Recurse
	for _, childTree := range groupIndexTree.GroupFields {
		// Find the child collection in the current node
		childCollection := node.GetOrCreateGroup(childTree.Suffix)

		if err := n.recursiveConvert(childTree, childCollection, currentPathIndexes); err != nil {
			return err
		}
	}

	return nil
}

func (n *Unflattener) Unflatten(source flattener.FlattenedTable) error {
	const FunctionName = "Unflatten"

	if n.columnFields == nil {
		if columnFields, err := fieldcolumns.NewColumnFieldsExtraction(n.metadataModel).Extract(); err != nil {
			return NewError().WithFunctionName(FunctionName).WithMessage("extract default columnFields failed").WithNestedError(err)
		} else {
			columnFields.Reposition()
			columnFields.Skip(nil, nil)
			n.columnFields = columnFields
		}
	}

	if n.recursiveIndexTree == nil {
		if value, err := n.recursiveInitGroupIndexTree(n.metadataModel, path.JSONPath(path.JsonpathKeyRoot)); err != nil {
			return err
		} else {
			n.recursiveIndexTree = value
		}
	}

	if n.index == nil {
		n.index = &GroupCollection{
			Instances: make(GroupCollectionInstances),
		}
	}

	for _, row := range source {
		n.currentSourceRow = row
		if err := n.recursiveConvert(n.recursiveIndexTree, n.index, []int{}); err != nil {
			return err
		}
	}

	return nil
}

func (n *Unflattener) recursiveInitGroupIndexTree(group any, groupJsonPathKey path.JSONPath) (*RecursiveGroupIndexTree, error) {
	const FunctionName = "recursiveInitGroupIndexTree"

	fieldGroupProp, err := core.AsJsonObject(group)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("group not JsonObject").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	groupFields, err := core.GetGroupFields(fieldGroupProp)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get group fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	groupReadOrderOfFields, err := core.GetGroupReadOrderOfFields(fieldGroupProp)
	if err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get group read order of fields failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	}

	groupIndexTree := &RecursiveGroupIndexTree{
		FieldColumnPosition: &fieldcolumns.FieldColumnPosition{
			FieldGroupJsonPathKey: groupJsonPathKey,
		},
		GroupFields: make([]*RecursiveGroupIndexTree, 0),
	}

	if value, err := fieldcolumns.NewGroupsColumnsIndexesRetrieval(n.columnFields).Get(fieldGroupProp); err != nil {
		return nil, NewError().WithFunctionName(FunctionName).WithMessage("get group column indexes failed").WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
	} else {
		groupIndexTree.GroupColumnIndexes = value
	}

	for _, fgKeySuffix := range groupReadOrderOfFields {
		fgProperty, err := core.AsJsonObject(groupFields[fgKeySuffix])
		if err != nil {
			return nil, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
		}

		fgJsonPathKey, err := core.AsJSONPath(fgProperty[core.FieldGroupJsonPathKey])
		if err != nil {
			return nil, NewError().WithFunctionName(FunctionName).WithMessage(fmt.Sprintf("get FieldGroupJsonPathKey for field with suffix key '%s' failed", fgKeySuffix)).WithNestedError(err).WithData(gojsoncore.JsonObject{"Group": group})
		}

		// Field is a Group
		if ok := core.IsFieldAGroup(fgProperty); ok {
			// Process group as a field
			if extractAsSingleField, ok := fgProperty[core.GroupExtractAsSingleField].(bool); ok && extractAsSingleField {
				continue
			}

			if !core.DoesFieldGroupFieldsContainNestedGroupFields(fgProperty) {
				// Process group as field with set of fields in separate columns
				if fgViewMaxNoOfValuesInSeparateColumns, _ := core.GetMaximumFlatNoOfColumns(fgProperty); fgViewMaxNoOfValuesInSeparateColumns > 0 {
					if _, err := core.GetGroupFields(fgProperty); err == nil {
						if _, err := core.GetGroupReadOrderOfFields(fgProperty); err == nil {
							continue
						}
					}
				}
			}

			if value, err := n.recursiveInitGroupIndexTree(fgProperty, fgJsonPathKey); err != nil {
				if !errors.Is(err, ErrNoGroupFields) {
					return nil, err
				}
			} else {
				value.Suffix = fgKeySuffix
				groupIndexTree.GroupFields = append(groupIndexTree.GroupFields, value)
			}
			continue
		}
	}

	return groupIndexTree, nil
}

func (n *Unflattener) WithDestination(value *object.Object) *Unflattener {
	n.SetDestination(value)
	return n
}

func (n *Unflattener) SetDestination(value *object.Object) {
	n.destination = value
}

func (n *Unflattener) WithColumnFields(value *fieldcolumns.ColumnFields) *Unflattener {
	n.SetColumnFields(value)
	return n
}

func (n *Unflattener) SetColumnFields(value *fieldcolumns.ColumnFields) {
	n.columnFields = value
}

func NewUnflattener(metadataModel gojsoncore.JsonObject, signature *Signature) *Unflattener {
	return &Unflattener{
		metadataModel: metadataModel,
		signature:     signature,
	}
}

/*
Unflattener converts a 2 dimension linear collection (like a 2D array) into a deeply nested mix of associative collections (like an array of objects).
*/
type Unflattener struct {
	metadataModel gojsoncore.JsonObject

	// columnFields extracted fields as table columns from metadataModel
	columnFields *fieldcolumns.ColumnFields

	// currentReadOrderOfFields retrieved from columnFields.
	//
	// The current read order of columns as it is in the currentSourceObject.
	currentReadOrderOfFields []int

	// currentSourceRow is the current row being processed.
	currentSourceRow flattener.FlattenedRow

	// destination where to write currentSourceObject to.
	destination *object.Object

	// Get signature to be used as index key based on primary key values.
	signature *Signature

	// Tree tracking the indexes at each level.
	//
	// Makes it possible to stream read currentSourceObject via Unflattener.Unflatten.
	index *GroupCollection

	// recursiveIndexTree data (tree of fields/groups) to use when converting the currentSourceObject.
	recursiveIndexTree *RecursiveGroupIndexTree
}

/*
RecursiveGroupIndexTree represents tree of field/groups to read for Unflattener.
*/
type RecursiveGroupIndexTree struct {
	FieldColumnPosition *fieldcolumns.FieldColumnPosition

	GroupColumnIndexes *fieldcolumns.GroupColumnIndexes

	GroupFields []*RecursiveGroupIndexTree

	Suffix string
}
