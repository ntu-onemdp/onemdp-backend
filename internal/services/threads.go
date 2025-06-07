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

	err := s.threadRepo.Create(thread)
	if err != nil {
		return "", err
	}

	postFactory := models.PostFactory{}
	post := postFactory.New(thread.Author, thread.ThreadID, thread.Title, content, nil, true)

	err = s.postRepo.Create(post)
	return thread.ThreadID, err
}

// Retrieve all threads after cursor
func (s *ThreadService) GetThreads(sort string, size int, descending bool, cursor time.Time, username string) ([]models.Thread, error) {
	// Convert sort string to ThreadColumn
	column := models.StrToThreadColumn(sort)

	// Retrieve threads from db
	threads, err := s.threadRepo.GetAll(column, cursor, size, descending)
	if err != nil {
		utils.Logger.Trace().Msg("Error getting threads from db")
		return nil, err
	}

	// Retrieve number of likes and replies for each thread
	for i := range threads {
		threads[i].NumLikes = s.likesRepo.GetNumLikes(threads[i].ThreadID)
		threads[i].NumReplies = s.postRepo.GetNumReplies(threads[i].ThreadID)
		threads[i].IsLiked = s.likesRepo.GetByUsernameAndContentId(username, threads[i].ThreadID)
	}

	return threads, nil
}

// Retrieve threads metadata
func (s *ThreadService) GetThreadsMetadata() (models.ThreadsMetadata, error) {
	return s.threadRepo.GetMetadata()
}

// Retrieve thread and all associated posts
func (s *ThreadService) GetThread(threadID string, username string) (*models.Thread, []models.Post, error) {
	// Retrieve thread from db
	thread, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		utils.Logger.Trace().Msg("Error getting thread")
		return nil, nil, err
	}

	// Get number of likes for thread
	thread.NumLikes = s.likesRepo.GetNumLikes(threadID)

	// Retrieve posts from db
	posts, err := s.postRepo.GetPostByThreadId(threadID)
	if err != nil {
		utils.Logger.Trace().Msg("Error getting posts from db")
		return nil, nil, err
	}

	// Get number of likes for each post
	for i := range posts {
		if posts[i].IsHeader {
			posts[i].IsLiked = s.likesRepo.GetByUsernameAndContentId(username, threadID)
			posts[i].NumLikes = s.likesRepo.GetNumLikes(threadID)
		} else {
			posts[i].IsLiked = s.likesRepo.GetByUsernameAndContentId(username, posts[i].PostID)
			posts[i].NumLikes = s.likesRepo.GetNumLikes(posts[i].PostID)
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
func (s *ThreadService) UpdateThread(threadID string, title string, content string, claim *models.JwtClaim) error {
	// Check if role is admin or staff
	if !HasStaffPermission(claim) {
		author, err := s.threadRepo.GetAuthor(threadID)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of thread")
			return err
		}

		// Check if author of thread matches the author in JWT claim
		if author != claim.Uid {
			return utils.NewErrUnauthorized()
		}
	}

	return s.threadRepo.Update(threadID, title, models.GetPreview(content))
}

// Delete thread and all associated posts
func (s *ThreadService) DeleteThread(threadID string, claim *models.JwtClaim) error {
	// Check if role is admin or staff
	if !HasStaffPermission(claim) {
		author, err := s.threadRepo.GetAuthor(threadID)
		if author == "" || err != nil {
			utils.Logger.Error().Err(err).Msg("Error getting author of thread")
			return err
		}

		// Check if author of thread matches the author in JWT claim
		if author != claim.Uid {
			return utils.NewErrUnauthorized()
		}
	}

	return s.threadRepo.Delete(threadID)
}
