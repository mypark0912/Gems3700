package infra

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func OpenDB(configDir, dbName string) (*sql.DB, error) {
	dbPath := filepath.Join(configDir, dbName)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", dbName, err)
	}
	db.SetMaxOpenConns(1) // SQLite is single-writer
	return db, nil
}

func EnsureLogDB(configDir string) (*sql.DB, error) {
	dbPath := filepath.Join(configDir, "log.db")
	isNew := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		isNew = true
	}

	db, err := OpenDB(configDir, "log.db")
	if err != nil {
		return nil, err
	}

	if isNew {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			logdate TEXT,
			account TEXT,
			userRole TEXT,
			action TEXT
		)`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("create log table: %w", err)
		}
	}

	return db, nil
}

func EnsureUserDB(configDir string) (*sql.DB, error) {
	dbPath := filepath.Join(configDir, "user.db")
	isNew := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		isNew = true
	}

	db, err := OpenDB(configDir, "user.db")
	if err != nil {
		return nil, err
	}

	if isNew {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user (
			account TEXT PRIMARY KEY,
			username TEXT,
			password TEXT,
			email TEXT,
			role TEXT DEFAULT 'user',
			api TEXT DEFAULT ''
		)`)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("create user table: %w", err)
		}
	}

	return db, nil
}
