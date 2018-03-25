package storage

import (
	"../models"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB interface {
	ListAll() ([]models.General, error)
	ListAllCourses() ([]models.Course, error)
	ListRecordsByIntent(tableName string) ([]models.Record, error)
	AddCourse(course *models.Course) error
	GetCourseByName(name string) (*models.Course, error)
	GetCourseByCode(code string) (*models.Course, error)
	PopulateCoursesData()
	DeleteCourseByCode(name string) error
}

type dbImpl struct {
	sqliteDB *sqlx.DB
}

func NewDB(fileName string) (DB, error) {
	db, err := sqlx.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}
	db.MustExec("PRAGMA foreign_keys = ON;")
	initDBSchema(db)
	return &dbImpl{db}, err
}
