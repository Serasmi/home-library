package db

import (
	"context"
	"github.com/Serasmi/home-library/internal/api/books"
	"github.com/Serasmi/home-library/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     logging.Logger
}

func NewMongoStorage(storage *mongo.Database, collection string, logger logging.Logger) books.Storage {
	return &db{storage.Collection(collection), logger}
}

func (d *db) Find(ctx context.Context) ([]books.Book, error) {
	//TODO implement me
	return []books.Book{}, nil
}
