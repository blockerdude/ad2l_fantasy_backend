package model

type LeagueParticipant struct {
	ID          int  `json:"-"`
	LeagueID    int  `json:"-"`
	AuthnID     int  `json:"-"`
	LeagueAdmin bool `json:"leagueAdmin"`
	Paid        bool `json:"paid"`
}
