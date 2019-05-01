package main

import (
	"database/sql"
	"os"
	"regexp"
	"strings"
)

// A DB encapsulates a database connection and provides a simple ORM.
type DB struct {
	DSN   string
	db    *sql.DB
	stmts []string
}

// Podcasts gets Podcasts from the database
func (db *DB) Podcasts() (results []*Podcast) {
	return
}

// FindPodcast in the database which matches the (non-zero) values in the
// given `where` Podcast.
func (db *DB) FindPodcast(where *Podcast) (result *Podcast, err error) {
	result = &Podcast{}
	return
}

// Open the DB. Subsequent calls will replace the underlying connection.
func (db *DB) Open() (err error) {
	db.db, err = sql.Open(db.ConnectionConfig())
	return
}

// ConnectionConfig gets the driverName and dataSourceName of the DB.
func (db *DB) ConnectionConfig() (driverName, dataSourceName string) {
	var found bool
	driverName = "sqlite3"
	dataSourceName = db.DSN
	if dataSourceName == "" {
		if dataSourceName, found = os.LookupEnv("DATABASE"); !found {
			dataSourceName = ":memory:"
			return
		}
	}
	if regexp.MustCompile("^sqlite3?:").MatchString(dataSourceName) {
		dataSourceName = strings.SplitN(dataSourceName, ":", 2)[1]
		return
	}
	if regexp.MustCompile("^postgres://").MatchString(dataSourceName) {
		driverName = "postgres"
	}
	return
}
