package model

import "time"

type Authn struct {
	ID           int       `json:"-"`
	SuperAdmin   bool      `json:"superAdmin"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"displayName"`
	LastAction   time.Time `json:"-"`
	SessionToken string    `json:"-"`
}
