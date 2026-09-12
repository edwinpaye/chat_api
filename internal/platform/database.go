package platform

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// OpenDB opens a *sql.DB using the DSN scheme.
// Expects a file:... DSN (modernc.org/sqlite driver).
func OpenDB(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite", dsn)
}
