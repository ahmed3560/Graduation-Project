package migrations

import (
	"io/fs"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var fileName = regexp.MustCompile(`^(\d{6})_([a-z0-9_]+)\.(up|down)\.sql$`)

// These tests need no database. They catch the mistakes that are easy to make by
// hand: a missing down file, a skipped number, a typo in the name.
func TestMigrationFiles(t *testing.T) {
	entries, err := fs.ReadDir(FS, ".")
	if err != nil {
		t.Fatal(err)
	}

	type pair struct{ up, down bool }
	versions := map[int]*pair{}

	for _, e := range entries {
		m := fileName.FindStringSubmatch(e.Name())
		if m == nil {
			t.Errorf("%s does not match NNNNNN_name.up.sql or NNNNNN_name.down.sql", e.Name())
			continue
		}
		v, _ := strconv.Atoi(m[1])
		if versions[v] == nil {
			versions[v] = &pair{}
		}
		if m[3] == "up" {
			versions[v].up = true
		} else {
			versions[v].down = true
		}

		body, err := fs.ReadFile(FS, e.Name())
		if err != nil {
			t.Fatal(err)
		}
		if strings.TrimSpace(string(body)) == "" {
			t.Errorf("%s is empty", e.Name())
		}
	}

	if len(versions) == 0 {
		t.Fatal("no migration files found")
	}
	for v := 1; v <= len(versions); v++ {
		p := versions[v]
		switch {
		case p == nil:
			t.Errorf("version %06d is missing: versions must be consecutive from 000001", v)
		case !p.up:
			t.Errorf("version %06d has no .up.sql file", v)
		case !p.down:
			t.Errorf("version %06d has no .down.sql file", v)
		}
	}
}
