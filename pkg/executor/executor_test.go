package executor

import (
	"testing"

	"github.com/migeru111/rdbms_go/pkg/parser"
	"github.com/migeru111/rdbms_go/pkg/storage"
	"github.com/migeru111/rdbms_go/pkg/types"
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

func TestTransactionCommit(t *testing.T) {
	db := storage.NewMemoryStorage()
	exec := NewExecutor(db)

	// Setup
	executeQuery(exec, "CREATE TABLE accounts (id INTEGER, balance INTEGER)")
	executeQuery(exec, "INSERT INTO accounts VALUES (1, 1000)")

	// Begin transaction
	result, err := executeQuery(exec, "BEGIN")
	if err != nil {
		t.Fatalf("begin error: %v", err)
	}
	if result.Message != "Transaction started" {
		t.Errorf("unexpected message: %s", result.Message)
	}

	// Update within transaction
	executeQuery(exec, "UPDATE accounts SET balance = 900 WHERE id = 1")

	// Commit
	result, err = executeQuery(exec, "COMMIT")
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}
	if result.Message != "Transaction committed" {
		t.Errorf("unexpected message: %s", result.Message)
	}

	// Verify changes persisted
	result, _ = executeQuery(exec, "SELECT balance FROM accounts WHERE id = 1")
	if len(result.Rows) != 1 || result.Rows[0][0] != "900" {
		t.Errorf("expected balance 900, got %v", result.Rows)
	}
}

func TestTransactionRollback(t *testing.T) {
	db := storage.NewMemoryStorage()
	exec := NewExecutor(db)

	// Setup
	executeQuery(exec, "CREATE TABLE accounts (id INTEGER, balance INTEGER)")
	executeQuery(exec, "INSERT INTO accounts VALUES (1, 1000)")

	// Begin transaction
	executeQuery(exec, "BEGIN TRANSACTION")

	// Make changes
	executeQuery(exec, "UPDATE accounts SET balance = 500 WHERE id = 1")
	executeQuery(exec, "INSERT INTO accounts VALUES (2, 2000)")

	// Verify changes within transaction
	result, _ := executeQuery(exec, "SELECT * FROM accounts")
	if len(result.Rows) != 2 {
		t.Errorf("expected 2 rows within transaction, got %d", len(result.Rows))
	}

	// Rollback
	result, err := executeQuery(exec, "ROLLBACK")
	if err != nil {
		t.Fatalf("rollback error: %v", err)
	}
	if result.Message != "Transaction rolled back" {
		t.Errorf("unexpected message: %s", result.Message)
	}

	// Verify changes were rolled back
	result, _ = executeQuery(exec, "SELECT balance FROM accounts WHERE id = 1")
	if len(result.Rows) != 1 || result.Rows[0][0] != "1000" {
		t.Errorf("expected balance 1000 after rollback, got %v", result.Rows)
	}

	// Verify inserted row was rolled back
	result, _ = executeQuery(exec, "SELECT * FROM accounts")
	if len(result.Rows) != 1 {
		t.Errorf("expected 1 row after rollback, got %d", len(result.Rows))
	}
}

func TestTransactionErrors(t *testing.T) {
	db := storage.NewMemoryStorage()
	exec := NewExecutor(db)

	// Commit without begin
	_, err := executeQuery(exec, "COMMIT")
	if err == nil {
		t.Error("expected error for commit without begin")
	}

	// Rollback without begin
	_, err = executeQuery(exec, "ROLLBACK")
	if err == nil {
		t.Error("expected error for rollback without begin")
	}

	// Nested transaction
	executeQuery(exec, "BEGIN")
	_, err = executeQuery(exec, "BEGIN")
	if err == nil {
		t.Error("expected error for nested begin")
	}
	executeQuery(exec, "ROLLBACK")
}

func executeQuery(exec *Executor, query string) (*types.Result, error) {
	p := parser.NewParser(query)
	stmt, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return exec.Execute(stmt)
}
