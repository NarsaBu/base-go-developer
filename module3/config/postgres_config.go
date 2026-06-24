package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection(cfg *Config) (*sql.DB, error) {
	dbCfg := cfg.Database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Server, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.DatabaseName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("Error while connecting to Postgres Database: ", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS urls(
		id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		url TEXT NOT NULL,
		alias TEXT NOT NULL UNIQUE);
	`)
	if err != nil {
		return nil, fmt.Errorf("Error while preparing statement: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("Error while executing statement: %w", err)
	}

	return db, nil
}
