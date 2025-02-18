package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type UserService struct {
	UsersRepo *repositories.UsersRepository
}

// Create new user and insert into the repository
func (s *UserService) CreateNewUser(username string, name string, semester int) error {
	user := models.User{
		Username:        username,
		Name:            name,
		DateCreated:     "",
		DateRemoved:     "",
		Semester:        semester,
		PasswordChanged: false,
		ProfilePhoto:    "",
		Status:          "",
	}

	err := s.UsersRepo.InsertOneUser(&user)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("User not inserted into database")
		return err
	}

	return nil
}

// Check if user's password has been changed
func (s *UserService) HasPasswordChanged(username string) (bool, error) {
	return s.UsersRepo.GetUserPasswordChanged(username)
}
