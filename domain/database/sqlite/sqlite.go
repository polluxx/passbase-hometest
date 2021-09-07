package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"passbase-hometest/config"
)

type SQLite struct {
	connection *sql.DB
}

var logger = zap.S().Named("SQLite")

func New(conf config.Database) *SQLite {
	db, err := sql.Open("sqlite3", conf.Name)
	if err != nil {
		logger.Fatalf("DB layer connection: %s", err)
	}

	return &SQLite{
		connection: db,
	}
}