package repository

import "github.com/Serasmi/home-library/internal/repository/mongorepo"

type Books interface {
	GetAllBooks() ([]string, error)
}

type Repository struct {
	Books
}

func New() *Repository {
	return &Repository{Books: mongorepo.NewBooks()}
}
