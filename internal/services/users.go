package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"slices"

	c "github.com/ntu-onemdp/onemdp-backend/config"
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
func (s *UserService) GetProfile(uid string) (*models.UserProfile, error) {
	return s.UsersRepo.GetUserProfile(uid)
}

// Get user profile photo
func (s *UserService) GetProfilePhoto(uid string) ([]byte, error) {
	return s.UsersRepo.GetProfilePhoto(uid)
}

// Check if user is pending registration
func (s *UserService) IsUserPending(email string) (bool, error) {
	return s.UsersRepo.IsUserPending(email)
}

// Admin: Get user information
func (s *UserService) GetUserAdmin(uid string) (*models.User, error) {
	return s.UsersRepo.GetUserByUidAdmin(uid)
}

// Admin: Get all users information
func (s *UserService) GetAllUsersAdmin() ([]models.User, error) {
	return s.UsersRepo.GetAllUsers()
}

// Get user role
func (s *UserService) GetRole(uid string) (models.UserRole, error) {
	role, err := s.UsersRepo.GetUserRole(uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting role by UID")
		return models.Unknown, err
	}

	return models.ParseRole(role)
}

// Get top 10 students with highest karma for given semester
func (s *UserService) GetTopKarma(semester string) ([]models.UserProfile, error) {
	TOP_N := 10 // Top 10 students

	return s.UsersRepo.GetTopKarma(semester, TOP_N)
}

// Check if user has staff permission
func (s *UserService) HasStaffPermission(uid string) (bool, error) {
	roleStr, err := s.UsersRepo.GetUserRole(uid)

	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting user role")
		return false, err
	}

	role, err := models.ParseRole(roleStr)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing user role")
		return false, err
	}

	return role >= models.Staff, nil
}

// Check if user has admin permission
func (s *UserService) HasAdminPermission(uid string) (bool, error) {
	roleStr, err := s.UsersRepo.GetUserRole(uid)

	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error getting user role")
		return false, err
	}

	role, err := models.ParseRole(roleStr)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing user role")
		return false, err
	}

	return role >= models.Admin, nil

}

// Update user's profile photo
// We do not need to validate whether original user is editing the profile photo as the UID is obtained from JWT.
func (s *UserService) UpdateProfilePhoto(uid string, file *multipart.FileHeader) error {
	// Reject if image size is too large
	if file.Size > c.MAX_PROFILE_IMG_SIZE {
		// Future improvement: automatically rescale image
		err := fmt.Sprintf("image size exceeds %d MB", c.MAX_PROFILE_IMG_SIZE/(1024*1024))
		utils.Logger.Error().Msgf("images size of %d exceeds %d MB", file.Size, c.MAX_PROFILE_IMG_SIZE/(1024*1024))
		return errors.New(err)
	}

	// Reject if image type is not supported (includes GIF)
	if !isValidType(file) || file.Header.Get("Content-Type") == "image/gif" {
		utils.Logger.Error().Msg("Unsupported image type")
		return errors.New("unsupported image type")
	}

	// Sanitize image
	image, err := utils.SanitizeImage(file)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to sanitize image")
		return err
	}

	utils.Logger.Debug().Msgf("Image parsing and sanitization complete for %s", uid)

	// Insert into db
	return s.UsersRepo.UpdateProfilePhoto(uid, image)
}

// Admin: update user's role
func (s *UserService) UpdateRole(uid string, role string) error {
	valid_roles := []string{"student", "staff", "admin"}
	if !slices.Contains(valid_roles, role) {
		return errors.New("invalid role provided")
	}

	return s.UsersRepo.UpdateUserRole(uid, role)
}
