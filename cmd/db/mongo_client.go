package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	client *mongo.Client
	ctx    context.Context
}

func NewMongoDBClient(connectionString string) (*MongoDBClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &MongoDBClient{client: client, ctx: ctx}, nil
}

func (mc *MongoDBClient) Close() {
	mc.client.Disconnect(mc.ctx)
}

func (mc *MongoDBClient) GetDatabase(databaseName string) *mongo.Database {
	return mc.client.Database(databaseName)
}

// Example usage:
// dbClient, err := mongodb.NewMongoDBClient("mongodb://localhost:27017")
// if err != nil {
//     log.Fatal(err)
// }
// defer dbClient.Close()
// db := dbClient.GetDatabase("mydatabase")
