package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS accounts (
    puuid TEXT PRIMARY KEY NOT NULL,
    game_name TEXT NOT NULL,
    tag_line TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS match_ids (
    match_id TEXT PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS counts (
    name TEXT PRIMARY KEY NOT NULL,
    value INTEGER NOT NULL
);

INSERT OR IGNORE INTO counts (name, value) VALUES
    ('account_row', 0),
    ('match_id_row', 0)
`)
	return db, err
}
