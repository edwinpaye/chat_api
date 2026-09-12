package platform

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenDB opens a *sql.DB using the DSN scheme.
// Expects a postgres:// URL (pgx stdlib driver).

import (
    "context"
    "errors"
    "fmt"
    "time"
)

// verifies that the database is reachable before returning. It configures the pool with
func OpenDB(dsn string) (*sql.DB, error) {
    if dsn == "" {
        return nil, errors.New("postgres DSN is empty")
    }

    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, fmt.Errorf("open postgres: %w", err)
    }

    // Tune the pool for a dev environment.
    db.SetMaxOpenConns(10)           // limit simultaneous connections
    db.SetMaxIdleConns(5)            // keep a few around for reuse
    db.SetConnMaxLifetime(30 * time.Minute)
    db.SetConnMaxIdleTime(10 * time.Minute)

    // Ping with a short timeout so you fail fast if Postgres is down.
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        _ = db.Close() // clean up the partially opened pool
        return nil, fmt.Errorf("ping postgres: %w", err)
    }

    return db, nil
}
