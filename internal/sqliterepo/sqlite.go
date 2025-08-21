package sqliterepo

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewSql(fileName string) (*sql.DB, error) {
	// change using sqlite
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
