package services

import (
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type AuthService struct {
	AuthRepo  *repositories.AuthRepository
	UsersRepo *repositories.UsersRepository
}

// Returns user information only if:
// 1. Username exists in both auth and user tables
// 2. Password is correct
// 3. Status is active
// If authenticated, this function returns (true, *user, role)
// else, this function returns (false, nil, "")
// This service is designed to minimize the number of DB read operations
func (s *AuthService) AuthenticateUser(username string, password string) (bool, *models.User, string) {
	// Query database for username and password
	auth, err := s.AuthRepo.GetAuthByUsername(username)
	if err != nil {
		utils.Logger.Error().Err(err)
		return false, nil, ""
	}

	// Compare user's plaintext pw and stored hashed pw
	match, err := argon2id.ComparePasswordAndHash(password, auth.Password)
	if !match || err != nil {
		return false, nil, ""
	}

	utils.Logger.Debug().Msg(fmt.Sprintf("%t", match))

	// Query for user's information if user is active
	user, err := s.UsersRepo.GetUserByUsername(username)
	if err != nil {
		utils.Logger.Err(err)
		// return false, nil, ""
	}
	if user == nil {
		utils.Logger.Error().Msg("User is nil after GetUserByUsername")
		return false, nil, ""
	}
	utils.Logger.Debug().Msg(fmt.Sprintf("%s, %s", user.Name, auth.Role))
	return true, user, auth.Role
}

// Update user's role. Admin can promote another user to admin.
// Service can be used for both promote and demote.
// This service checks that the role given is valid.
//
// This service does not check if the new role is the same as the old role.
func (s *AuthService) UpdateUserRole(username string, new_role string) error {
	// Convert new role to lowercase
	new_role = strings.ToLower(new_role)

	// Check if new role is valid
	if !isValidRole(new_role) {
		return fmt.Errorf("invalid role: %s", new_role)
	}

	return s.AuthRepo.UpdateUserRole(username, new_role)
}

// Check if user role is valid
func isValidRole(role string) bool {
	validRoles := map[string]bool{
		"student": true,
		"staff":   true,
		"admin":   true,
	}

	return validRoles[role]
}
