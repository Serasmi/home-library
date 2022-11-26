package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/Serasmi/home-library/internal/api/books"
	"github.com/Serasmi/home-library/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type db struct {
	collection *mongo.Collection
	logger     logging.Logger
}

func NewMongoStorage(storage *mongo.Database, collection string, logger logging.Logger) books.Storage {
	return &db{storage.Collection(collection), logger}
}

func (d *db) Find(ctx context.Context) (books []books.Book, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return books, errors.New("books not found")
		}
		return books, fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err = result.All(ctx, &books); err != nil {
		return books, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return books, nil
}

func (d *db) FindOne(ctx context.Context, id string) (b books.Book, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return b, fmt.Errorf("failed to convert hex to objectId. error: %w", err)
	}

	filter := bson.M{"_id": objectId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := d.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return b, errors.New("book not found")
		}
		return b, fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err = result.Decode(&b); err != nil {
		return b, fmt.Errorf("failed to decode document. error: %w", err)
	}

	return b, nil
}

func (d *db) Insert(ctx context.Context, book books.Book) (id string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := d.collection.InsertOne(ctx, book)
	if err != nil {
		return id, fmt.Errorf("failed to create new book. error: %w", err)
	}

	bookId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return id, fmt.Errorf("failed to convet objectid to hex")
	}

	id = bookId.Hex()

	return id, nil
}

func (d *db) Update(ctx context.Context, book books.Book) error {
	//TODO implement me
	return nil
}

func (d *db) Remove(ctx context.Context, id string) error {
	//TODO implement me
	return nil
}
