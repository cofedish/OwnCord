package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Dialect identifies the SQL dialect used by the current database backend.
type Dialect string

const (
	DialectSQLite   Dialect = "sqlite"
	DialectPostgres Dialect = "postgres"
)

var (
	insertTargetRE      = regexp.MustCompile(`(?is)^\s*INSERT\s+INTO\s+("?[\w.]+"?)`)
	insertOrIgnoreRE    = regexp.MustCompile(`(?is)^\s*INSERT\s+OR\s+IGNORE\s+INTO\s+`)
	nocaseCompareRE     = regexp.MustCompile(`(?i)=\s*\?\s+COLLATE\s+NOCASE`)
	strftimeColumnRE    = regexp.MustCompile(`(?i)strftime\('%s',\s*([a-z_][a-z0-9_]*)\s*\)`)
	serialInsertTables  = map[string]struct{}{
		"users":      {},
		"sessions":   {},
		"channels":   {},
		"messages":   {},
		"invites":    {},
		"roles":      {},
		"audit_log":  {},
		"reactions":  {},
		"emoji":      {},
		"sounds":     {},
	}
)

const postgresNowExpr = `to_char(timezone('utc', now()), 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`

type sqlCompatDB struct {
	raw     *sql.DB
	dialect Dialect
}

type Tx struct {
	raw     *sql.Tx
	dialect Dialect
}

type compatResult struct {
	lastInsertID int64
	rowsAffected int64
	hasLastID    bool
}

func (r compatResult) LastInsertId() (int64, error) {
	if !r.hasLastID {
		return 0, errors.New("LastInsertId is not supported for this statement")
	}
	return r.lastInsertID, nil
}

func (r compatResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func (d *sqlCompatDB) Ping() error {
	return d.raw.Ping()
}

func (d *sqlCompatDB) Close() error {
	return d.raw.Close()
}

func (d *sqlCompatDB) QueryRow(query string, args ...any) *sql.Row {
	return d.raw.QueryRow(normalizeQuery(d.dialect, query), args...)
}

func (d *sqlCompatDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.raw.QueryRowContext(ctx, normalizeQuery(d.dialect, query), args...)
}

func (d *sqlCompatDB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.raw.Query(normalizeQuery(d.dialect, query), args...)
}

func (d *sqlCompatDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.raw.QueryContext(ctx, normalizeQuery(d.dialect, query), args...)
}

func (d *sqlCompatDB) Exec(query string, args ...any) (sql.Result, error) {
	return execCompat(d.raw, d.dialect, query, args...)
}

func (d *sqlCompatDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return execCompatContext(ctx, d.raw, d.dialect, query, args...)
}

func (d *sqlCompatDB) Begin() (*Tx, error) {
	tx, err := d.raw.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{raw: tx, dialect: d.dialect}, nil
}

func (d *sqlCompatDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := d.raw.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{raw: tx, dialect: d.dialect}, nil
}

func (tx *Tx) QueryRow(query string, args ...any) *sql.Row {
	return tx.raw.QueryRow(normalizeQuery(tx.dialect, query), args...)
}

func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return tx.raw.QueryRowContext(ctx, normalizeQuery(tx.dialect, query), args...)
}

func (tx *Tx) Query(query string, args ...any) (*sql.Rows, error) {
	return tx.raw.Query(normalizeQuery(tx.dialect, query), args...)
}

func (tx *Tx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.raw.QueryContext(ctx, normalizeQuery(tx.dialect, query), args...)
}

func (tx *Tx) Exec(query string, args ...any) (sql.Result, error) {
	return execCompat(tx.raw, tx.dialect, query, args...)
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return execCompatContext(ctx, tx.raw, tx.dialect, query, args...)
}

func (tx *Tx) Commit() error {
	return tx.raw.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.raw.Rollback()
}

func execCompat(execable interface {
	Exec(string, ...any) (sql.Result, error)
	QueryRow(string, ...any) *sql.Row
}, dialect Dialect, query string, args ...any) (sql.Result, error) {
	normalized := normalizeQuery(dialect, query)
	if dialect == DialectPostgres {
		if table, ok := insertTargetTable(normalized); ok && shouldReturnInsertID(normalized, table) {
			var id int64
			if err := execable.QueryRow(normalized+" RETURNING id", args...).Scan(&id); err != nil {
				return nil, err
			}
			return compatResult{lastInsertID: id, rowsAffected: 1, hasLastID: true}, nil
		}
	}
	return execable.Exec(normalized, args...)
}

func execCompatContext(ctx context.Context, execable interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}, dialect Dialect, query string, args ...any) (sql.Result, error) {
	normalized := normalizeQuery(dialect, query)
	if dialect == DialectPostgres {
		if table, ok := insertTargetTable(normalized); ok && shouldReturnInsertID(normalized, table) {
			var id int64
			if err := execable.QueryRowContext(ctx, normalized+" RETURNING id", args...).Scan(&id); err != nil {
				return nil, err
			}
			return compatResult{lastInsertID: id, rowsAffected: 1, hasLastID: true}, nil
		}
	}
	return execable.ExecContext(ctx, normalized, args...)
}

func normalizeQuery(dialect Dialect, query string) string {
	if dialect != DialectPostgres {
		return query
	}

	normalized := query
	normalized = strings.ReplaceAll(normalized, "datetime('now')", postgresNowExpr)
	normalized = strings.ReplaceAll(normalized, "strftime('%s', 'now')", "EXTRACT(EPOCH FROM timezone('utc', now()))::bigint")
	normalized = strftimeColumnRE.ReplaceAllString(normalized, "EXTRACT(EPOCH FROM ($1)::timestamptz)::bigint")
	normalized = nocaseCompareRE.ReplaceAllString(normalized, "ILIKE ?")
	if insertOrIgnoreRE.MatchString(normalized) {
		normalized = insertOrIgnoreRE.ReplaceAllString(normalized, "INSERT INTO ")
		if !strings.Contains(strings.ToUpper(normalized), " ON CONFLICT ") && !strings.Contains(strings.ToUpper(normalized), " RETURNING ") {
			normalized += " ON CONFLICT DO NOTHING"
		}
	}
	return rebindPostgres(normalized)
}

func rebindPostgres(query string) string {
	var out strings.Builder
	out.Grow(len(query) + 8)
	argIndex := 1
	for _, ch := range query {
		if ch == '?' {
			out.WriteByte('$')
			out.WriteString(fmt.Sprintf("%d", argIndex))
			argIndex++
			continue
		}
		out.WriteRune(ch)
	}
	return out.String()
}

func insertTargetTable(query string) (string, bool) {
	matches := insertTargetRE.FindStringSubmatch(query)
	if len(matches) != 2 {
		return "", false
	}
	name := strings.Trim(matches[1], `"`)
	if dot := strings.LastIndex(name, "."); dot >= 0 {
		name = name[dot+1:]
	}
	name = strings.ToLower(name)
	_, ok := serialInsertTables[name]
	return name, ok
}

func shouldReturnInsertID(query, table string) bool {
	if _, ok := serialInsertTables[table]; !ok {
		return false
	}
	upper := strings.ToUpper(query)
	return !strings.Contains(upper, " ON CONFLICT ") && !strings.Contains(upper, " RETURNING ")
}

