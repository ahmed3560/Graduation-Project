package migrator

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"io/fs"
	"net/url"
	"os"
	"testing"
	"testing/fstest"

	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver for the test helpers
)

// Two tiny migrations that exist only for these tests.
var twoMigrations = fstest.MapFS{
	"000001_create_widgets.up.sql":   {Data: []byte("CREATE TABLE widgets (id int PRIMARY KEY);")},
	"000001_create_widgets.down.sql": {Data: []byte("DROP TABLE widgets;")},
	"000002_create_gadgets.up.sql":   {Data: []byte("CREATE TABLE gadgets (id int PRIMARY KEY);")},
	"000002_create_gadgets.down.sql": {Data: []byte("DROP TABLE gadgets;")},
}

func TestDriverURL(t *testing.T) {
	tests := []struct{ in, want string }{
		{"postgres://u:p@localhost:5432/db?sslmode=disable", "pgx5://u:p@localhost:5432/db?sslmode=disable"},
		{"postgresql://u:p@localhost/db", "pgx5://u:p@localhost/db"},
		{"pgx5://u:p@localhost/db", "pgx5://u:p@localhost/db"},
	}
	for _, tc := range tests {
		if got := driverURL(tc.in); got != tc.want {
			t.Errorf("driverURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// setup gives a test its own empty schema inside the test database, so tests never
// see each other's tables. It skips the test when TEST_DATABASE_URL is not set.
func setup(t *testing.T, files fs.FS) (*Migrator, *sql.DB, string) {
	t.Helper()
	base := os.Getenv("TEST_DATABASE_URL")
	if base == "" {
		t.Skip("TEST_DATABASE_URL is not set; skipping database tests")
	}

	db, err := sql.Open("pgx", base)
	if err != nil {
		t.Fatal(err)
	}
	suffix := make([]byte, 6)
	if _, err := rand.Read(suffix); err != nil {
		t.Fatal(err)
	}
	schema := "migrator_test_" + hex.EncodeToString(suffix)
	if _, err := db.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec("DROP SCHEMA " + schema + " CASCADE")
		_ = db.Close()
	})

	u, err := url.Parse(base)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema) // the migrator's tables land in this schema
	u.RawQuery = q.Encode()

	m, err := New(u.String(), files)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = m.Close() })
	return m, db, schema
}

func tableExists(t *testing.T, db *sql.DB, schema, name string) bool {
	t.Helper()
	var exists bool
	if err := db.QueryRow("SELECT to_regclass($1) IS NOT NULL", schema+"."+name).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	return exists
}

func mustVersion(t *testing.T, m *Migrator) (uint, bool) {
	t.Helper()
	v, dirty, err := m.Version()
	if err != nil {
		t.Fatalf("Version: %v", err)
	}
	return v, dirty
}

func TestVersion_EmptyDatabaseIsZero(t *testing.T) {
	m, _, _ := setup(t, twoMigrations)

	if v, dirty := mustVersion(t, m); v != 0 || dirty {
		t.Errorf("version = %d dirty = %v, want 0 false", v, dirty)
	}
}

func TestUp_AppliesEveryMigration(t *testing.T) {
	m, db, schema := setup(t, twoMigrations)

	if err := m.Up(); err != nil {
		t.Fatalf("Up: %v", err)
	}
	if v, dirty := mustVersion(t, m); v != 2 || dirty {
		t.Errorf("version = %d dirty = %v, want 2 false", v, dirty)
	}
	for _, table := range []string{"widgets", "gadgets"} {
		if !tableExists(t, db, schema, table) {
			t.Errorf("table %s was not created", table)
		}
	}
}

func TestUp_RunningTwiceIsHarmless(t *testing.T) {
	m, _, _ := setup(t, twoMigrations)

	if err := m.Up(); err != nil {
		t.Fatalf("first Up: %v", err)
	}
	if err := m.Up(); err != nil {
		t.Fatalf("second Up should do nothing, got: %v", err)
	}
	if v, _ := mustVersion(t, m); v != 2 {
		t.Errorf("version = %d, want 2", v)
	}
}

func TestDown_UndoesOnlyTheLastN(t *testing.T) {
	m, db, schema := setup(t, twoMigrations)
	if err := m.Up(); err != nil {
		t.Fatal(err)
	}

	if err := m.Down(1); err != nil {
		t.Fatalf("Down(1): %v", err)
	}
	if v, _ := mustVersion(t, m); v != 1 {
		t.Errorf("version = %d, want 1", v)
	}
	if tableExists(t, db, schema, "gadgets") {
		t.Error("gadgets should be gone")
	}
	if !tableExists(t, db, schema, "widgets") {
		t.Error("widgets should still exist")
	}

	if err := m.Down(1); err != nil {
		t.Fatalf("second Down(1): %v", err)
	}
	if v, _ := mustVersion(t, m); v != 0 {
		t.Errorf("version = %d, want 0", v)
	}
}

func TestDown_TooManyChangesNothing(t *testing.T) {
	m, db, schema := setup(t, twoMigrations)
	if err := m.Up(); err != nil {
		t.Fatal(err)
	}

	if err := m.Down(3); err == nil {
		t.Fatal("Down(3) with 2 applied should fail")
	}
	if v, _ := mustVersion(t, m); v != 2 {
		t.Errorf("version = %d, want 2 (nothing should have changed)", v)
	}
	if !tableExists(t, db, schema, "widgets") || !tableExists(t, db, schema, "gadgets") {
		t.Error("tables should be untouched")
	}
}

func TestDown_RejectsZeroAndNegative(t *testing.T) {
	m, _, _ := setup(t, twoMigrations)
	if err := m.Up(); err != nil {
		t.Fatal(err)
	}
	for _, n := range []int{0, -1} {
		if err := m.Down(n); err == nil {
			t.Errorf("Down(%d) should fail", n)
		}
	}
	if v, _ := mustVersion(t, m); v != 2 {
		t.Errorf("version = %d, want 2", v)
	}
}

func TestFailingMigration_LeavesDatabaseDirtyAndBlocksDown(t *testing.T) {
	files := fstest.MapFS{
		"000001_create_widgets.up.sql":   {Data: []byte("CREATE TABLE widgets (id int PRIMARY KEY);")},
		"000001_create_widgets.down.sql": {Data: []byte("DROP TABLE widgets;")},
		"000002_broken.up.sql":           {Data: []byte("CREATE TABLE gadgets (id int PRIMARY KEY); THIS IS NOT SQL;")},
		"000002_broken.down.sql":         {Data: []byte("DROP TABLE gadgets;")},
	}
	m, db, schema := setup(t, files)

	if err := m.Up(); err == nil {
		t.Fatal("Up should fail on the broken migration")
	}
	if v, dirty := mustVersion(t, m); v != 2 || !dirty {
		t.Errorf("version = %d dirty = %v, want 2 true", v, dirty)
	}
	if tableExists(t, db, schema, "gadgets") {
		t.Error("the broken file ran as one unit, so gadgets should have been rolled back")
	}
	if err := m.Down(1); err == nil {
		t.Error("Down must refuse to run while the database is dirty")
	}
}
