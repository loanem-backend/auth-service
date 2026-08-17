package entity

import "time"

type Assistant struct {
	ID           int
	Name         string
	Phone        string
	HashPassword string
	Email        string
	Active       bool
	Period       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
