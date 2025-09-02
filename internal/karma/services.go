package karma

import (
	"github.com/jackc/pgx/v5"
	"github.com/ntu-onemdp/onemdp-backend/internal/semester"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type KarmaService struct {
	repo *karmaRepository

	// Cache current settings for faster access
	currentSettings *Karma
}

var Service *KarmaService

func NewKarmaService(repo *karmaRepository) *KarmaService {
	settings, err := repo.getSettings()
	if err != nil {
		if err == pgx.ErrNoRows {
			// Karma settings have not been initialized for current semester yet.
			utils.Logger.Warn().Msg("Karma settings not found, initializing with default values")

			// Default values for karma
			settings = &Karma{
				Semester:         *semester.Service.GetCurrentSem(),
				CreateThreadPts:  10,
				CreateArticlePts: 20,
				CreateCommentPts: 2,
				CreatePostPts:    5,
				LikePts:          1,
			}

			if err := repo.insert(settings.Semester); err != nil {
				utils.Logger.Error().Err(err).Msg("Error inserting default karma settings for new semester")
				return nil
			}

			if err := repo.update(*settings); err != nil {
				utils.Logger.Error().Err(err).Msg("Error setting default karma settings for current semester")
				return nil
			}
		} else {
			// Other error
			utils.Logger.Error().Err(err).Msg("Error retrieving karma settings")
			return nil
		}
	}

	return &KarmaService{
		repo:            repo,
		currentSettings: settings,
	}
}

// Get the current settings
func (s *KarmaService) GetKarmaSettings() *Karma {
	return s.currentSettings
}

// Insert new karma settings into database.
// This function should be called when a new semester is created
func (s *KarmaService) Insert(semester string) error {
	return s.repo.insert(semester)
}

// Retrieve current karma settings from cache
func (s *KarmaService) GetSettings() *Karma {
	return s.currentSettings
}

// Update karma settings
func (s *KarmaService) UpdateSettings(settings Karma) error {
	if err := s.repo.update(settings); err != nil {
		return err
	}

	// Update cache
	s.currentSettings = &settings
	utils.Logger.Info().Interface("karma settings", settings).Msg("Karma settings successfully updated")
	return nil
}
