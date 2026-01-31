package executor

import (
	"testing"

	"github.com/migeru111/rdbms_go/pkg/parser"
	"github.com/migeru111/rdbms_go/pkg/storage"
)

func TestCreateTableAndInsert(t *testing.T) {
	db := storage.NewMemoryStorage()
	exec := NewExecutor(db)

	// Create table
	p := parser.NewParser("CREATE TABLE users (id INTEGER, name TEXT)")
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := exec.Execute(stmt)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}

	if result.Message != "Table users created" {
		t.Errorf("unexpected message: %s", result.Message)
	}

	// Insert row
	p = parser.NewParser("INSERT INTO users VALUES (1, 'Alice')")
	stmt, err = p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err = exec.Execute(stmt)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}

	if result.Message != "1 row inserted" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

func TestSelectAll(t *testing.T) {
	db := storage.NewMemoryStorage()
	exec := NewExecutor(db)

	// Create and populate table
	queries := []string{
		"CREATE TABLE users (id INTEGER, name TEXT)",
		"INSERT INTO users VALUES (1, 'Alice')",
		"INSERT INTO users VALUES (2, 'Bob')",
	}

	for _, q := range queries {
		p := parser.NewParser(q)
		stmt, _ := p.Parse()
		exec.Execute(stmt)
	}

	// Select all
	p := parser.NewParser("SELECT * FROM users")
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := exec.Execute(stmt)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}

	if len(result.Columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(result.Columns))
	}

	if len(result.Rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(result.Rows))
	}
}

func TestSelectWithWhere(t *testing.T) {
	db := storage.NewMemoryStorage()
	exec := NewExecutor(db)

	// Create and populate table
	queries := []string{
		"CREATE TABLE users (id INTEGER, name TEXT)",
		"INSERT INTO users VALUES (1, 'Alice')",
		"INSERT INTO users VALUES (2, 'Bob')",
		"INSERT INTO users VALUES (3, 'Charlie')",
	}

	for _, q := range queries {
		p := parser.NewParser(q)
		stmt, _ := p.Parse()
		exec.Execute(stmt)
	}

	// Select with WHERE
	p := parser.NewParser("SELECT name FROM users WHERE id = 2")
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := exec.Execute(stmt)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}

	if result.Rows[0][0] != "Bob" {
		t.Errorf("expected 'Bob', got '%s'", result.Rows[0][0])
	}
}

func TestDelete(t *testing.T) {
	db := storage.NewMemoryStorage()
	exec := NewExecutor(db)

	// Create and populate table
	queries := []string{
		"CREATE TABLE users (id INTEGER, name TEXT)",
		"INSERT INTO users VALUES (1, 'Alice')",
		"INSERT INTO users VALUES (2, 'Bob')",
	}

	for _, q := range queries {
		p := parser.NewParser(q)
		stmt, _ := p.Parse()
		exec.Execute(stmt)
	}

	// Delete
	p := parser.NewParser("DELETE FROM users WHERE id = 1")
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := exec.Execute(stmt)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}

	if result.Message != "1 row(s) deleted" {
		t.Errorf("unexpected message: %s", result.Message)
	}

	// Verify deletion
	p = parser.NewParser("SELECT * FROM users")
	stmt, _ = p.Parse()
	result, _ = exec.Execute(stmt)

	if len(result.Rows) != 1 {
		t.Errorf("expected 1 row after delete, got %d", len(result.Rows))
	}
}

func TestUpdate(t *testing.T) {
	db := storage.NewMemoryStorage()
	exec := NewExecutor(db)

	// Create and populate table
	queries := []string{
		"CREATE TABLE users (id INTEGER, name TEXT)",
		"INSERT INTO users VALUES (1, 'Alice')",
		"INSERT INTO users VALUES (2, 'Bob')",
	}

	for _, q := range queries {
		p := parser.NewParser(q)
		stmt, _ := p.Parse()
		exec.Execute(stmt)
	}

	// Update
	p := parser.NewParser("UPDATE users SET name = 'Alicia' WHERE id = 1")
	stmt, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := exec.Execute(stmt)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}

	if result.Message != "1 row(s) updated" {
		t.Errorf("unexpected message: %s", result.Message)
	}

	// Verify update
	p = parser.NewParser("SELECT name FROM users WHERE id = 1")
	stmt, _ = p.Parse()
	result, _ = exec.Execute(stmt)

	if len(result.Rows) != 1 || result.Rows[0][0] != "Alicia" {
		t.Errorf("expected 'Alicia', got %v", result.Rows)
	}
}
