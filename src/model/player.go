package model

type Player struct {
	ID      int    `json:"-"`
	SteamID string `json:"-"`
	Name    string `json:"name"`
	// In the future mmr, position, and other attributes might be added
}
