package model

type AuthnPick struct {
	LeagueParticipantID int   `json:"-"`
	TimeframeID         int   `json:"-"`
	PlayerIDs           []int `json:"-"`
}
