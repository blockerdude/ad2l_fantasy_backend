package model

type Roster struct {
	ID       int `json:"-"`
	SeasonID int `json:"-"`
	TeamID   int `json:"-"`
}

type RosterPlayer struct {
	RosterID int `json:"-"`
	PlayerID int `json:"-"`
	Position int `json:"position"`
}
