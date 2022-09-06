package mongorepo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoClient struct {
	client *mongo.Client
}

func New() *MongoClient {
	return &MongoClient{}
}

func (c *MongoClient) Init(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))
	if err != nil {
		return err
	}

	c.client = client

	return nil
}

func (c *MongoClient) Close(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return c.client.Disconnect(ctx)
}
