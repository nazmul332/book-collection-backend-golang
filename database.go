package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	if dataSourceName != ":memory:" {
		dataSourceName += "?_auth&_auth_user=admin&_auth_pass=admin123&_auth_crypt=sha1&parseTime=true"
	}
	if !strings.Contains(dataSourceName, "?") {
		dataSourceName += "?parseTime=true"
	} else {
		dataSourceName += "&parseTime=true"
	}

	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	err = createTables()
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	authorsTable := `
	CREATE TABLE IF NOT EXISTS authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		bio TEXT
	);`

	genresTable := `
	CREATE TABLE IF NOT EXISTS genres (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT
	);`

	booksTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author_id INTEGER,
		genre_id INTEGER,
		published_at DATETIME,
		publisher TEXT,
		isbn TEXT,
		FOREIGN KEY (author_id) REFERENCES authors(id),
		FOREIGN KEY (genre_id) REFERENCES genres(id)
	);`

	tables := []string{usersTable, authorsTable, genresTable, booksTable}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	return nil
}