package types

import "fmt"

// DataType represents the supported data types
type DataType int

const (
	TypeInteger DataType = iota
	TypeText
	TypeBoolean
)

func (dt DataType) String() string {
	switch dt {
	case TypeInteger:
		return "INTEGER"
	case TypeText:
		return "TEXT"
	case TypeBoolean:
		return "BOOLEAN"
	default:
		return "UNKNOWN"
	}
}

// ParseDataType converts a string to DataType
func ParseDataType(s string) (DataType, error) {
	switch s {
	case "INTEGER", "INT":
		return TypeInteger, nil
	case "TEXT", "VARCHAR", "STRING":
		return TypeText, nil
	case "BOOLEAN", "BOOL":
		return TypeBoolean, nil
	default:
		return 0, fmt.Errorf("unknown data type: %s", s)
	}
}

// Column represents a column definition
type Column struct {
	Name     string
	DataType DataType
}

// Schema represents a table schema
type Schema struct {
	Columns []Column
}

// Row represents a single row of data
type Row struct {
	Values []Value
}

// Value represents a single value in a cell
type Value struct {
	Type    DataType
	IsNull  bool
	IntVal  int64
	TextVal string
	BoolVal bool
}

// NewIntValue creates a new integer value
func NewIntValue(v int64) Value {
	return Value{Type: TypeInteger, IntVal: v}
}

// NewTextValue creates a new text value
func NewTextValue(v string) Value {
	return Value{Type: TypeText, TextVal: v}
}

// NewBoolValue creates a new boolean value
func NewBoolValue(v bool) Value {
	return Value{Type: TypeBoolean, BoolVal: v}
}

// NewNullValue creates a new null value
func NewNullValue(dt DataType) Value {
	return Value{Type: dt, IsNull: true}
}

// String returns a string representation of the value
func (v Value) String() string {
	if v.IsNull {
		return "NULL"
	}
	switch v.Type {
	case TypeInteger:
		return fmt.Sprintf("%d", v.IntVal)
	case TypeText:
		return v.TextVal
	case TypeBoolean:
		if v.BoolVal {
			return "TRUE"
		}
		return "FALSE"
	default:
		return "UNKNOWN"
	}
}

// Table represents a table with schema and data
type Table struct {
	Name   string
	Schema Schema
	Rows   []Row
}

// Result represents the result of a query
type Result struct {
	Columns []string
	Rows    [][]string
	Message string
}
