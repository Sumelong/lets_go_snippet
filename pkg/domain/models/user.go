package models

import "time"

// User  type.
type User struct {
	ID             int
	Created        time.Time
	updated        time.Time
	Name           string
	Email          string
	HashedPassword []byte
	Active         bool
}

func (u User) GetID() uint {
	return uint(u.ID)
}
