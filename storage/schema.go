package storage

import (
	"github.com/jmoiron/sqlx"
)

var createGeneralQuery string = `create table if not exists General
(
	qword TEXT,
    query TEXT,
	entity TEXT,
    value TEXT,
    PRIMARY KEY(query)
)`

func initDBSchema(db *sqlx.DB) error {
	_, err := db.Exec(createGeneralQuery)
	return err
}