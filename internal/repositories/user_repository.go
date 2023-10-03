package repositories

import (
	"time"

	"github.com/kkwon1/apod-forum-backend/internal/db/dao"
	"github.com/kkwon1/apod-forum-backend/internal/models"
	"github.com/kkwon1/apod-forum-backend/internal/models/requests"
)

type UserRepository struct {
	postUpvoteDao *dao.PostUpvoteDao
}

func NewUserRepository(postUpvoteDao *dao.PostUpvoteDao) (*UserRepository, error) {
	return &UserRepository{postUpvoteDao: postUpvoteDao}, nil
}

func (userRepo *UserRepository) GetUpvotedPostIds(userId string) []string {
	return userRepo.postUpvoteDao.GetUpvotedPostIds(userId)
}

func (userRepo *UserRepository) UpvotePost(req requests.UpvotePostRequest) {
	userRepo.postUpvoteDao.UpvotePost(models.Upvote{
		PostId: req.PostId,
		UserSub: req.UserSub,
		Timestamp: time.Now(),
	})
}