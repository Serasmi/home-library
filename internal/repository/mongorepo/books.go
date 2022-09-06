package mongorepo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Books struct {
	collection *mongo.Collection
}

func NewBooks(mongoClient *MongoClient) *Books {
	return &Books{
		collection: mongoClient.client.Database("HomeLibrary").Collection("Books"),
	}
}

func (b *Books) GetAllBooks() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cur, err := b.collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := cur.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
		fmt.Println(result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return []string{"Book 1", "Book 2", "Book 3"}, nil
}
