package storage

import (
	"github.com/jmoiron/sqlx"
)

var createGeneralQuery string = `create table if not exists General
(
	qword TEXT,
    query TEXT,
	intent TEXT,
    value TEXT,
    PRIMARY KEY(query)
)`

var createCoursesQuery string = `create table if not exists Courses
(
	code TEXT,
	name TEXT,
	au INTEGER,
	preReq TEXT,
	description TEXT
)`



func initDBSchema(db *sqlx.DB) error {
	_, err := db.Exec(createGeneralQuery)
	_, err = db.Exec(createCoursesQuery)
	return err
}