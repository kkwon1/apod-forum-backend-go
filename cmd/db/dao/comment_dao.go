package dao

import (
	"context"
	"log"

	"github.com/kkwon1/apod-forum-backend/cmd/db"
	"github.com/kkwon1/apod-forum-backend/cmd/models"
	"go.mongodb.org/mongo-driver/bson"
)

type CommentDao struct {
	dbClient *db.MongoDBClient
}

func NewCommentDao(client *db.MongoDBClient) (*CommentDao, error) {
	return &CommentDao{dbClient: client}, nil
}

func (dao *CommentDao) GetCommentsByPostId(postId string) ([]models.Comment, error) {
	commentsCollection := dao.dbClient.GetDatabase("apodDB").Collection("comments")
	filter := bson.M{"postId": postId}
	cursor, err := commentsCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var comments []models.Comment
	for cursor.Next(context.Background()) {
		var comment models.Comment
		if err := cursor.Decode(&comment); err != nil {
			log.Fatal("Failed to decode Comment", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}