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
