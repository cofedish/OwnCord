// Package migrations holds embedded SQL migration files for the OwnCord server.
package migrations

import (
	"embed"
	"io/fs"
)

// SQLiteFS holds the legacy SQLite migration set used by tests and
// compatibility tooling.
//
//go:embed *.sql
var SQLiteFS embed.FS

// postgresFS holds the PostgreSQL migration set.
//
//go:embed postgres/*.sql
var postgresFS embed.FS

// PostgresFS exposes the embedded PostgreSQL migration files at the directory
// root expected by the migration runner.
var PostgresFS fs.FS = mustSub(postgresFS, "postgres")

func mustSub(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
