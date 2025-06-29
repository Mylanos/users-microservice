package models

import (
	"time"

	"github.com/google/uuid"
)

// BU representation for users
type User struct {
	ID    uuid.UUID
	Name  string
	Email string
	// convert to age maybe
	DateOfBirth time.Time
}

func NewUser(id uuid.UUID, name string, email string, dateOfBirth time.Time) *User {
	return &User{
		ID:          id,
		DateOfBirth: dateOfBirth,
		Name:        name,
		Email:       email,
	}
}
