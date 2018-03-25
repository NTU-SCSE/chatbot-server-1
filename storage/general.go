package storage

import (
	"../models"
)
var listAllQuery string = "select * from General"


func (db *dbImpl) ListAll() ([]models.General, error) {
	res := []models.General{}

	err := db.sqliteDB.Select(&res, listAllQuery)
	return res, err
}