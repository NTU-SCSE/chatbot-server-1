package models


type General struct {
    QWord string `json:"qword" db:"qword"`
    Query string `json:"query" db:"query"`
    Entity string `json:"entity" db:"entity"`
    Value string `json:"value" db:"value"`
}