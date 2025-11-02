// XSchr - Experimental Scheduler in Go
// "Make it work, make it right, make it fast — in that order."

package db

// Field represents a single column in a table schema
type Field struct {
	Name    string
	Type    string
	Options string
}

// TableSchema defines the structure of a database table
type TableSchema struct {
	Name   string
	Fields []Field
	Extra  string // PRIMARY KEY, indexes, etc.
}

// JobsTableSchema defines the schema for the jobs table
// Inspired by SLURM's table definition approach
var JobsTableSchema = TableSchema{
	Name: "jobs",
	Fields: []Field{
		{Name: "id", Type: "INTEGER", Options: "PRIMARY KEY AUTOINCREMENT"},
		{Name: "command", Type: "TEXT", Options: "NOT NULL"},
		{Name: "cpus", Type: "INTEGER", Options: "NOT NULL"},
		{Name: "state", Type: "TEXT", Options: "NOT NULL"},
		{Name: "pid", Type: "INTEGER", Options: ""},
		{Name: "submitted_at", Type: "INTEGER", Options: "NOT NULL"},
		{Name: "started_at", Type: "INTEGER", Options: ""},
		{Name: "finished_at", Type: "INTEGER", Options: ""},
		{Name: "exit_code", Type: "INTEGER", Options: ""},
	},
	Extra: "",
}

// SchemaVersionTableSchema defines the meta-table that tracks schema versions
// This is equivalent to SLURM's table_defs_table
var SchemaVersionTableSchema = TableSchema{
	Name: "schema_versions",
	Fields: []Field{
		{Name: "table_name", Type: "TEXT", Options: "PRIMARY KEY"},
		{Name: "definition", Type: "TEXT", Options: "NOT NULL"},
		{Name: "created_at", Type: "INTEGER", Options: "NOT NULL"},
		{Name: "updated_at", Type: "INTEGER", Options: "NOT NULL"},
	},
	Extra: "",
}

// buildCreateTableSQL generates a CREATE TABLE statement from a schema
func (s *TableSchema) buildCreateTableSQL() string {
	query := "CREATE TABLE IF NOT EXISTS " + s.Name + " ("
	for i, field := range s.Fields {
		if i > 0 {
			query += ", "
		}
		query += field.Name + " " + field.Type
		if field.Options != "" {
			query += " " + field.Options
		}
	}
	if s.Extra != "" {
		query += ", " + s.Extra
	}
	query += ")"
	return query
}

// buildSchemaDefinition creates a canonical string representation of the schema
// Used to detect if schema has changed (like SLURM's definition column)
func (s *TableSchema) buildSchemaDefinition() string {
	def := ""
	for i, field := range s.Fields {
		if i > 0 {
			def += ","
		}
		def += field.Name + " " + field.Type + " " + field.Options
	}
	if s.Extra != "" {
		def += "," + s.Extra
	}
	return def
}
