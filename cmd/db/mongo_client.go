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
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
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
