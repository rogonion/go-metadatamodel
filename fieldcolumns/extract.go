package fieldcolumns

import (
	gojsoncore "github.com/rogonion/go-json/core"
	"github.com/rogonion/go-json/path"
	"github.com/rogonion/go-json/schema"
	"github.com/rogonion/go-metadatamodel/core"
)

func (n *ExtractFields) Extract() error {
	n.fieldColumns = new(FieldColumns)
}

func (n *ExtractFields) recursiveExtract(group any)

func (n *ExtractFields) GetFieldColumns() *FieldColumns {
	return n.fieldColumns
}

func (n *ExtractFields) WithSchema(value schema.Schema) *ExtractFields {
	n.SetSchema(value)
	return n
}

func (n *ExtractFields) SetSchema(value schema.Schema) {
	n.schema = value
}

func (n *ExtractFields) WithAdd(value core.FieldGroupPropertiesMatch) *ExtractFields {
	n.SetAdd(value)
	return n
}

func (n *ExtractFields) SetAdd(value core.FieldGroupPropertiesMatch) {
	n.add = value
}

func (n *ExtractFields) WithSkip(value core.FieldGroupPropertiesMatch) *ExtractFields {
	n.SetSkip(value)
	return n
}

func (n *ExtractFields) SetSkip(value core.FieldGroupPropertiesMatch) {
	n.skip = value
}

func NewExtractFields(metadataModel gojsoncore.JsonObject) *ExtractFields {
	n := new(ExtractFields)
	n.metadataModel = metadataModel
	return n
}

type ExtractFields struct {
	metadataModel gojsoncore.JsonObject

	schema schema.Schema

	fieldColumns *FieldColumns

	// skip a field column if its properties matches one of the entries values.
	skip core.FieldGroupPropertiesMatch

	// add a field column if its properties matches one of the entries values.
	add core.FieldGroupPropertiesMatch
}

type FieldColumns struct {
	Fields                              map[string]*FieldColumn
	ReadOrderOfFieldColumnsJsonPathKeys []*ReadOrderOfFieldColumnsJsonPathKeys
}

type ReadOrderOfFieldColumnsJsonPathKeys struct {
	FieldGroupJsonPathKey         path.JSONPath
	ViewInSeparateColumns         bool
	IndexOfValueInSeparateColumns int
}

type FieldColumn struct {
	Property gojsoncore.JsonObject
	Schema   schema.Schema
}
