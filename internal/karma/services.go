package karma

import "github.com/ntu-onemdp/onemdp-backend/internal/utils"

type KarmaService struct {
	repo *karmaRepository

	// Cache current settings for faster access
	currentSettings *Karma
}

var Service *KarmaService

func NewKarmaService(repo *karmaRepository) *KarmaService {
	settings, err := repo.getSettings()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error retrieving karma settings")
		return nil
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
