package models

type User struct {
	Username        string `json:"username"`
	Name            string `json:"name"`
	DateCreated     string `json:"date_created"`
	DateRemoved     string `json:"date_removed"`
	Semester        int    `json:"semester"`
	PasswordChanged bool   `json:"password_changed"`
	ProfilePhoto    string `json:"profile_photo"`
	Status          string `json:"status"`
}

// Initialize a new user for insertion into database
func CreateUser(username string, name string, semester int) *User {
	return &User{
		Username:        username,
		Name:            name,
		DateCreated:     "",
		DateRemoved:     "",
		Semester:        semester,
		PasswordChanged: false,
		ProfilePhoto:    "",
		Status:          "",
	}
}
