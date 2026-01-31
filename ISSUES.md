# Future Enhancement Issues

以下はRDBMS-Goの将来の拡張機能として計画されているissueです。

---

## 1. Feature: Disk-based persistent storage

### Description
Implement disk-based persistent storage to save data between sessions.

### Requirements
- Save table schemas and data to disk files
- Load data on startup
- Support for WAL (Write-Ahead Logging) for crash recovery
- File format design (e.g., page-based storage)

### Benefits
- Data persistence across restarts
- Larger datasets than memory allows

---

## 2. Feature: B-Tree index support

### Description
Implement B-Tree indexes for faster query performance.

### Requirements
- CREATE INDEX statement
- Automatic index usage in WHERE clauses
- Support for primary key indexes
- Index maintenance on INSERT/UPDATE/DELETE

### Benefits
- O(log n) lookup instead of O(n) full table scans
- Improved query performance for large tables

---

## 3. Feature: JOIN operations

### Description
Implement JOIN operations for querying multiple tables.

### Requirements
- INNER JOIN
- LEFT JOIN / RIGHT JOIN
- CROSS JOIN
- Multiple table joins
- Join condition optimization

### Example
```sql
SELECT users.name, orders.total
FROM users
INNER JOIN orders ON users.id = orders.user_id
```

---

## 4. Feature: Transaction support (ACID)

### Description
Implement transaction support with ACID properties.

### Requirements
- BEGIN TRANSACTION / COMMIT / ROLLBACK statements
- Atomicity: All or nothing
- Consistency: Data integrity maintained
- Isolation: Concurrent transaction handling
- Durability: Committed transactions persist

### Example
```sql
BEGIN TRANSACTION;
INSERT INTO accounts VALUES (1, 1000);
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
COMMIT;
```

---

## 5. Feature: Aggregate functions

### Description
Implement aggregate functions for data analysis.

### Requirements
- COUNT()
- SUM()
- AVG()
- MAX()
- MIN()
- GROUP BY clause
- HAVING clause

### Example
```sql
SELECT department, COUNT(*), AVG(salary)
FROM employees
GROUP BY department
HAVING COUNT(*) > 5
```

---

## 6. Feature: ORDER BY and LIMIT

### Description
Implement result ordering and limiting.

### Requirements
- ORDER BY column [ASC|DESC]
- Multiple column ordering
- LIMIT n
- OFFSET n

### Example
```sql
SELECT * FROM users ORDER BY created_at DESC LIMIT 10 OFFSET 20
```

---

## 7. Feature: Primary Key and Foreign Key constraints

### Description
Implement referential integrity constraints.

### Requirements
- PRIMARY KEY constraint (unique, not null)
- FOREIGN KEY constraint with references
- ON DELETE CASCADE / SET NULL
- ON UPDATE CASCADE
- Constraint validation on INSERT/UPDATE/DELETE

### Example
```sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
)
```

---

## 8. Feature: NOT NULL and UNIQUE constraints

### Description
Implement column-level constraints.

### Requirements
- NOT NULL constraint
- UNIQUE constraint
- Constraint validation on INSERT/UPDATE
- Error messages for constraint violations

### Example
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL
)
```

---

## 9. Feature: Subqueries

### Description
Implement support for nested queries.

### Requirements
- Subqueries in WHERE clause
- Subqueries in FROM clause (derived tables)
- Correlated subqueries
- IN, EXISTS, ANY, ALL operators

### Example
```sql
SELECT * FROM users
WHERE id IN (SELECT user_id FROM orders WHERE total > 100)
```

---

## 10. Feature: Query optimizer

### Description
Implement a query optimizer for better performance.

### Requirements
- Query plan generation
- Cost-based optimization
- Index selection
- Join order optimization
- EXPLAIN statement

### Example
```sql
EXPLAIN SELECT * FROM users WHERE id = 1
```
