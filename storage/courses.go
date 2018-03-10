package storage

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"strings"

	"../models"
)

var listAllCourses string = "select * from Courses"
var addCourseQuery string = "INSERT INTO Courses (code, name, au, preReq, description) VALUES(:code, :name, :au, :preReq, :description)"
var getCourseByCodeQuery string = "Select * from Courses where code = ?"
var getCourseByNameQuery string = "Select * from Courses where name = ?"

func (db *dbImpl) PopulateCoursesData() {
	var courses []models.Course
	file, _ := ioutil.ReadFile("./ce.json")
	json.Unmarshal(file, &courses)
	for _, elem := range courses {
		course := models.Course{strings.ToLower(elem.Code), strings.ToLower(elem.Name), elem.AU, strings.ToLower(elem.PreReq), strings.ToLower(elem.Description)}
		db.AddCourse(&course)
	}
}

func (db *dbImpl) ListAllCourses() ([]models.Course, error) {
	res := []models.Course{}

	err := db.sqliteDB.Select(&res, listAllCourses)
	return res, err
}

func (db *dbImpl) AddCourse(course *models.Course) error {
	_, _ = db.sqliteDB.NamedExec(addCourseQuery, course)
	return nil
}

func (db *dbImpl) GetCourseByCode(code string) (*models.Course, error) {
	result := models.Course{}
	err := db.sqliteDB.Get(&result, getCourseByCodeQuery, code)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (db *dbImpl) GetCourseByName(name string) (*models.Course, error) {
	result := models.Course{}
	err := db.sqliteDB.Get(&result, getCourseByNameQuery, name)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}
