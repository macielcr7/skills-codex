package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	databaseURL := flag.String("database-url", os.Getenv("DATABASE_URL"), "database URL (ex.: postgres://user:pass@localhost:5432/db?sslmode=disable)")
	migrationsDir := flag.String("migrations", "migrations", "path to migrations dir (default: migrations)")
	direction := flag.String("direction", "up", "migration direction: up|down")
	steps := flag.Int("steps", 0, "number of steps (0 = all)")
	flag.Parse()

	if *databaseURL == "" {
		slog.Error("missing DATABASE_URL (or --database-url)")
		os.Exit(2)
	}

	sourceURL := fmt.Sprintf("file://%s", *migrationsDir)
	m, err := migrate.New(sourceURL, *databaseURL)
	if err != nil {
		slog.Error("failed to create migrator", "error", err)
		os.Exit(1)
	}

	defer func() {
		_, _ = m.Close()
	}()

	var runErr error
	switch *direction {
	case "up":
		if *steps == 0 {
			runErr = m.Up()
		} else {
			runErr = m.Steps(*steps)
		}
	case "down":
		if *steps == 0 {
			runErr = m.Down()
		} else {
			runErr = m.Steps(-*steps)
		}
	default:
		slog.Error("invalid --direction", "direction", *direction)
		os.Exit(2)
	}

	if runErr != nil && runErr != migrate.ErrNoChange {
		slog.Error("migration failed", "error", runErr)
		os.Exit(1)
	}

	slog.Info("migrations done", "direction", *direction, "steps", *steps)
}

