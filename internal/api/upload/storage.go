package upload

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Serasmi/home-library/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	CreateMeta(ctx context.Context, meta Meta) (string, error)
}

type mongoStorage struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewMongoStorage(storage *mongo.Database, collection string, logger *logging.Logger) Storage {
	return &mongoStorage{collection: storage.Collection(collection), logger: logger}
}

func (m *mongoStorage) CreateMeta(ctx context.Context, meta Meta) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := m.collection.InsertOne(ctx, meta)
	if err != nil {
		return "", fmt.Errorf("execute query: %w", err)
	}

	metaId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("convet objectid to hex")
	}

	return metaId.Hex(), nil
}
