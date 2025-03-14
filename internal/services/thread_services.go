package services

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ThreadService struct {
	threadRepo *repositories.ThreadRepository
	postRepo   *repositories.PostsRepository
	likesRepo  *repositories.LikesRepository

	threadFactory *models.ThreadFactory
	postFactory   *models.PostFactory
}

func NewThreadService(threadRepo *repositories.ThreadRepository, postRepo *repositories.PostsRepository, likesRepo *repositories.LikesRepository) *ThreadService {
	return &ThreadService{
		threadRepo:    threadRepo,
		postRepo:      postRepo,
		likesRepo:     likesRepo,
		threadFactory: models.NewThreadFactory(),
		postFactory:   models.NewPostFactory(),
	}
}

// Create new thread and insert into the repository
func (s *ThreadService) CreateNewThread(author string, title string, content string) error {
	thread := s.threadFactory.New(author, title, content)

	err := s.threadRepo.Create(thread)
	if err != nil {
		return err
	}

	postFactory := models.PostFactory{}
	post := postFactory.New(thread.Author, thread.ThreadID, thread.Title, thread.Preview, nil, true)

	err = s.postRepo.Create(post)
	return err
}

// Retrieve thread and all associated posts
func (s *ThreadService) GetThread(threadID string) (*models.Thread, []models.Post, error) {
	// Retrieve thread from db
	thread, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		utils.Logger.Trace().Msg("Error getting thread")
		return nil, nil, err
	}

	// Get number of likes for thread
	thread.NumLikes, err = s.likesRepo.GetNumLikes(threadID)
	if err != nil {
		utils.Logger.Trace().Msg("Error getting number of likes")
		return nil, nil, err
	}

	// Retrieve posts from db
	posts, err := s.postRepo.GetPostByThreadId(threadID)
	if err != nil {
		utils.Logger.Trace().Msg("Error getting posts from db")
		return nil, nil, err
	}

	// Get number of likes for each post
	for i := range posts {
		posts[i].NumLikes, err = s.likesRepo.GetNumLikes(posts[i].PostID)
		if err != nil {
			utils.Logger.Trace().Msg("Error getting number of likes")
			return nil, nil, err
		}
	}

	return thread, posts, nil
}

// Check if thread exists
func (s *ThreadService) ThreadExists(threadID string) bool {
	return s.threadRepo.IsAvailable(threadID)
}

// Update thread's last activity
func (s *ThreadService) UpdateThreadLastActivity(threadID string) error {
	return s.threadRepo.UpdateActivity(threadID)
}

// Update thread's title and preview
func (s *ThreadService) UpdateThread(threadID string, title string, content string, claim *utils.JwtClaim) error {
	// Check if role is admin or staff
	if !HasStaffPermission(claim) {
		author, err := s.threadRepo.GetAuthor(threadID)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of thread")
			return err
		}

		// Check if author of thread matches the author in JWT claim
		if author != claim.Username {
			return utils.NewErrUnauthorized()
		}
	}

	return s.threadRepo.Update(threadID, title, getPreview(content))
}

// Delete thread and all associated posts
func (s *ThreadService) DeleteThread(threadID string, claim *utils.JwtClaim) error {
	// Check if role is admin or staff
	if !HasStaffPermission(claim) {
		author, err := s.threadRepo.GetAuthor(threadID)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of thread")
			return err
		}

		// Check if author of thread matches the author in JWT claim
		if author != claim.Username {
			return utils.NewErrUnauthorized()
		}
	}

	err := s.threadRepo.Delete(threadID)
	if err != nil {
		utils.Logger.Trace().Msg("Error deleting thread")
		return err
	}

	return s.postRepo.DeletePostsByThread(threadID)
}

// Utility function to get preview from content
func getPreview(content string) string {
	const MAX_PREVIEW_LENGTH = 100

	if len(content) <= MAX_PREVIEW_LENGTH {
		return content
	}
	return content[:MAX_PREVIEW_LENGTH]
}
