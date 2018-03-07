package storage
// "github.com/satori/go.uuid"
import (
    "github.com/jmoiron/sqlx"
    "../models"    
    _ "github.com/mattn/go-sqlite3"
    
)

// TODO: Use singleton for db
type DB interface {
    // CreateUser(user models.User) error
    // GetUser(userName string) (*models.User, error)
    ListAll() ([]models.General, error)
    ListAllCourses() ([]models.Course, error)
    ListRecordsByIntent(tableName string) ([]models.Record, error)
    AddCourse(course *models.Course) (error)
    GetCourseByName(name string) (*models.Course, error)
    GetCourseByCode(code string) (*models.Course, error)
    PopulateCoursesData()
    // DeleteUser(userName string) error
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