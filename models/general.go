package models

type General struct {
	QWord  string `json:"qword" db:"qword"`
	Query  string `json:"query" db:"query"`
	Intent string `json:"intent" db:"intent"`
	Value  string `json:"value" db:"value"`
}
