package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ArticleService struct {
	articleRepo *repositories.ArticleRepository

	articleFactory *models.ArticleFactory
}

var Articles *ArticleService

func NewArticleService(articleRepo *repositories.ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
	}
}

// Create new article and insert into the repository
// Returns article id on success
func (s *ArticleService) CreateNewArticle(author string, title string, content string) (string, error) {
	article := s.articleFactory.New(author, title, content)

	err := s.articleRepo.Insert(article)
	if err != nil {
		return "", err
	}

	return article.ArticleID, nil
}

// Retrieve article and all related comments
func (s *ArticleService) GetArticle(articleID string, uid string) (*models.Article, error) {
	return s.articleRepo.GetByID(articleID, uid)
}

// Delete article and all associated comments
func (s *ArticleService) DeleteArticle(articleID string, uid string) error {
	// Check if role is admin or staff
	hasStaffPermission, err := Users.HasStaffPermission(uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error checkecing for staff permission")
		return err
	}

	// Not staff/admin: check if user is author of the article
	if !hasStaffPermission {
		// Get uid of article author
		author, err := s.articleRepo.GetAuthor(articleID)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting article author")
			return err
		}
		if author != uid {
			utils.Logger.Warn().Msgf("User %s is not author of article %s", uid, articleID)
			return utils.NewErrUnauthorized()
		}
	}

	return s.articleRepo.Delete(articleID)
}
