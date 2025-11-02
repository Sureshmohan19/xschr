// XSchr - Experimental Scheduler in Go
// "Make it work, make it right, make it fast — in that order."

package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// DBConn wraps a SQLite connection with thread-safe access
// Inspired by SLURM's mysql_conn_t structure
type DBConn struct {
	db   *sql.DB
	mu   sync.Mutex
	path string
}

// OpenDB opens a SQLite database connection and initializes schema
func OpenDB(path string) (*DBConn, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Set busy timeout to handle SQLITE_BUSY
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}

	conn := &DBConn{
		db:   db,
		path: path,
	}

	// Initialize schema
	if err := conn.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return conn, nil
}

// Close closes the database connection
func (c *DBConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.db.Close()
}

// initSchema creates the schema_versions table and jobs table
func (c *DBConn) initSchema() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create schema_versions table (like SLURM's table_defs_table)
	query := SchemaVersionTableSchema.buildCreateTableSQL()
	if _, err := c.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create schema_versions table: %w", err)
	}

	// Create jobs table with auto-migration
	if err := c.createTableWithMigration(JobsTableSchema); err != nil {
		return fmt.Errorf("failed to create jobs table: %w", err)
	}

	return nil
}

// createTableWithMigration creates a table and handles schema migration
// This is equivalent to SLURM's mysql_db_create_table + _mysql_make_table_current
func (c *DBConn) createTableWithMigration(schema TableSchema) error {
	// Create the table if it doesn't exist
	createSQL := schema.buildCreateTableSQL()
	if _, err := c.db.Exec(createSQL); err != nil {
		return fmt.Errorf("failed to create table %s: %w", schema.Name, err)
	}

	// Check if schema has changed by comparing with schema_versions
	currentDef := schema.buildSchemaDefinition()

	var storedDef string
	err := c.db.QueryRow(
		"SELECT definition FROM schema_versions WHERE table_name = ?",
		schema.Name,
	).Scan(&storedDef)

	now := time.Now().Unix()

	if err == sql.ErrNoRows {
		// First time seeing this table, store the definition
		_, err = c.db.Exec(
			"INSERT INTO schema_versions (table_name, definition, created_at, updated_at) VALUES (?, ?, ?, ?)",
			schema.Name, currentDef, now, now,
		)
		return err
	} else if err != nil {
		return fmt.Errorf("failed to check schema version: %w", err)
	}

	// If definition changed, we would run ALTER TABLE here
	// For v1.0.0, we skip auto-migration and just update the definition
	// (Full migration logic can be added later when needed)
	if storedDef != currentDef {
		fmt.Printf("Schema changed for table %s (migration not yet implemented)\n", schema.Name)
		_, err = c.db.Exec(
			"UPDATE schema_versions SET definition = ?, updated_at = ? WHERE table_name = ?",
			currentDef, now, schema.Name,
		)
		return err
	}

	return nil
}

// Exec executes a query without returning results
func (c *DBConn) Exec(query string, args ...interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.db.Exec(query, args...)
	return err
}

// Query executes a query and returns rows
func (c *DBConn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.db.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (c *DBConn) QueryRow(query string, args ...interface{}) *sql.Row {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.db.QueryRow(query, args...)
}

// InsertReturnID executes an INSERT and returns the last inserted ID
// Equivalent to SLURM's mysql_insert_ret_id
func (c *DBConn) InsertReturnID(query string, args ...interface{}) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	result, err := c.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
