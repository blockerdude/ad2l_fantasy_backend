package model

type Team struct {
	ID           int    `json:"-"`
	ConferenceID int    `json:"-"`
	Name         string `json:"name"`
	// TODO: Flag_Url is a field in the db, but likely leaving as empty for now
}
