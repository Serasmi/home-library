package repository

import (
	"github.com/Serasmi/home-library/internal/repository/mongorepo"
)

func NewMongoRepository(mongoClient *mongorepo.MongoClient) *Repository {
	return &Repository{
		Books: mongorepo.NewBooks(mongoClient),
	}
}
