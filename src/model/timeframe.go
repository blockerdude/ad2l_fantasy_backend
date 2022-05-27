package model

import "time"

type Timeframe struct {
	ID        int       `json:"-"`
	SeasonID  int       `json:"-"`
	Name      string    `json:"name"`
	Open      bool      `json:"open"`
	CloseDate time.Time `json:"startDate"`
}
