package repositories

import "github.com/kkwon1/apod-forum-backend/internal/db/dao"

type UserRepository struct {
	postUpvoteDao *dao.PostUpvoteDao
}

func NewUserRepository(postUpvoteDao *dao.PostUpvoteDao) (*UserRepository, error) {
	return &UserRepository{postUpvoteDao: postUpvoteDao}, nil
}

func (userRepo *UserRepository) GetUpvotedPostIds(userId string) []string {
	return userRepo.postUpvoteDao.GetUpvotedPostIds(userId)
}
