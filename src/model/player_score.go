package model

type PlayerScore struct {
	TimeframeID            int      `json:"-"`
	PlayerID               int      `json:"-"`
	MatchID                string   `json:"-"`
	SubstituteID           *int     `json:"-"`
	Total                  *float64 `json:"total"` // TODO: think about changing this name to 'TotalOverride'
	Kills                  int      `json:"kills"`
	Deaths                 int      `json:"deaths"`
	LastHits               int      `json:"lastHits"`
	Denies                 int      `json:"denies"`
	TeamfightParticipation float64  `json:"teamfightParticipation"`
	GoldPerMinute          int      `json:"goldPerMinute"`
	TowerKills             int      `json:"towerKills"`
	RoshKills              int      `json:"roshKills"`
	ObserversPlaced        int      `json:"observersPlaced"`
	CampsStacked           int      `json:"campsStacked"`
	RunesTaken             int      `json:"runesTaken"`
	FirstBlood             bool     `json:"firstBlood"`
	StunSeconds            float64  `json:"stunSeconds"`
}
