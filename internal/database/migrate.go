package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

// Migrations content is compiled into the binary so it can run at startup
// without a separate migration tooling step in production.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

// runMigrations applies pending SQL migrations using goose.
// A Postgres session lock prevents concurrent instances from migrating at once.
func runMigrations(sqlDB *sql.DB) error {
	migrationsDir, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("open embedded migrations: %w", err)
	}

	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return fmt.Errorf("create session locker: %w", err)
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		sqlDB,
		migrationsDir,
		goose.WithSessionLocker(locker),
	)
	if err != nil {
		return fmt.Errorf("create goose provider: %w", err)
	}

	applied, err := provider.Up(context.Background())
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	if len(applied) > 0 {
		log.Printf("✅ Migrations applied: %d", len(applied))
	} else {
		log.Println("🛢️ No pending migrations")
	}
	return nil
}
