package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type UserService struct {
	UsersRepo *repositories.UsersRepository
}

var Users *UserService

// Create new user and insert into the repository
func (s *UserService) CreateNewUser(email string, semester string, role string) error {
	user := models.CreatePendingUser(email, semester, role)

	err := s.UsersRepo.InsertOneUser(user)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("User not inserted into database")
		return err
	}

	return nil
}

// Register user by moving them from pending_users to users table
func (s *UserService) RegisterUser(uid string, email string, name string) error {
	return s.UsersRepo.RegisterUser(uid, email, name)
}

// Get user profile
func (s *UserService) GetProfile(email string) (*models.UserProfile, error) {
	return s.UsersRepo.GetUserProfile(email)
}

// Check if user is pending registration
func (s *UserService) IsUserPending(email string) (bool, error) {
	return s.UsersRepo.IsUserPending(email)
}

// Admin: Get user information
func (s *UserService) GetUserInformation(username string) (*models.User, error) {
	return s.UsersRepo.GetUserByUsernameAdmin(username)
}

// Admin: Get all users information
func (s *UserService) GetAllUsersInformation() ([]models.User, error) {
	return s.UsersRepo.GetAllUsers()
}
