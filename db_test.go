package main

import (
	"testing"
)

func TestConnectionConfig(t *testing.T) {
	for _, tc := range []struct{ in, dn, dsn string }{
		{"", "sqlite3", ":memory:"},
		{"sqlite3:foo.db", "sqlite3", "foo.db"},
		{":memory:", "sqlite3", ":memory:"},
		{"a/b/foo.db", "sqlite3", "a/b/foo.db"},
		{"postgres://u:p@h/db?x=y", "postgres", "postgres://u:p@h/db?x=y"},
		{"blah://u:p@h/db?x=y", "sqlite3", "blah://u:p@h/db?x=y"},
	} {
		t.Run(tc.in, func(t *testing.T) {
			db := &DB{DSN: tc.in}
			dn, dsn := db.ConnectionConfig()
			if dn != tc.dn {
				t.Errorf("expected driverName to be %q, got %q", tc.dn, dn)
			}
			if dsn != tc.dsn {
				t.Errorf("expected dataSourceName to be %q, got %q", tc.dsn, dsn)
			}
		})
	}
}

func TestPodcasts(t *testing.T) {
}
