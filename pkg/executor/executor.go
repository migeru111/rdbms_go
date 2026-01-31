package executor

import (
	"fmt"

	"github.com/migeru111/rdbms_go/pkg/parser"
	"github.com/migeru111/rdbms_go/pkg/storage"
	"github.com/migeru111/rdbms_go/pkg/types"
)

// Executor executes SQL statements
type Executor struct {
	storage storage.Storage
}

// NewExecutor creates a new executor
func NewExecutor(s storage.Storage) *Executor {
	return &Executor{storage: s}
}

// Execute executes a SQL statement
func (e *Executor) Execute(stmt parser.Statement) (*types.Result, error) {
	switch s := stmt.(type) {
	case *parser.CreateTableStatement:
		return e.executeCreateTable(s)
	case *parser.InsertStatement:
		return e.executeInsert(s)
	case *parser.SelectStatement:
		return e.executeSelect(s)
	case *parser.DeleteStatement:
		return e.executeDelete(s)
	case *parser.UpdateStatement:
		return e.executeUpdate(s)
	case *parser.BeginStatement:
		return e.executeBegin()
	case *parser.CommitStatement:
		return e.executeCommit()
	case *parser.RollbackStatement:
		return e.executeRollback()
	default:
		return nil, fmt.Errorf("unknown statement type")
	}
}

func (e *Executor) executeBegin() (*types.Result, error) {
	if err := e.storage.Begin(); err != nil {
		return nil, err
	}
	return &types.Result{
		Message: "Transaction started",
	}, nil
}

func (e *Executor) executeCommit() (*types.Result, error) {
	if err := e.storage.Commit(); err != nil {
		return nil, err
	}
	return &types.Result{
		Message: "Transaction committed",
	}, nil
}

func (e *Executor) executeRollback() (*types.Result, error) {
	if err := e.storage.Rollback(); err != nil {
		return nil, err
	}
	return &types.Result{
		Message: "Transaction rolled back",
	}, nil
}

func (e *Executor) executeCreateTable(stmt *parser.CreateTableStatement) (*types.Result, error) {
	schema := types.Schema{
		Columns: make([]types.Column, len(stmt.Columns)),
	}
	for i, col := range stmt.Columns {
		schema.Columns[i] = types.Column{
			Name:     col.Name,
			DataType: col.DataType,
		}
	}

	if err := e.storage.CreateTable(stmt.TableName, schema); err != nil {
		return nil, err
	}

	return &types.Result{
		Message: fmt.Sprintf("Table %s created", stmt.TableName),
	}, nil
}

func (e *Executor) executeInsert(stmt *parser.InsertStatement) (*types.Result, error) {
	table, err := e.storage.GetTable(stmt.TableName)
	if err != nil {
		return nil, err
	}

	row := types.Row{
		Values: make([]types.Value, len(table.Schema.Columns)),
	}

	// If columns are specified, map values to correct positions
	if len(stmt.Columns) > 0 {
		// Initialize with null values
		for i, col := range table.Schema.Columns {
			row.Values[i] = types.NewNullValue(col.DataType)
		}

		// Map specified columns
		colIndex := make(map[string]int)
		for i, col := range table.Schema.Columns {
			colIndex[col.Name] = i
		}

		for i, colName := range stmt.Columns {
			idx, ok := colIndex[colName]
			if !ok {
				return nil, fmt.Errorf("unknown column: %s", colName)
			}
			val, err := e.exprToValue(stmt.Values[i], table.Schema.Columns[idx].DataType)
			if err != nil {
				return nil, err
			}
			row.Values[idx] = val
		}
	} else {
		// Values in order
		if len(stmt.Values) != len(table.Schema.Columns) {
			return nil, fmt.Errorf("value count mismatch")
		}
		for i, expr := range stmt.Values {
			val, err := e.exprToValue(expr, table.Schema.Columns[i].DataType)
			if err != nil {
				return nil, err
			}
			row.Values[i] = val
		}
	}

	if err := e.storage.InsertRow(stmt.TableName, row); err != nil {
		return nil, err
	}

	return &types.Result{
		Message: "1 row inserted",
	}, nil
}

