package model

type League struct {
	ID          int    `json:"-"`
	ObjectID    string `json:"objectID"`
	SeasonID    int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
