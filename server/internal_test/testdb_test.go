package internal_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/ijry/sanjie/internal/db"
)

func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	if _, err := database.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}
	if err := db.Seed(database); err != nil {
		t.Fatalf("seed test db: %v", err)
	}
	return database
}
