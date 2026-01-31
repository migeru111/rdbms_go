package parser

import "github.com/migeru111/rdbms_go/pkg/types"

// Statement represents a SQL statement
type Statement interface {
	statementNode()
}

// Expression represents an expression
type Expression interface {
	expressionNode()
}

// CreateTableStatement represents CREATE TABLE
type CreateTableStatement struct {
	TableName string
	Columns   []ColumnDef
}

func (s *CreateTableStatement) statementNode() {}

// ColumnDef represents a column definition
type ColumnDef struct {
	Name     string
	DataType types.DataType
}

// InsertStatement represents INSERT INTO
type InsertStatement struct {
	TableName string
	Columns   []string
	Values    []Expression
}

func (s *InsertStatement) statementNode() {}

// SelectStatement represents SELECT
type SelectStatement struct {
	Columns   []string
	TableName string
	Where     Expression
}

func (s *SelectStatement) statementNode() {}

// DeleteStatement represents DELETE
type DeleteStatement struct {
	TableName string
	Where     Expression
}

func (s *DeleteStatement) statementNode() {}

// UpdateStatement represents UPDATE
type UpdateStatement struct {
	TableName string
	Sets      []SetClause
	Where     Expression
}

func (s *UpdateStatement) statementNode() {}

// SetClause represents SET column = value
type SetClause struct {
	Column string
	Value  Expression
}

// Literal expressions
type IntegerLiteral struct {
	Value int64
}

func (e *IntegerLiteral) expressionNode() {}

type StringLiteral struct {
	Value string
}

func (e *StringLiteral) expressionNode() {}

type BooleanLiteral struct {
	Value bool
}

func (e *BooleanLiteral) expressionNode() {}

type NullLiteral struct{}

func (e *NullLiteral) expressionNode() {}

// Identifier represents a column reference
type Identifier struct {
	Name string
}

func (e *Identifier) expressionNode() {}

// ComparisonExpr represents a comparison expression
type ComparisonExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (e *ComparisonExpr) expressionNode() {}

// BinaryExpr represents AND/OR expressions
type BinaryExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (e *BinaryExpr) expressionNode() {}

// BeginStatement represents BEGIN TRANSACTION
type BeginStatement struct{}

func (s *BeginStatement) statementNode() {}

// CommitStatement represents COMMIT
type CommitStatement struct{}

func (s *CommitStatement) statementNode() {}

// RollbackStatement represents ROLLBACK
type RollbackStatement struct{}

func (s *RollbackStatement) statementNode() {}
