package uploads

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

var _ Storage = (*mongoStorage)(nil)

type Storage interface {
	CreateUpload(ctx context.Context, upload Upload) (string, error)
	GetUploadByID(ctx context.Context, id string) (Upload, error)
	GetUploadNameByBookID(ctx context.Context, bookId string) (string, error)
	DeleteUpload(ctx context.Context, id string) error
	UpdateUploadStatus(ctx context.Context, id string, status Status) error
}

type mongoStorage struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewMongoStorage(storage *mongo.Database, collection string, logger *logging.Logger) Storage {
	return &mongoStorage{collection: storage.Collection(collection), logger: logger}
}

func (s *mongoStorage) CreateUpload(ctx context.Context, upload Upload) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := s.collection.InsertOne(ctx, upload)
	if err != nil {
		return "", fmt.Errorf("execute query: %w", err)
	}

	uploadId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("convet objectid to hex")
	}

	return uploadId.Hex(), nil
}

func (s *mongoStorage) GetUploadByID(ctx context.Context, id string) (u Upload, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectID. error: %w", err)
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := s.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return u, errors.New("upload not found")
		}

		return u, fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode document. error: %w", err)
	}

	return u, nil
}

func (s *mongoStorage) GetUploadNameByBookID(ctx context.Context, bookId string) (name string, err error) {
	filter := bson.M{"bookId": bookId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := s.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return name, errors.New("book not found")
		}

		return name, fmt.Errorf("failed to execute query. error: %w", err)
	}

	var upload Upload
	if err = result.Decode(&upload); err != nil {
		return name, fmt.Errorf("failed to decode document. error: %w", err)
	}

	return upload.Filename, nil
}

func (s *mongoStorage) DeleteUpload(ctx context.Context, id string) error {
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

func (s *mongoStorage) UpdateUploadStatus(ctx context.Context, id string, status Status) error {
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
