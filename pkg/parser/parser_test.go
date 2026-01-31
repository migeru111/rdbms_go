package parser

import (
	"testing"

	"github.com/migeru111/rdbms_go/pkg/types"
)

func TestParseCreateTable(t *testing.T) {
	input := "CREATE TABLE users (id INTEGER, name TEXT, active BOOLEAN)"
	p := NewParser(input)
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	create, ok := stmt.(*CreateTableStatement)
	if !ok {
		t.Fatalf("expected CreateTableStatement, got %T", stmt)
	}

	if create.TableName != "users" {
		t.Errorf("expected table name 'users', got '%s'", create.TableName)
	}

	if len(create.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(create.Columns))
	}

	expected := []struct {
		name     string
		dataType types.DataType
	}{
		{"id", types.TypeInteger},
		{"name", types.TypeText},
		{"active", types.TypeBoolean},
	}

	for i, exp := range expected {
		if create.Columns[i].Name != exp.name {
			t.Errorf("column %d: expected name '%s', got '%s'", i, exp.name, create.Columns[i].Name)
		}
		if create.Columns[i].DataType != exp.dataType {
			t.Errorf("column %d: expected type %v, got %v", i, exp.dataType, create.Columns[i].DataType)
		}
	}
}

func TestParseInsert(t *testing.T) {
	input := "INSERT INTO users VALUES (1, 'Alice', TRUE)"
	p := NewParser(input)
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insert, ok := stmt.(*InsertStatement)
	if !ok {
		t.Fatalf("expected InsertStatement, got %T", stmt)
	}

	if insert.TableName != "users" {
		t.Errorf("expected table name 'users', got '%s'", insert.TableName)
	}

	if len(insert.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(insert.Values))
	}
}

func TestParseSelect(t *testing.T) {
	input := "SELECT id, name FROM users WHERE id = 1"
	p := NewParser(input)
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	sel, ok := stmt.(*SelectStatement)
	if !ok {
		t.Fatalf("expected SelectStatement, got %T", stmt)
	}

	if sel.TableName != "users" {
		t.Errorf("expected table name 'users', got '%s'", sel.TableName)
	}

	if len(sel.Columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(sel.Columns))
	}

	if sel.Where == nil {
		t.Error("expected WHERE clause")
	}
}

func TestParseSelectStar(t *testing.T) {
	input := "SELECT * FROM users"
	p := NewParser(input)
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	sel, ok := stmt.(*SelectStatement)
	if !ok {
		t.Fatalf("expected SelectStatement, got %T", stmt)
	}

	if len(sel.Columns) != 1 || sel.Columns[0] != "*" {
		t.Errorf("expected ['*'], got %v", sel.Columns)
	}
}

func TestParseDelete(t *testing.T) {
	input := "DELETE FROM users WHERE id = 1"
	p := NewParser(input)
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	del, ok := stmt.(*DeleteStatement)
	if !ok {
		t.Fatalf("expected DeleteStatement, got %T", stmt)
	}

	if del.TableName != "users" {
		t.Errorf("expected table name 'users', got '%s'", del.TableName)
	}

	if del.Where == nil {
		t.Error("expected WHERE clause")
	}
}

func TestParseUpdate(t *testing.T) {
	input := "UPDATE users SET name = 'Bob' WHERE id = 1"
	p := NewParser(input)
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	upd, ok := stmt.(*UpdateStatement)
	if !ok {
		t.Fatalf("expected UpdateStatement, got %T", stmt)
	}

	if upd.TableName != "users" {
		t.Errorf("expected table name 'users', got '%s'", upd.TableName)
	}

	if len(upd.Sets) != 1 {
		t.Errorf("expected 1 SET clause, got %d", len(upd.Sets))
	}

	if upd.Where == nil {
		t.Error("expected WHERE clause")
	}
}
