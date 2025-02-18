package models

import "time"

type User struct {
	Username        string     `json:"username" db:"username"`
	Name            string     `json:"name" db:"name"`
	DateCreated     time.Time  `json:"date_created" db:"date_created"`
	DateRemoved     *time.Time `json:"date_removed,omitempty" db:"date_removed"`
	Semester        int        `json:"semester" db:"semester"`
	PasswordChanged bool       `json:"password_changed" db:"password_changed"`
	ProfilePhoto    *string    `json:"profile_photo" db:"profile_photo"`
	Status          string     `json:"status" db:"status"`
}

// Initialize a new user for insertion into database
func CreateUser(username string, name string, semester int) *User {
	return &User{
		Username:        username,
		Name:            name,
		DateCreated:     time.Now(),
		Semester:        semester,
		PasswordChanged: false,
		ProfilePhoto:    nil,
		Status:          "active",
	}
}
