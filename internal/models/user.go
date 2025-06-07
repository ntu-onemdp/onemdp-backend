package models

import "time"

type User struct {
	Uid          *string    `json:"uid" db:"uid"`
	Name         *string    `json:"name" db:"name"`
	Email        string     `json:"email" db:"email"`
	Role         string     `json:"role" db:"role"`
	DateCreated  *time.Time `json:"date_created" db:"date_created"`
	DateRemoved  *time.Time `json:"date_removed,omitempty" db:"date_removed"`
	Semester     string     `json:"semester" db:"semester"`
	ProfilePhoto *[]byte    `json:"profile_photo" db:"profile_photo"`
	Status       string     `json:"status" db:"status"`
	Karma        int        `json:"karma" db:"karma"`
}

// Initialize a new user for insertion into database
// Optional parameters:
// - role: If not provided, defaults to "student"
func CreateUser(email string, semester string, role string) *User {
	// Defaults to student if not provided
	if role == "" {
		role = "student"
	}

	return &User{
		Uid:          nil,
		Name:         nil,
		Email:        email,
		Role:         role,
		DateCreated:  nil,
		Semester:     semester,
		ProfilePhoto: nil,
		Status:       "active",
		Karma:        0,
	}
}

// Public user profile details returned by HandleGetUserProfile
type UserProfile struct {
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	ProfilePhoto *[]byte `json:"profile_photo" db:"profile_photo"`
	Semester     string  `json:"semester"`
	Karma        int     `json:"karma"`
	Role         string  `json:"role"`
}
