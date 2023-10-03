package dao

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kkwon1/apod-forum-backend/internal/db"
	"github.com/kkwon1/apod-forum-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (dao *ApodDao) GetApodFromTo(startDate time.Time, endDate time.Time) []models.Apod {
	apodCollection := dao.dbClient.GetDatabase("apodDB").Collection("apod")
	filter := bson.M{
		"date": bson.M{
			"$gte": startDate.Format("2006-01-02"),
			"$lte": endDate.Format("2006-01-02"),
		},
	}
	cursor, err := apodCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal("Failed to read apod page")
	}
	defer cursor.Close(context.Background())

	var results []models.Apod
	for cursor.Next(context.Background()) {
		var apod models.Apod
		if err := cursor.Decode(&apod); err != nil {
			log.Fatal("Failed to decode APOD")
		}
		results = append(results, apod)
	}

	return results
}

func (dao *ApodDao) SearchApods(searchString string) []models.Apod {
	apodCollection := dao.dbClient.GetDatabase("apodDB").Collection("apod")
	pipeline := mongo.Pipeline{
		{{
			Key: "$search", Value: bson.M{
				"index": "textSearch",
				"text": bson.M{
					"query": searchString,
					"path": bson.M{
						"wildcard": "*",
					},
				},
			},
		}},
	}
	cursor, _ := apodCollection.Aggregate(context.Background(), pipeline)

	defer cursor.Close(context.Background())

	// Process the results
	var results []models.Apod
	for cursor.Next(context.Background()) {
		var apod models.Apod
		if err := cursor.Decode(&apod); err != nil {
			log.Fatal("Failed to decode APOD")
		}
		results = append(results, apod)
	}

	return results
}

func (dao *ApodDao) GetRandomApod() models.Apod {
	apodCollection := dao.dbClient.GetDatabase("apodDB").Collection("apod")

	pipeline := mongo.Pipeline{
		{{Key: "$sample", Value: bson.M{"size": 1}}},
	}
	cursor, _ := apodCollection.Aggregate(context.Background(), pipeline)

	defer cursor.Close(context.Background())

	var apod models.Apod
	// Check if there are results
	if cursor.Next(context.Background()) {
		err := cursor.Decode(&apod)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No random document found.")
	}
	return apod
}

func (dao *ApodDao) IncrementUpvoteCount(postId string) {
	apodCollection := dao.dbClient.GetDatabase("apodDB").Collection("apod")

	filter := bson.M{
		"date": postId,
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "upvoteCount", Value: 1}}}}
	opts := options.Update().SetUpsert(true)

	result, err := apodCollection.UpdateOne(context.Background(), filter, update, opts)

	log.Print(result)
	if err != nil {
		log.Fatal(err)
	}
}