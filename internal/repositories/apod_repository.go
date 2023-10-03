package repositories

import (
	"time"

	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/kkwon1/apod-forum-backend/internal/db/dao"
	"github.com/kkwon1/apod-forum-backend/internal/models"
)

var apodCache *lru.Cache[string, models.Apod]

type ApodRepository struct {
	apodDao *dao.ApodDao
}

func NewApodRepository(apodDao *dao.ApodDao) (*ApodRepository, error) {
	// initialize LRU Cache with 3000 items
	apodCache, _ = lru.New[string, models.Apod](3000)

	return &ApodRepository{apodDao: apodDao}, nil
}

func (apodRepo *ApodRepository) GetApod(date string) models.Apod {
	if apodCache.Contains(date) {
		var apod, _ = apodCache.Get(date)
		return apod
	} else {
		apod := apodRepo.apodDao.FindApod(date)
		apodCache.Add(date, apod)
		return apod
	}
}

func (apodRepo *ApodRepository) GetApodsBetweenDates(startDate time.Time, endDate time.Time) []models.Apod {
	apods := apodRepo.apodDao.GetApodFromTo(startDate, endDate)
	for _, apod := range apods {
		apodCache.Add(apod.Date, apod)
	}
	return apods
}

func (apodRepo *ApodRepository) SearchApods(searchString string) []models.Apod {
	apods := apodRepo.apodDao.SearchApods(searchString)
	for _, apod := range apods {
		apodCache.Add(apod.Date, apod)
	}
	return apods
}

func (apodRepo *ApodRepository) GetRandomApod() models.Apod {
	apod := apodRepo.apodDao.GetRandomApod()
	apodCache.Add(apod.Date, apod)
	return apod
}

func (apodRepo *ApodRepository) GetApodPost(postId string) models.ApodPost {
	var post models.ApodPost
	if apodCache.Contains(postId) {
		var apod, _ = apodCache.Get(postId)

		post = models.ApodPost{
			NasaApod: apod,
		}
	} else {
		apod := apodRepo.apodDao.FindApod(postId)
		apodCache.Add(apod.Date, apod)

		// Stub comment tree for now
		post = models.ApodPost{
			NasaApod: apod,
			Comments: models.CommentTree{
				CommentID:    postId,
				Children:     []models.CommentTree{},
				CreateDate:   "2023-09-30",
				ModifiedDate: "2023-09-30",
				Comment:      "Test Comment",
				Author:       "Test Author",
				IsDeleted:    false,
				IsLeaf:       false,
			},
		}
	}

	return post
}

func (apodRepo *ApodRepository) IncrementUpvoteCount(postId string) {
	apodRepo.apodDao.IncrementUpvoteCount(postId)
	apodCache.Remove(postId)
}
