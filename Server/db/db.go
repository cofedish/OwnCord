// Package db provides database access for the OwnCord server.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/owncord/server/config"
	"github.com/owncord/server/migrations"
	_ "modernc.org/sqlite" // register the sqlite3 driver
)

// DB wraps *sql.DB and exposes the subset of methods needed by the server.
type DB struct {
	sqlDB   *sqlCompatDB
	dialect Dialect
}

// Open opens (or creates) a SQLite database at path. It remains for tests and
// transitional tooling; production startup should use OpenConfig.
func Open(path string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite db: %w", err)
	}

	// Verify the connection is actually usable.
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("pinging sqlite db: %w", err)
	}

	// SQLite only allows one writer at a time. Pin to a single connection
	// so concurrent goroutines queue on the Go side rather than getting
	// SQLITE_BUSY. For :memory: databases this also ensures all callers
	// share the same in-memory state.
	sqlDB.SetMaxOpenConns(1)

	// Enable WAL mode for better concurrent read performance.
	if _, err := sqlDB.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("enabling WAL mode: %w", err)
	}

	// Wait up to 5 seconds for the write lock instead of failing instantly.
	if _, err := sqlDB.Exec("PRAGMA busy_timeout=5000;"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("setting busy_timeout: %w", err)
	}

	// Enforce foreign key constraints.
	if _, err := sqlDB.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("enabling foreign keys: %w", err)
	}

	// Performance tuning (safe with WAL mode).
	if _, err := sqlDB.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("setting synchronous mode: %w", err)
	}
	if _, err := sqlDB.Exec("PRAGMA temp_store=MEMORY;"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("setting temp_store: %w", err)
	}
	if _, err := sqlDB.Exec("PRAGMA mmap_size=268435456;"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("setting mmap_size: %w", err)
	}
	if _, err := sqlDB.Exec("PRAGMA cache_size=-64000;"); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("setting cache_size: %w", err)
	}

	return &DB{
		sqlDB:   &sqlCompatDB{raw: sqlDB, dialect: DialectSQLite},
		dialect: DialectSQLite,
	}, nil
}

// OpenConfig opens the configured production database backend.
func OpenConfig(cfg config.DatabaseConfig) (*DB, error) {
	switch Dialect(cfg.Driver) {
	case DialectPostgres:
		return openPostgres(cfg)
	case "", DialectSQLite:
		return nil, fmt.Errorf("sqlite is no longer supported for server startup; configure database.driver=postgres")
	default:
		return nil, fmt.Errorf("unsupported database driver %q", cfg.Driver)
	}
}

func openPostgres(cfg config.DatabaseConfig) (*DB, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("postgres database URL is required")
	}

	pgCfg, err := pgx.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing postgres db URL: %w", err)
	}

	raw := stdlib.OpenDB(*pgCfg)
	raw.SetMaxOpenConns(cfg.MaxOpenConns)
	raw.SetMaxIdleConns(cfg.MaxIdleConns)
	if cfg.ConnMaxLifetimeMinutes > 0 {
		raw.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinutes) * time.Minute)
	}

	if err := raw.Ping(); err != nil {
		_ = raw.Close()
		return nil, fmt.Errorf("pinging postgres db: %w", err)
	}

	return &DB{
		sqlDB:   &sqlCompatDB{raw: raw, dialect: DialectPostgres},
		dialect: DialectPostgres,
	}, nil
}

// Migrate runs all SQL migration files from the embedded migrations FS in
// lexicographic order, applying each file exactly once.  It delegates to
// MigrateFS (defined in migrate.go) which maintains the schema_versions
// tracking table.
func Migrate(database *DB) error {
	switch database.dialect {
	case DialectPostgres:
		return MigrateFS(database, migrations.PostgresFS)
	default:
		return MigrateFS(database, migrations.SQLiteFS)
	}
}

// Close releases the underlying database connection.
func (d *DB) Close() error {
	if d.dialect == DialectSQLite {
		// Run PRAGMA optimize to analyze and update query planner statistics.
		_, _ = d.sqlDB.Exec("PRAGMA optimize;")
	}
	return d.sqlDB.Close()
}

// QueryRow executes a query that returns at most one row.
func (d *DB) QueryRow(query string, args ...any) *sql.Row {
	return d.sqlDB.QueryRow(query, args...)
}

// QueryRowContext executes a query that returns at most one row, with context.
func (d *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.sqlDB.QueryRowContext(ctx, query, args...)
}

// Exec executes a query that doesn't return rows.
func (d *DB) Exec(query string, args ...any) (sql.Result, error) {
	return d.sqlDB.Exec(query, args...)
}

// ExecContext executes a query that doesn't return rows, with context.
func (d *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.sqlDB.ExecContext(ctx, query, args...)
}

// Query executes a query that returns multiple rows.
func (d *DB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.sqlDB.Query(query, args...)
}

// QueryContext executes a query that returns multiple rows, with context.
func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.sqlDB.QueryContext(ctx, query, args...)
}

// Begin starts a database transaction.
func (d *DB) Begin() (*Tx, error) {
	return d.sqlDB.Begin()
}

// BeginTx starts a database transaction with context and options.
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	return d.sqlDB.BeginTx(ctx, opts)
}

// SQLDb returns the underlying *sql.DB for cases requiring direct access.
func (d *DB) SQLDb() *sql.DB {
	return d.sqlDB.raw
}

// Dialect reports the active SQL dialect.
func (d *DB) Dialect() Dialect {
	return d.dialect
}
