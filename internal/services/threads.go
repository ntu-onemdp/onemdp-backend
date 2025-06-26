package services

import (
	"time"

	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/repositories"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type ThreadService struct {
	threadRepo *repositories.ThreadsRepository
	postRepo   *repositories.PostsRepository
	likesRepo  *repositories.LikesRepository

	threadFactory *models.ThreadFactory
	postFactory   *models.PostFactory
}

var Threads *ThreadService

func NewThreadService(threadRepo *repositories.ThreadsRepository, postRepo *repositories.PostsRepository, likesRepo *repositories.LikesRepository) *ThreadService {
	return &ThreadService{
		threadRepo:    threadRepo,
		postRepo:      postRepo,
		likesRepo:     likesRepo,
		threadFactory: models.NewThreadFactory(),
		postFactory:   models.NewPostFactory(),
	}
}

// Create new thread and insert into the repository
// Returns thread id on success
func (s *ThreadService) CreateNewThread(author string, title string, content string) (string, error) {
	thread := s.threadFactory.New(author, title, content)

	err := s.threadRepo.Insert(thread)
	if err != nil {
		return "", err
	}

	postFactory := models.PostFactory{}
	post := postFactory.New(thread.AuthorUid, thread.ThreadID, thread.Title, content, nil, true)

	err = s.postRepo.Create(post)
	return thread.ThreadID, err
}

// Retrieve all threads after cursor
func (s *ThreadService) GetThreads(sort string, size int, descending bool, cursor time.Time, uid string) ([]models.Thread, error) {
	// Convert sort string to ThreadColumn
	column := models.StrToThreadColumn(sort)

	// Retrieve threads from db
	threads, err := s.threadRepo.GetAll(column, uid, cursor, size, descending)
	if err != nil {
		utils.Logger.Trace().Msg("Error getting threads from db")
		return nil, err
	}

	return threads, nil
}

// Retrieve threads metadata
func (s *ThreadService) GetThreadsMetadata() (models.ThreadsMetadata, error) {
	return s.threadRepo.GetMetadata()
}

// Retrieve thread and all associated posts
func (s *ThreadService) GetThread(threadID string, uid string) (*models.Thread, []models.Post, error) {
	// Retrieve thread from db
	thread, err := s.threadRepo.GetByID(threadID, uid)
	if err != nil {
		utils.Logger.Trace().Msg("Error getting thread")
		return nil, nil, err
	}

	// Retrieve posts from db
	posts, err := s.postRepo.GetPostsByThreadId(threadID, uid)
	if err != nil {
		utils.Logger.Trace().Msg("Error getting posts from db")
		return nil, nil, err
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
func (s *ThreadService) UpdateThread(threadID string, title string, content string, uid string) error {
	// Check if role is admin or staff
	hasStaffPermission, err := Users.HasStaffPermission(uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error checking staff permission")
		return err
	}

	if !hasStaffPermission {
		author, err := s.threadRepo.GetAuthor(threadID)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of thread")
			return err
		}

		// Check if author of thread matches the author in JWT claim
		if author != uid {
			return utils.NewErrUnauthorized()
		}
	}

	return s.threadRepo.Update(threadID, title, models.GetPreview(content))
}

// Delete thread and all associated posts
func (s *ThreadService) DeleteThread(threadID string, uid string) error {
	// Check if role is admin or staff
	hasStaffPermission, err := Users.HasStaffPermission(uid)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error checking staff permission")
		return err
	}

	if !hasStaffPermission {
		author, err := s.threadRepo.GetAuthor(threadID)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of thread")
			return err
		}

		// Check if author of thread matches the author in JWT claim
		if author != uid {
			return utils.NewErrUnauthorized()
		}
	}

	return s.threadRepo.Delete(threadID)
}
