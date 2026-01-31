# RDBMS-Go

A minimal relational database management system written in Go.

## Features

- **SQL Parser**: Supports basic SQL statements
  - `CREATE TABLE`
  - `INSERT INTO`
  - `SELECT` (with `WHERE` clause)
  - `UPDATE`
  - `DELETE`

- **Data Types**:
  - `INTEGER` / `INT`
  - `TEXT` / `VARCHAR` / `STRING`
  - `BOOLEAN` / `BOOL`

- **WHERE Clause**:
  - Comparison operators: `=`, `<>`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical operators: `AND`, `OR`

- **In-Memory Storage**: Thread-safe storage engine

- **Interactive REPL**: Command-line interface for SQL queries

## Installation

```bash
go build -o rdbms ./cmd/rdbms
```

## Usage

```bash
./rdbms
```

### Example Session

```sql
rdbms> CREATE TABLE users (id INTEGER, name TEXT, active BOOLEAN)
Table users created

rdbms> INSERT INTO users VALUES (1, 'Alice', TRUE)
1 row inserted

rdbms> INSERT INTO users VALUES (2, 'Bob', FALSE)
1 row inserted

rdbms> SELECT * FROM users
+----+-------+--------+
| id | name  | active |
+----+-------+--------+
| 1  | Alice | TRUE   |
| 2  | Bob   | FALSE  |
+----+-------+--------+
2 row(s)

rdbms> SELECT name FROM users WHERE active = TRUE
+-------+
| name  |
+-------+
| Alice |
+-------+
1 row(s)

rdbms> UPDATE users SET active = TRUE WHERE id = 2
1 row(s) updated

rdbms> DELETE FROM users WHERE id = 1
1 row(s) deleted

rdbms> exit
Goodbye!
```

## Project Structure

```
.
├── cmd/rdbms/          # Main application
│   └── main.go
├── pkg/
│   ├── types/          # Data types and schema definitions
│   │   └── types.go
│   ├── parser/         # SQL lexer and parser
│   │   ├── lexer.go
│   │   ├── ast.go
│   │   ├── parser.go
│   │   └── parser_test.go
│   ├── storage/        # Storage engine
│   │   └── storage.go
│   └── executor/       # Query executor
│       ├── executor.go
│       └── executor_test.go
├── go.mod
├── ISSUES.md           # Future enhancement plans
└── README.md
```

## Testing

```bash
go test ./...
```

## Future Enhancements

See [ISSUES.md](ISSUES.md) for planned features:

1. Disk-based persistent storage
2. B-Tree index support
3. JOIN operations
4. Transaction support (ACID)
5. Aggregate functions (COUNT, SUM, AVG, etc.)
6. ORDER BY and LIMIT
7. Primary Key and Foreign Key constraints
8. NOT NULL and UNIQUE constraints
9. Subqueries
10. Query optimizer

## License

MIT
