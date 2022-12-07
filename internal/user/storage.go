package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Serasmi/home-library/pkg/logging"
)

type Storage interface {
	FindByName(ctx context.Context, username string) (User, error)
}

type mongoStorage struct {
	collection *mongo.Collection
	logger     logging.Logger
}

func NewMongoStorage(database *mongo.Database, collection string, logger logging.Logger) Storage {
	return &mongoStorage{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (s *mongoStorage) FindByName(ctx context.Context, username string) (u User, err error) {
	s.logger.Info("find user by username")

	filter := bson.M{"username": username}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result User

	err = s.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Error("user not found")
			return
		}

		s.logger.Error("db query error")

		return
	}

	return result, nil
}
