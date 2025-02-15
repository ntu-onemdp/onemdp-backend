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

// For dev use only
// Initialize an empty user for testing
func CreateUser(username string, name string) *User {
	return &User{
		Username:        username,
		Name:            name,
		DateCreated:     "",
		DateRemoved:     "",
		Semester:        1,
		PasswordChanged: false,
		ProfilePhoto:    "",
		Status:          "",
	}
}
