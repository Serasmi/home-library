package upload

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Serasmi/home-library/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	CreateMeta(ctx context.Context, meta Meta) (string, error)
	GetMetaById(ctx context.Context, id string) (Meta, error)
	DeleteMeta(ctx context.Context, id string) error
	UpdateMetaStatus(ctx context.Context, id string, status Status) error
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

func (s *mongoStorage) GetMetaById(ctx context.Context, id string) (m Meta, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return m, fmt.Errorf("failed to convert hex to objectID. error: %w", err)
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := s.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return m, errors.New("meta not found")
		}

		return m, fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err = result.Decode(&m); err != nil {
		return m, fmt.Errorf("failed to decode document. error: %w", err)
	}

	return m, nil
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

func (s *mongoStorage) UpdateMetaStatus(ctx context.Context, id string, status Status) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex to objectID. error: %w", err)
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"status": status}}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("book not found")
	}

	return nil
}
