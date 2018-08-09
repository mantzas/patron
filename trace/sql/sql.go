package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

// Conn represents a single database connection.
type Conn struct {
	conn *sql.Conn
}

// BeginTx starts a transaction.
func (c *Conn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := c.conn.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

// Close returns the connection to the connection pool.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// ExecContext executes a query without returning any rows.
func (c *Conn) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.conn.ExecContext(ctx, query, args...)
}

// PingContext verifies the connection to the database is still alive.
func (c *Conn) PingContext(ctx context.Context) error {
	return c.conn.PingContext(ctx)
}

// PrepareContext creates a prepared statement for later queries or executions.
func (c *Conn) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := c.conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt}, nil
}

// QueryContext executes a query that returns rows.
func (c *Conn) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return c.conn.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
func (c *Conn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return c.conn.QueryRowContext(ctx, query, args...)
}

// DB contains the underlying db to be traced.
type DB struct {
	db *sql.DB
}

// Open opens a database.
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

// OpenDB opens a database.
func OpenDB(c driver.Connector) *DB {
	db := sql.OpenDB(c)
	return &DB{db: db}
}

// Begin starts a transaction.
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

// BeginTx starts a transaction.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

// Close closes the database, releasing any open resources.
func (db *DB) Close() error {
	return db.db.Close()
}

// Conn returns a connection.
func (db *DB) Conn(ctx context.Context) (*Conn, error) {
	conn, err := db.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &Conn{conn: conn}, nil
}

// Driver returns the database's underlying driver.
func (db *DB) Driver() driver.Driver {
	return db.db.Driver()
}

// Exec executes a query without returning any rows.
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

// ExecContext executes a query without returning any rows.
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
func (db *DB) Ping() error {
	return db.db.Ping()
}

// PingContext verifies a connection to the database is still alive, establishing a connection if necessary.
func (db *DB) PingContext(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

// Prepare creates a prepared statement for later queries or executions.
func (db *DB) Prepare(query string) (*Stmt, error) {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt}, nil
}

// PrepareContext creates a prepared statement for later queries or executions.
func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := db.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt}, nil
}

// Query executes a query that returns rows.
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

// QueryContext executes a query that returns rows.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.db.QueryRow(query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.db.QueryRowContext(ctx, query, args...)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (db *DB) SetConnMaxLifetime(d time.Duration) {
	db.db.SetConnMaxLifetime(d)
}

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
func (db *DB) SetMaxIdleConns(n int) {
	db.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
func (db *DB) SetMaxOpenConns(n int) {
	db.db.SetMaxOpenConns(n)
}

// Stats returns database statistics.
func (db *DB) Stats() sql.DBStats {
	return db.db.Stats()
}

// Stmt is a prepared statement.
type Stmt struct {
	stmt *sql.Stmt
}

// Close closes the statement.
func (s *Stmt) Close() error {
	return s.stmt.Close()
}

// Exec executes a prepared statement.
func (s *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.stmt.Exec(args...)
}

// ExecContext executes a prepared statement.
func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return s.stmt.ExecContext(ctx, args...)
}

// Query executes a prepared query statement.
func (s *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	return s.stmt.Query(args...)
}

// QueryContext executes a prepared query statement.
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	return s.stmt.QueryContext(ctx, args...)
}

// QueryRow executes a prepared query statement.
func (s *Stmt) QueryRow(args ...interface{}) *sql.Row {
	return s.stmt.QueryRow(args...)
}

// QueryRowContext executes a prepared query statement.
func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	return s.stmt.QueryRowContext(ctx, args...)
}

// Tx is an in-progress database transaction.
type Tx struct {
	tx *sql.Tx
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

// Exec executes a query that doesn't return rows.
func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.tx.Exec(query, args...)
}

// ExecContext executes a query that doesn't return rows.
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, args...)
}

// Prepare creates a prepared statement for use within a transaction.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	stmt, err := tx.tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt}, nil
}

// PrepareContext creates a prepared statement for use within a transaction.
func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := tx.tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt}, nil
}

// Query executes a query that returns rows.
func (tx *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.tx.Query(query, args...)
}

// QueryContext executes a query that returns rows.
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (tx *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	return tx.tx.QueryRow(query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return tx.tx.QueryRowContext(ctx, query, args...)
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

// Stmt returns a transaction-specific prepared statement from an existing statement.
func (tx *Tx) Stmt(stmt *Stmt) *Stmt {
	return &Stmt{stmt: tx.tx.Stmt(stmt.stmt)}
}

// StmtContext returns a transaction-specific prepared statement from an existing statement.
func (tx *Tx) StmtContext(ctx context.Context, stmt *Stmt) *Stmt {
	return &Stmt{stmt: tx.tx.StmtContext(ctx, stmt.stmt)}
}
