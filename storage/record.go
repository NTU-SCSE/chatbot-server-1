package storage


import (
	"../models"
	"database/sql"
)
var listRecordsByIntentQuery string = "select * from "

func (db *dbImpl) ListRecordsByIntent(tableName string) ([]models.Record, error) {
	result := []models.Record{}
	// TODO: Proper handling here
	if tableName != "scse" && tableName != "scholarship" && tableName != "hostel" {
		return result, nil
	}
    err := db.sqliteDB.Select(&result, listRecordsByIntentQuery + tableName)
    if err == sql.ErrNoRows {
        return nil, nil
    }

    return result, err
}