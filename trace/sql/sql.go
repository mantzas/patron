package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
)

type connInfo struct {
	instance, user string
}

// Conn represents a single database connection.
type Conn struct {
	connInfo
	conn *sql.Conn
}

func (c *Conn) startSpan(
	ctx context.Context,
	opName, stmt string,
	tags ...opentracing.Tag,
) (opentracing.Span, context.Context) {
	return trace.StartSQLSpan(ctx, opName, "sql", "rdbms", c.instance, c.user, stmt)
}

// BeginTx starts a transaction.
func (c *Conn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	sp, _ := c.startSpan(ctx, "conn.BeginTx", "")
	tx, err := c.conn.BeginTx(ctx, opts)
	if err != nil {
		trace.FinishSpanWithError(sp)
		return nil, err
	}

	trace.FinishSpanWithSuccess(sp)
	return &Tx{tx: tx}, nil
}

// Close returns the connection to the connection pool.
func (c *Conn) Close(ctx context.Context) error {
	sp, _ := c.startSpan(ctx, "conn.Close", "")
	err := c.conn.Close()
	if err != nil {
		trace.FinishSpanWithError(sp)
		return err
	}
	trace.FinishSpanWithSuccess(sp)
	return nil
}

// ExecContext executes a query without returning any rows.
func (c *Conn) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	sp, _ := c.startSpan(ctx, "conn.ExecContext", query)
	res, err := c.conn.ExecContext(ctx, query, args...)
	if err != nil {
		trace.FinishSpanWithError(sp)
		return nil, err
	}
	trace.FinishSpanWithSuccess(sp)
	return res, nil
}

// PingContext verifies the connection to the database is still alive.
func (c *Conn) PingContext(ctx context.Context) error {
	sp, _ := c.startSpan(ctx, "conn.PingContext", "")
	err := c.conn.PingContext(ctx)
	if err != nil {
		trace.FinishSpanWithError(sp)
		return err
	}
	trace.FinishSpanWithSuccess(sp)
	return nil
}

// PrepareContext creates a prepared statement for later queries or executions.
func (c *Conn) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	sp, _ := c.startSpan(ctx, "conn.PrepareContext", query)
	stmt, err := c.conn.PrepareContext(ctx, query)
	if err != nil {
		trace.FinishSpanWithError(sp)
		return nil, err
	}
	trace.FinishSpanWithSuccess(sp)
	return &Stmt{stmt: stmt}, nil
}

// QueryContext executes a query that returns rows.
func (c *Conn) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	sp, _ := c.startSpan(ctx, "conn.QueryContext", query)
	rows, err := c.conn.QueryContext(ctx, query, args...)
	if err != nil {
		trace.FinishSpanWithError(sp)
		return nil, err
	}
	trace.FinishSpanWithSuccess(sp)
	return rows, nil
}

// QueryRowContext executes a query that is expected to return at most one row.
func (c *Conn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	sp, _ := c.startSpan(ctx, "conn.QueryRowContext", query)
	row := c.conn.QueryRowContext(ctx, query, args...)
	trace.FinishSpanWithSuccess(sp)
	return row
}

// DB contains the underlying db to be traced.
type DB struct {
	connInfo
	db *sql.DB
}

func (db *DB) startSpan(
	ctx context.Context,
	opName, stmt string,
	tags ...opentracing.Tag,
) (opentracing.Span, context.Context) {
	return trace.StartSQLSpan(ctx, opName, "sql", "rdbms", db.instance, db.user, stmt)
}

// Open opens a database.
func Open(ctx context.Context, driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

// OpenDB opens a database.
func OpenDB(ctx context.Context, c driver.Connector) *DB {
	db := sql.OpenDB(c)
	return &DB{db: db}
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
func (db *DB) Close(ctx context.Context) error {
	return db.db.Close()
}

// Conn returns a connection.
func (db *DB) Conn(ctx context.Context) (*Conn, error) {
	conn, err := db.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &Conn{conn: conn, connInfo: db.connInfo}, nil
}

// Driver returns the database's underlying driver.
func (db *DB) Driver(ctx context.Context) driver.Driver {
	return db.db.Driver()
}

// ExecContext executes a query without returning any rows.
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// PingContext verifies a connection to the database is still alive, establishing a connection if necessary.
func (db *DB) PingContext(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

// PrepareContext creates a prepared statement for later queries or executions.
func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := db.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt}, nil
}

// QueryContext executes a query that returns rows.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
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
func (db *DB) Stats(ctx context.Context) sql.DBStats {
	return db.db.Stats()
}

// Stmt is a prepared statement.
type Stmt struct {
	connInfo
	stmt *sql.Stmt
}

// Close closes the statement.
func (s *Stmt) Close(ctx context.Context) error {
	return s.stmt.Close()
}

// ExecContext executes a prepared statement.
func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return s.stmt.ExecContext(ctx, args...)
}

// QueryContext executes a prepared query statement.
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	return s.stmt.QueryContext(ctx, args...)
}

// QueryRowContext executes a prepared query statement.
func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	return s.stmt.QueryRowContext(ctx, args...)
}

// Tx is an in-progress database transaction.
type Tx struct {
	connInfo
	tx *sql.Tx
}

// Commit commits the transaction.
func (tx *Tx) Commit(ctx context.Context) error {
	return tx.tx.Commit()
}

// ExecContext executes a query that doesn't return rows.
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, args...)
}

// PrepareContext creates a prepared statement for use within a transaction.
func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := tx.tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt}, nil
}

// QueryContext executes a query that returns rows.
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return tx.tx.QueryRowContext(ctx, query, args...)
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

// StmtContext returns a transaction-specific prepared statement from an existing statement.
func (tx *Tx) StmtContext(ctx context.Context, stmt *Stmt) *Stmt {
	return &Stmt{stmt: tx.tx.StmtContext(ctx, stmt.stmt)}
}
