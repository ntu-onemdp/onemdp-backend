package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type CommentService struct {
	repo *repositories.CommentsRepository

	commentFactory *models.CommentFactory
}

var Comments *CommentService

// Create a new comment and insert into the repository.
// Returns comment id on success
func (s *CommentService) Create(authorUID string, articleUID string, content string) (string, error) {
	comment := s.commentFactory.New(authorUID, articleUID, content)

	if err := s.repo.Create(comment); err != nil {
		return "", err
	}

	return comment.CommentID, nil
}

// Delete comment if uid matches author
func (s *CommentService) Delete(commentID string, uid string) error {
	// Check if role is admin or staff
	hasStaffPermission, err := Users.HasStaffPermission(uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error checkecing for staff permission")
		return err
	}

	// Not staff/admin: check if user is author of the article
	if !hasStaffPermission {
		// Get uid of comment's author
		author, err := s.repo.GetAuthor(commentID)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting comment's author")
			return err
		}

		// uid of user does not match author
		if uid != author {
			utils.Logger.Warn().Msgf("User %s is not author of comment %s", uid, commentID)
			return utils.NewErrUnauthorized()
		}
	}

	return s.repo.Delete(commentID)
}
