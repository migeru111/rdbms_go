package storage

import (
	"fmt"
	"sync"

	"github.com/migeru111/rdbms_go/pkg/types"
)

// Storage interface defines storage operations
type Storage interface {
	CreateTable(name string, schema types.Schema) error
	GetTable(name string) (*types.Table, error)
	TableExists(name string) bool
	InsertRow(tableName string, row types.Row) error
	GetRows(tableName string) ([]types.Row, error)
	DeleteRows(tableName string, indices []int) error
	UpdateRows(tableName string, indices []int, updates map[int]types.Value) error
}

// MemoryStorage implements in-memory storage
type MemoryStorage struct {
	mu     sync.RWMutex
	tables map[string]*types.Table
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tables: make(map[string]*types.Table),
	}
}

// CreateTable creates a new table
func (s *MemoryStorage) CreateTable(name string, schema types.Schema) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tables[name]; exists {
		return fmt.Errorf("table %s already exists", name)
	}

	s.tables[name] = &types.Table{
		Name:   name,
		Schema: schema,
		Rows:   []types.Row{},
	}
	return nil
}

// GetTable returns a table by name
func (s *MemoryStorage) GetTable(name string) (*types.Table, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	table, exists := s.tables[name]
	if !exists {
		return nil, fmt.Errorf("table %s not found", name)
	}
	return table, nil
}

// TableExists checks if a table exists
func (s *MemoryStorage) TableExists(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.tables[name]
	return exists
}

// InsertRow inserts a row into a table
func (s *MemoryStorage) InsertRow(tableName string, row types.Row) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, exists := s.tables[tableName]
	if !exists {
		return fmt.Errorf("table %s not found", tableName)
	}

	if len(row.Values) != len(table.Schema.Columns) {
		return fmt.Errorf("column count mismatch: expected %d, got %d",
			len(table.Schema.Columns), len(row.Values))
	}

	table.Rows = append(table.Rows, row)
	return nil
}

// GetRows returns all rows from a table
func (s *MemoryStorage) GetRows(tableName string) ([]types.Row, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	table, exists := s.tables[tableName]
	if !exists {
		return nil, fmt.Errorf("table %s not found", tableName)
	}
	return table.Rows, nil
}

// DeleteRows deletes rows at specified indices
func (s *MemoryStorage) DeleteRows(tableName string, indices []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, exists := s.tables[tableName]
	if !exists {
		return fmt.Errorf("table %s not found", tableName)
	}

	// Create a set of indices to delete
	toDelete := make(map[int]bool)
	for _, idx := range indices {
		toDelete[idx] = true
	}

	// Build new rows slice excluding deleted indices
	newRows := make([]types.Row, 0, len(table.Rows)-len(indices))
	for i, row := range table.Rows {
		if !toDelete[i] {
			newRows = append(newRows, row)
		}
	}
	table.Rows = newRows
	return nil
}

// UpdateRows updates rows at specified indices
func (s *MemoryStorage) UpdateRows(tableName string, indices []int, updates map[int]types.Value) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, exists := s.tables[tableName]
	if !exists {
		return fmt.Errorf("table %s not found", tableName)
	}

	for _, rowIdx := range indices {
		if rowIdx < 0 || rowIdx >= len(table.Rows) {
			continue
		}
		for colIdx, val := range updates {
			if colIdx >= 0 && colIdx < len(table.Rows[rowIdx].Values) {
				table.Rows[rowIdx].Values[colIdx] = val
			}
		}
	}
	return nil
}