func (e *Executor) executeSelect(stmt *parser.SelectStatement) (*types.Result, error) {
	table, err := e.storage.GetTable(stmt.TableName)
	if err != nil {
		return nil, err
	}

	rows, err := e.storage.GetRows(stmt.TableName)
	if err != nil {
		return nil, err
	}

	// Determine which columns to select
	var colIndices []int
	var columns []string

	if len(stmt.Columns) == 1 && stmt.Columns[0] == "*" {
		for i, col := range table.Schema.Columns {
			colIndices = append(colIndices, i)
			columns = append(columns, col.Name)
		}
	} else {
		colIndex := make(map[string]int)
		for i, col := range table.Schema.Columns {
			colIndex[col.Name] = i
		}
		for _, colName := range stmt.Columns {
			idx, ok := colIndex[colName]
			if !ok {
				return nil, fmt.Errorf("unknown column: %s", colName)
			}
			colIndices = append(colIndices, idx)
			columns = append(columns, colName)
		}
	}

	// Filter and project rows
	var resultRows [][]string
	for _, row := range rows {
		if stmt.Where != nil {
			match, err := e.evaluateWhere(stmt.Where, row, table.Schema)
			if err != nil {
				return nil, err
			}
			if !match {
				continue
			}
		}

		var resultRow []string
		for _, idx := range colIndices {
			resultRow = append(resultRow, row.Values[idx].String())
		}
		resultRows = append(resultRows, resultRow)
	}

	return &types.Result{
		Columns: columns,
		Rows:    resultRows,
	}, nil
}

func (e *Executor) executeDelete(stmt *parser.DeleteStatement) (*types.Result, error) {
	table, err := e.storage.GetTable(stmt.TableName)
	if err != nil {
		return nil, err
	}

	rows, err := e.storage.GetRows(stmt.TableName)
	if err != nil {
		return nil, err
	}

	var indicesToDelete []int
	for i, row := range rows {
		if stmt.Where != nil {
			match, err := e.evaluateWhere(stmt.Where, row, table.Schema)
			if err != nil {
				return nil, err
			}
			if match {
				indicesToDelete = append(indicesToDelete, i)
			}
		} else {
			indicesToDelete = append(indicesToDelete, i)
		}
	}

	if err := e.storage.DeleteRows(stmt.TableName, indicesToDelete); err != nil {
		return nil, err
	}

	return &types.Result{
		Message: fmt.Sprintf("%d row(s) deleted", len(indicesToDelete)),
	}, nil
}

func (e *Executor) executeUpdate(stmt *parser.UpdateStatement) (*types.Result, error) {
	table, err := e.storage.GetTable(stmt.TableName)
	if err != nil {
		return nil, err
	}

	rows, err := e.storage.GetRows(stmt.TableName)
	if err != nil {
		return nil, err
	}

	// Build column index map
	colIndex := make(map[string]int)
	for i, col := range table.Schema.Columns {
		colIndex[col.Name] = i
	}

	// Build updates map
	updates := make(map[int]types.Value)
	for _, set := range stmt.Sets {
		idx, ok := colIndex[set.Column]
		if !ok {
			return nil, fmt.Errorf("unknown column: %s", set.Column)
		}
		val, err := e.exprToValue(set.Value, table.Schema.Columns[idx].DataType)
		if err != nil {
			return nil, err
		}
		updates[idx] = val
	}

	var indicesToUpdate []int
	for i, row := range rows {
		if stmt.Where != nil {
			match, err := e.evaluateWhere(stmt.Where, row, table.Schema)
			if err != nil {
				return nil, err
			}
			if match {
				indicesToUpdate = append(indicesToUpdate, i)
			}
		} else {
			indicesToUpdate = append(indicesToUpdate, i)
		}
	}

	if err := e.storage.UpdateRows(stmt.TableName, indicesToUpdate, updates); err != nil {
		return nil, err
	}

	return &types.Result{
		Message: fmt.Sprintf("%d row(s) updated", len(indicesToUpdate)),
	}, nil
}

