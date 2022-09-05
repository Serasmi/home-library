package service

import "github.com/Serasmi/home-library/internal/repository"

type Books interface {
	GetAllBooks() ([]string, error)
}

type Service struct {
	Books
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Books: NewBooksService(repo.Books),
	}
}
