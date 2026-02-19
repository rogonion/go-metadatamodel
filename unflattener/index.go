package unflattener

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/flattener"
)

// GenerateSignature creates a deterministic unique key for a set of columns (the PKs).
// It uses the Signature.joinSymbol separator to prevent concatenation collisions.
func (n *Signature) GenerateSignature(row flattener.FlattenedRow, readOrderOfRow []int) string {
	// Optimization: Singleton groups (no PK) always return empty string
	if len(readOrderOfRow) == 0 {
		return ""
	}

	var b strings.Builder
	// Heuristic: Pre-allocate 32 bytes per key column to minimize resize allocations
	b.Grow(len(readOrderOfRow) * 32)

	for i, colIdx := range readOrderOfRow {
		if i > 0 {
			b.WriteByte(n.joinSymbol) // Separator
		}

		// 1. Safety Checks
		if colIdx >= len(row) {
			b.WriteString("nil")
			continue
		}

		val := row[colIdx]
		if !val.IsValid() {
			b.WriteString("nil")
			continue
		}

		// 2. Fast Path (Primitives)
		// We handle common types directly to avoid the overhead of the Converter module.
		switch val.Kind() {
		case reflect.String:
			b.WriteString(val.String())
			continue
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b.WriteString(strconv.FormatInt(val.Int(), 10))
			continue
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			b.WriteString(strconv.FormatUint(val.Uint(), 10))
			continue
		case reflect.Bool:
			if val.Bool() {
				b.WriteString("t")
			} else {
				b.WriteString("f")
			}
			continue
		default:
			if convertedVal, err := n.converter.ConvertNode(val, n.signatureSchema); err == nil {
				b.WriteString(convertedVal.String())
				continue
			}
		}

		// 4. Fallback
		// If all else fails, rely on Go's default formatting (e.g. for floats or unregistered structs)
		b.WriteString(fmt.Sprint(val.Interface()))
	}

	return b.String()
}

// WithConverter sets the schema converter used for generating signatures from non-primitive types.
func (n *Signature) WithConverter(value *schema.Conversion) *Signature {
	n.SetConverter(value)
	return n
}

// SetConverter sets the schema converter.
func (n *Signature) SetConverter(value *schema.Conversion) {
	n.converter = value
}

// WithJoinSymbol sets the separator symbol used in signatures.
func (n *Signature) WithJoinSymbol(value byte) *Signature {
	n.SetJoinSymbol(value)
	return n
}

// SetJoinSymbol sets the separator symbol.
func (n *Signature) SetJoinSymbol(value byte) {
	n.joinSymbol = value
}

// NewSignature creates a new Signature instance with default settings.
func NewSignature() *Signature {
	n := new(Signature)
	n.SetJoinSymbol('|')
	n.converter = schema.NewConversion()
	n.signatureSchema = &schema.DynamicSchemaNode{
		Kind: reflect.String,
		Type: reflect.TypeOf(""),
	}
	return n
}

// Signature handles the generation of unique keys (signatures) for rows based on primary key columns.
type Signature struct {
	converter       *schema.Conversion
	signatureSchema *schema.DynamicSchemaNode
	joinSymbol      byte
}

// GetOrCreateGroup retrieves or creates a child GroupCollection for a specific nested group suffix.
func (n *GroupIndexNode) GetOrCreateGroup(suffix string) *GroupCollection {
	if groupCollection, exists := n.Groups[suffix]; exists {
		return groupCollection
	}

	newGroupCollection := &GroupCollection{
		NextIndex: 0,
		Instances: make(GroupCollectionInstances),
	}
	n.Groups[suffix] = newGroupCollection

	return newGroupCollection
}

// GroupIndexNode represents a specific instance of an element (e.g., "Employee #1").
type GroupIndexNode struct {
	JsonPathKey path.JSONPath

	// The resolved index of THIS node in its parent's list.
	MyIndex int

	// Nested groups belonging to this instance.
	// Key is the group suffix (e.g., "Address", "Profile").
	Groups GroupIndexNodeGroups
}

// GroupIndexNodeGroups is a map of nested group collections.
type GroupIndexNodeGroups map[string]*GroupCollection

// GetOrCreateInstance retrieves or creates a GroupIndexNode for a given signature (Primary Key).
func (n *GroupCollection) GetOrCreateInstance(signature string, jsonPath path.JSONPath) (*GroupIndexNode, int) {
	if instance, exists := n.Instances[signature]; exists {
		return instance, instance.MyIndex
	}

	// Create New
	newIdx := n.NextIndex
	newInstance := &GroupIndexNode{
		JsonPathKey: jsonPath,
		MyIndex:     newIdx,
		Groups:      make(GroupIndexNodeGroups),
	}

	n.Instances[signature] = newInstance
	n.NextIndex++ // Increment for next time

	return newInstance, newIdx
}

// GroupCollection represents a "List" of child nodes of a specific type.
// (e.g., The list of Addresses belonging to Employee #1)
//
// It holds the state required to append new items to this specific list.
type GroupCollection struct {
	// The Counter: Tracks the next available index for this list.
	NextIndex int

	// The Registry: Hash map of current elements.
	// Key is the unique signature (Primary Key).
	Instances GroupCollectionInstances
}

// GroupCollectionInstances is a map of group instances keyed by their signature.
type GroupCollectionInstances map[string]*GroupIndexNode
