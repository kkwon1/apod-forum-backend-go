package repositories

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/kkwon1/apod-forum-backend/internal/db"
	"github.com/kkwon1/apod-forum-backend/internal/models"
)

var apodCache *lru.Cache[string, models.Apod]

type ApodRepository struct {
	dbClient *db.MongoDBClient
}

func NewApodRepository(client *db.MongoDBClient) (*ApodRepository, error) {
	// initialize LRU Cache with 3000 items
	apodCache, _ = lru.New[string, models.Apod](3000)

	return &ApodRepository{dbClient: client}, nil
}

func (ar *ApodRepository) GetApod(date string) models.Apod {
	if apodCache.Contains(date) {
		var apod, _ = apodCache.Get(date)
		return apod
	} else {
		apodCollection := ar.dbClient.GetDatabase("apodDB").Collection("apod")
		var apod models.Apod
		filter := bson.M{"date": date}
		apodCollection.FindOne(context.Background(), filter).Decode(&apod)

		apodCache.Add(date, apod)
		return apod
	}
}
