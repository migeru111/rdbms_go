package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/migeru111/rdbms_go/pkg/types"
)

// Parser parses SQL statements
type Parser struct {
	lexer   *Lexer
	curTok  Token
	peekTok Token
}

// NewParser creates a new parser
func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(literal string) bool {
	return strings.ToUpper(p.curTok.Literal) == literal
}

func (p *Parser) expectToken(literal string) error {
	if !p.curTokenIs(literal) {
		return fmt.Errorf("expected %s, got %s", literal, p.curTok.Literal)
	}
	p.nextToken()
	return nil
}

// Parse parses a SQL statement
func (p *Parser) Parse() (Statement, error) {
	switch strings.ToUpper(p.curTok.Literal) {
	case "CREATE":
		return p.parseCreateTable()
	case "INSERT":
		return p.parseInsert()
	case "SELECT":
		return p.parseSelect()
	case "DELETE":
		return p.parseDelete()
	case "UPDATE":
		return p.parseUpdate()
	default:
		return nil, fmt.Errorf("unknown statement: %s", p.curTok.Literal)
	}
}

func (p *Parser) parseCreateTable() (*CreateTableStatement, error) {
	if err := p.expectToken("CREATE"); err != nil {
		return nil, err
	}
	if err := p.expectToken("TABLE"); err != nil {
		return nil, err
	}

	stmt := &CreateTableStatement{}
	stmt.TableName = p.curTok.Literal
	p.nextToken()

	if err := p.expectToken("("); err != nil {
		return nil, err
	}

	for {
		col := ColumnDef{}
		col.Name = p.curTok.Literal
		p.nextToken()

		dt, err := types.ParseDataType(strings.ToUpper(p.curTok.Literal))
		if err != nil {
			return nil, err
		}
		col.DataType = dt
		p.nextToken()

		stmt.Columns = append(stmt.Columns, col)

		if p.curTok.Type == TokenComma {
			p.nextToken()
			continue
		}
		break
	}

	if err := p.expectToken(")"); err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) parseInsert() (*InsertStatement, error) {
	if err := p.expectToken("INSERT"); err != nil {
		return nil, err
	}
	if err := p.expectToken("INTO"); err != nil {
		return nil, err
	}

	stmt := &InsertStatement{}
	stmt.TableName = p.curTok.Literal
	p.nextToken()

	// Optional column list
	if p.curTok.Type == TokenLParen {
		p.nextToken()
		for {
			stmt.Columns = append(stmt.Columns, p.curTok.Literal)
			p.nextToken()
			if p.curTok.Type == TokenComma {
				p.nextToken()
				continue
			}
			break
		}
		if err := p.expectToken(")"); err != nil {
			return nil, err
		}
	}

	if err := p.expectToken("VALUES"); err != nil {
		return nil, err
	}
	if err := p.expectToken("("); err != nil {
		return nil, err
	}

	for {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Values = append(stmt.Values, expr)

		if p.curTok.Type == TokenComma {
			p.nextToken()
			continue
		}
		break
	}

	if err := p.expectToken(")"); err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) parseSelect() (*SelectStatement, error) {
	if err := p.expectToken("SELECT"); err != nil {
		return nil, err
	}

	stmt := &SelectStatement{}

	// Parse columns
	if p.curTok.Type == TokenStar {
		stmt.Columns = []string{"*"}
		p.nextToken()
	} else {
		for {
			stmt.Columns = append(stmt.Columns, p.curTok.Literal)
			p.nextToken()
			if p.curTok.Type == TokenComma {
				p.nextToken()
				continue
			}
			break
		}
	}

	if err := p.expectToken("FROM"); err != nil {
		return nil, err
	}

	stmt.TableName = p.curTok.Literal
	p.nextToken()

	// Optional WHERE clause
	if p.curTokenIs("WHERE") {
		p.nextToken()
		where, err := p.parseWhereExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	return stmt, nil
}

func (p *Parser) parseDelete() (*DeleteStatement, error) {
	if err := p.expectToken("DELETE"); err != nil {
		return nil, err
	}
	if err := p.expectToken("FROM"); err != nil {
		return nil, err
	}

	stmt := &DeleteStatement{}
	stmt.TableName = p.curTok.Literal
	p.nextToken()

	// Optional WHERE clause
	if p.curTokenIs("WHERE") {
		p.nextToken()
		where, err := p.parseWhereExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	return stmt, nil
}

func (p *Parser) parseUpdate() (*UpdateStatement, error) {
	if err := p.expectToken("UPDATE"); err != nil {
		return nil, err
	}

	stmt := &UpdateStatement{}
	stmt.TableName = p.curTok.Literal
	p.nextToken()

	if err := p.expectToken("SET"); err != nil {
		return nil, err
	}

	// Parse SET clauses
	for {
		set := SetClause{}
		set.Column = p.curTok.Literal
		p.nextToken()

		if err := p.expectToken("="); err != nil {
			return nil, err
		}

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		set.Value = expr
		stmt.Sets = append(stmt.Sets, set)

		if p.curTok.Type == TokenComma {
			p.nextToken()
			continue
		}
		break
	}

	// Optional WHERE clause
	if p.curTokenIs("WHERE") {
		p.nextToken()
		where, err := p.parseWhereExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	return stmt, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	switch p.curTok.Type {
	case TokenNumber:
		val, err := strconv.ParseInt(p.curTok.Literal, 10, 64)
		if err != nil {
			return nil, err
		}
		p.nextToken()
		return &IntegerLiteral{Value: val}, nil
	case TokenString:
		val := p.curTok.Literal
		p.nextToken()
		return &StringLiteral{Value: val}, nil
	case TokenKeyword:
		switch strings.ToUpper(p.curTok.Literal) {
		case "TRUE":
			p.nextToken()
			return &BooleanLiteral{Value: true}, nil
		case "FALSE":
			p.nextToken()
			return &BooleanLiteral{Value: false}, nil
		case "NULL":
			p.nextToken()
			return &NullLiteral{}, nil
		}
	case TokenIdentifier:
		name := p.curTok.Literal
		p.nextToken()
		return &Identifier{Name: name}, nil
	}
	return nil, fmt.Errorf("unexpected token: %s", p.curTok.Literal)
}

func (p *Parser) parseWhereExpression() (Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.curTokenIs("AND") || p.curTokenIs("OR") {
		op := strings.ToUpper(p.curTok.Literal)
		p.nextToken()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

func (p *Parser) parseComparison() (Expression, error) {
	left, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	var op string
	switch p.curTok.Type {
	case TokenEquals:
		op = "="
	case TokenNotEquals:
		op = "<>"
	case TokenLessThan:
		op = "<"
	case TokenGreaterThan:
		op = ">"
	case TokenLessEquals:
		op = "<="
	case TokenGreaterEquals:
		op = ">="
	default:
		return left, nil
	}

	p.nextToken()
	right, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ComparisonExpr{Left: left, Operator: op, Right: right}, nil
}
