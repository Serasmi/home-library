package service

import "github.com/Serasmi/home-library/internal/repository"

type BooksService struct {
	repo repository.Books
}

func (b *BooksService) GetAllBooks() ([]string, error) {
	return b.repo.GetAllBooks()
}

func NewBooksService(repo repository.Books) *BooksService {
	return &BooksService{repo: repo}
}
