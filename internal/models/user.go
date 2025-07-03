package models

import (
	"fmt"
	"time"
)

type User struct {
	Uid          string     `json:"uid" db:"uid"`
	Name         string     `json:"name" db:"name"`
	Email        string     `json:"email" db:"email"`
	Role         string     `json:"role" db:"role"`
	DateCreated  time.Time  `json:"date_created" db:"date_created"`
	DateRemoved  *time.Time `json:"date_removed,omitempty" db:"date_removed"`
	Semester     string     `json:"semester" db:"semester"`
	ProfilePhoto *[]byte    `json:"profile_photo" db:"profile_photo"`
	Status       string     `json:"status" db:"status"`
	Karma        int        `json:"karma" db:"karma"`
}

// User pending registration
type PendingUser struct {
	Email       string    `json:"email" db:"email"`
	Role        string    `json:"role" db:"role"`
	Semester    string    `json:"semester" db:"semester"`
	TimeCreated time.Time `json:"time_created" db:"time_created"`
}

// Initialize a new user for insertion into user table after registration
func CreateUser(uid string, name string, email string, semester string, role string) *User {
	return &User{
		Uid:          uid,
		Name:         name,
		Email:        email,
		Role:         role,
		DateCreated:  time.Now(),
		Semester:     semester,
		ProfilePhoto: nil,
		Status:       "active",
		Karma:        0,
	}
}

// Create a pending user for registration
func CreatePendingUser(email string, semester string, role string) *PendingUser {
	// Defaults to student if not provided
	if role == "" {
		role = "student"
	}

	return &PendingUser{
		Email:       email,
		Role:        role,
		Semester:    semester,
		TimeCreated: time.Now(),
	}
}

// Public user profile details returned by HandleGetUserProfile
type UserProfile struct {
	Uid          string  `json:"uid"`
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	ProfilePhoto *[]byte `json:"profile_photo" db:"profile_photo"`
	Semester     string  `json:"semester"`
	Karma        int     `json:"karma"`
	Role         string  `json:"role"`
}

// Role levels
type UserRole int

const (
	Unknown          = 0
	Student UserRole = iota
	Bot
	Staff
	Admin
)

func ParseRole(role string) (UserRole, error) {
	switch role {
	case "student":
		return Student, nil
	case "staff":
		return Staff, nil
	case "admin":
		return Admin, nil
	default:
		return Student, fmt.Errorf("unknown role: %s", role)
	}
}
