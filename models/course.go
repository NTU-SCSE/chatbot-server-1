package models

type Course struct {
	Code        string `json:"code" db:"code"`
	Name        string `json:"name" db:"name"`
	AU          int    `json:"AU" db:"au"`
	PreReq      string `json:"preReq" db:"preReq"`
	Description string `json:"description" db:"description"`
}
