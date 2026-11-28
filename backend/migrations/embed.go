// Package migrations holds the SQL migration files and embeds them in the binary,
// so cmd/migrate needs no loose files next to it when it runs on a server.
//
// File naming rules (migrations_test.go checks them):
//
//	000001_short_name.up.sql     applies the change
//	000001_short_name.down.sql   undoes it
//
// Versions are consecutive: 000001, 000002, 000003, with no gaps.
package migrations

import "embed"

// FS contains every *.sql file in this folder.
//
//go:embed *.sql
var FS embed.FS
