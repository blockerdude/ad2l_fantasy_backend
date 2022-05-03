package model

type Conference struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
