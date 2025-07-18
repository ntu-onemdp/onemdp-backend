package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ArticleService struct {
	articleRepo  *repositories.ArticleRepository
	commentsRepo *repositories.CommentsRepository

	articleFactory *models.ArticleFactory
}

var Articles *ArticleService

func NewArticleService(articleRepo *repositories.ArticleRepository, commentsRepo *repositories.CommentsRepository) *ArticleService {
	return &ArticleService{
		articleRepo:    articleRepo,
		commentsRepo:   commentsRepo,
		articleFactory: models.NewArticleFactory(),
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

// Retrieve all articles in given page.
func (s *ArticleService) GetArticles(sort string, size int, desc bool, page int, uid string) ([]models.Article, error) {
	// Convert sort string to ThreadColumn object
	column := models.StrToThreadColumn(sort)

	// Retrieve articles from db
	return s.articleRepo.GetAll(uid, column, page, size, desc)
}

// Retrieve article metadata
func (s *ArticleService) GetMetadata() (*models.ArticlesMetadata, error) {
	return s.articleRepo.GetMetadata()
}

// Retrieve article and all related comments
func (s *ArticleService) GetArticle(articleID string, uid string) (*models.Article, []models.Comment, error) {
	// Retrieve article
	article, err := s.articleRepo.GetByID(articleID, uid)
	if err != nil {
		utils.Logger.Trace().Err(err).Msg("Error retrieving article")
		return nil, nil, err
	}

	// Retrive comments
	comments, err := s.commentsRepo.GetCommentsByArticleID(articleID, uid)
	if err != nil {
		utils.Logger.Trace().Err(err).Msg("Error retrieving comments")
		return article, nil, err
	}

	return article, comments, nil
}

// Check if article exists/ is available
func (s *ArticleService) ArticleExists(articleID string) bool {
	return s.articleRepo.IsAvailable(articleID)
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
