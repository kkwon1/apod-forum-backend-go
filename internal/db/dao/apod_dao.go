package dao

import (
	"context"

	"github.com/kkwon1/apod-forum-backend/internal/db"
	"github.com/kkwon1/apod-forum-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ApodDao struct {
	dbClient *db.MongoDBClient
}

func NewApodDao(client *db.MongoDBClient) (*ApodDao, error) {
	return &ApodDao{dbClient: client}, nil
}

func (dao *ApodDao) FindApod(date string) models.Apod {
	apodCollection := dao.dbClient.GetDatabase("apodDB").Collection("apod")
	var apod models.Apod
	filter := bson.M{"date": date}
	apodCollection.FindOne(context.Background(), filter).Decode(&apod)

	return apod
}
