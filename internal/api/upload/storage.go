package upload

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Serasmi/home-library/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	CreateMeta(ctx context.Context, meta Meta) (string, error)
	DeleteMeta(ctx context.Context, id string) error
}

type mongoStorage struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewMongoStorage(storage *mongo.Database, collection string, logger *logging.Logger) Storage {
	return &mongoStorage{collection: storage.Collection(collection), logger: logger}
}

func (s *mongoStorage) CreateMeta(ctx context.Context, meta Meta) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := s.collection.InsertOne(ctx, meta)
	if err != nil {
		return "", fmt.Errorf("execute query: %w", err)
	}

	metaId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("convet objectid to hex")
	}

	return metaId.Hex(), nil
}

func (s *mongoStorage) DeleteMeta(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectID. error: %w", err)
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query")
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}
