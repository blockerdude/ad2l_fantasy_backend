package model

type Conference struct {
	ID          int    `json:"-"`
	ObjectID    string `json:"objectID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
