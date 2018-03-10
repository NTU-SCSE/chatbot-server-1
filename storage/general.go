package storage

import (
	"../models"
)

// "database/sql"

// var addUserQuery string = "replace into User values(:user_name, :hashed_password)"
// var getUserQuery string = "select * from User where user_name = ?"
var listAllQuery string = "select * from General"

// var deleteUserQuery string = "delete from User where user_name = ?"

// func (db *dbImpl) CreateUser(user models.User) error {
//     _, err := db.sqliteDB.NamedExec(addUserQuery, user)
//     return err
// }

// func (db *dbImpl) GetUser(userName string) (*models.User, error) {
//     result := models.User{}
//     err := db.sqliteDB.Get(&result, getUserQuery, userName)

//     if err == sql.ErrNoRows {
//         return nil, nil
//     }

//     return &result, err
// }

func (db *dbImpl) ListAll() ([]models.General, error) {
	res := []models.General{}

	err := db.sqliteDB.Select(&res, listAllQuery)
	return res, err
}

// func (db *dbImpl) DeleteUser(userName string) error {
// 	_, err := db.sqliteDB.Exec(deleteUserQuery, userName)
// 	return err
// }
