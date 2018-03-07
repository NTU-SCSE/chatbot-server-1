package models


type Record struct {
    Params string `json:"params" db:"params"`
    Response string `json:"response" db:"response"`
}