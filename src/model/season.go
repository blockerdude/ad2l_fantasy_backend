package model

import "time"

type Season struct {
	ID           int       `json:"-"`
	ConferenceID int       `json:"-"`
	Name         string    `json:"name"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
}
