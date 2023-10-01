package repositories

import (
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
