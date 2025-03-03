package services

import (
	"time"

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
		DateCreated:     time.Now(),
		Semester:        semester,
		PasswordChanged: false,
		ProfilePhoto:    nil,
		Status:          "active",
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

// Get user profile
func (s *UserService) GetUserProfile(username string) (*models.UserProfile, error) {
	return s.UsersRepo.GetUserProfile(username)
}

// Admin: Get user information
func (s *UserService) GetUserInformation(username string) (*models.User, error) {
	return s.UsersRepo.GetUserByUsernameAdmin(username)
}

// Admin: Get all users information
func (s *UserService) GetAllUsersInformation() ([]models.User, error) {
	return s.UsersRepo.GetAllUsers()
}
