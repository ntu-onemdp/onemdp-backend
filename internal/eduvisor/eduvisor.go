package eduvisor

import (
	"os"

	"github.com/google/uuid"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/services"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type EduvisorService struct {
	EduvisorModel *models.User
}

var Eduvisor *EduvisorService

func NewEduvisorService() *EduvisorService {
	// Retrieve eduvisor user object
	eduvisor := repositories.Users.GetEduvisor()

	// Eduvisor does not exist: create user and add
	if eduvisor == nil {
		utils.Logger.Trace().Msg("eduvisor not found in system, creating from scratch. ")

		uid, err := uuid.NewRandom()
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error generating UUID")
			utils.Logger.Warn().Msg("Eduvisor is not initialized")
			return nil
		}

		// Create eduvisor user
		eduvisor = models.CreateUser(uid.String(), constants.EDUVISOR_NAME, constants.EDUVISOR_EMAIL, "N.A.", "Bot")

		// Insert into users table directly. (Do not use services as this is not a normal student user.)
		if err = repositories.Users.DangerouslyInsertUser(eduvisor); err != nil {
			utils.Logger.Error().Err(err).Msg("Error inserting eduvisor into users")
			utils.Logger.Warn().Msg("Eduvisor is not initialized")
			return nil
		}
	}

	// Generate JWT token for Eduvisor
	claim := models.NewClaim(eduvisor.Uid)
	jwt, err := services.JwtHandler.GenerateJwt(claim)
	if err != nil {
		utils.Logger.Warn().Msg("Error generating jwt for eduvisor")
	}

	// Save jwt token to /eduvisor
	err = os.WriteFile("config/eduvisor-jwt.txt", []byte(jwt), 0644)
	if err != nil {
		utils.Logger.Warn().Msg("Error writing jwt key to file")
	}

	return &EduvisorService{
		EduvisorModel: eduvisor,
	}
}
