package dao

import (
	"context"
	"log"

	"github.com/kkwon1/apod-forum-backend/cmd/db"
	"github.com/kkwon1/apod-forum-backend/cmd/models"
	"go.mongodb.org/mongo-driver/bson"
)

type PostUpvoteDao struct {
	dbClient *db.MongoDBClient
}

func NewPostUpvoteDao(client *db.MongoDBClient) (*PostUpvoteDao, error) {
	return &PostUpvoteDao{dbClient: client}, nil
}

func (dao *PostUpvoteDao) GetUpvotedPostIds(userId string) []string {
	postUpvoteCollection := dao.dbClient.GetDatabase("apodDB").Collection("postUpvote")

	filter := bson.M{
		"userSub": userId,
	}

	cursor, err := postUpvoteCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal("Failed to read from postUpvoteCollection", err)
	}
	defer cursor.Close(context.Background())

	var upvotes []models.Upvote
	for cursor.Next(context.Background()) {
		var upvote models.Upvote
		if err := cursor.Decode(&upvote); err != nil {
			log.Fatal("Failed to decode Upvote", err)
		}
		upvotes = append(upvotes, upvote)
	}

	var postIds []string
	for _, upvote := range upvotes {
		postIds = append(postIds, upvote.PostId)
	}

	return postIds
}

func (dao *PostUpvoteDao) UpvotePost(upvote models.Upvote) {
	postUpvoteCollection := dao.dbClient.GetDatabase("apodDB").Collection("postUpvote")
	_, err := postUpvoteCollection.InsertOne(context.Background(), upvote)

	if err != nil {
		log.Fatal("Failed to insert upvote into collection", err)
	}
}