func (e *Executor) exprToValue(expr parser.Expression, dt types.DataType) (types.Value, error) {
	switch ex := expr.(type) {
	case *parser.IntegerLiteral:
		return types.NewIntValue(ex.Value), nil
	case *parser.StringLiteral:
		return types.NewTextValue(ex.Value), nil
	case *parser.BooleanLiteral:
		return types.NewBoolValue(ex.Value), nil
	case *parser.NullLiteral:
		return types.NewNullValue(dt), nil
	default:
		return types.Value{}, fmt.Errorf("unsupported expression type")
	}
}

func (e *Executor) evaluateWhere(expr parser.Expression, row types.Row, schema types.Schema) (bool, error) {
	switch ex := expr.(type) {
	case *parser.ComparisonExpr:
		return e.evaluateComparison(ex, row, schema)
	case *parser.BinaryExpr:
		left, err := e.evaluateWhere(ex.Left, row, schema)
		if err != nil {
			return false, err
		}
		right, err := e.evaluateWhere(ex.Right, row, schema)
		if err != nil {
			return false, err
		}
		switch ex.Operator {
		case "AND":
			return left && right, nil
		case "OR":
			return left || right, nil
		}
	}
	return false, fmt.Errorf("unsupported where expression")
}

func (e *Executor) evaluateComparison(expr *parser.ComparisonExpr, row types.Row, schema types.Schema) (bool, error) {
	leftVal, err := e.resolveExpr(expr.Left, row, schema)
	if err != nil {
		return false, err
	}
	rightVal, err := e.resolveExpr(expr.Right, row, schema)
	if err != nil {
		return false, err
	}

	if leftVal.IsNull || rightVal.IsNull {
		return false, nil
	}

	switch expr.Operator {
	case "=":
		return e.compareValues(leftVal, rightVal) == 0, nil
	case "<>", "!=":
		return e.compareValues(leftVal, rightVal) != 0, nil
	case "<":
		return e.compareValues(leftVal, rightVal) < 0, nil
	case ">":
		return e.compareValues(leftVal, rightVal) > 0, nil
	case "<=":
		return e.compareValues(leftVal, rightVal) <= 0, nil
	case ">=":
		return e.compareValues(leftVal, rightVal) >= 0, nil
	}
	return false, fmt.Errorf("unknown operator: %s", expr.Operator)
}

func (e *Executor) resolveExpr(expr parser.Expression, row types.Row, schema types.Schema) (types.Value, error) {
	switch ex := expr.(type) {
	case *parser.Identifier:
		for i, col := range schema.Columns {
			if col.Name == ex.Name {
				return row.Values[i], nil
			}
		}
		return types.Value{}, fmt.Errorf("unknown column: %s", ex.Name)
	case *parser.IntegerLiteral:
		return types.NewIntValue(ex.Value), nil
	case *parser.StringLiteral:
		return types.NewTextValue(ex.Value), nil
	case *parser.BooleanLiteral:
		return types.NewBoolValue(ex.Value), nil
	case *parser.NullLiteral:
		return types.NewNullValue(types.TypeText), nil
	}
	return types.Value{}, fmt.Errorf("unsupported expression")
}

func (e *Executor) compareValues(a, b types.Value) int {
	switch a.Type {
	case types.TypeInteger:
		if a.IntVal < b.IntVal {
			return -1
		} else if a.IntVal > b.IntVal {
			return 1
		}
		return 0
	case types.TypeText:
		if a.TextVal < b.TextVal {
			return -1
		} else if a.TextVal > b.TextVal {
			return 1
		}
		return 0
	case types.TypeBoolean:
		if a.BoolVal == b.BoolVal {
			return 0
		}
		if !a.BoolVal && b.BoolVal {
			return -1
		}
		return 1
	}
	return 0
}
