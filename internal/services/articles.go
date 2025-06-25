package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
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
