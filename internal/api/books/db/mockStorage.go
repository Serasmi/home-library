package db

import (
	"context"
	"errors"
	"github.com/Serasmi/home-library/internal/api/books"
	"github.com/Serasmi/home-library/pkg/logging"
)

type mockStorage struct {
	books  []books.Book
	logger logging.Logger
}

func NewMockStorage(logger logging.Logger) books.Storage {
	return &mockStorage{
		books:  initBooks(),
		logger: logger,
	}
}

func initBooks() []books.Book {
	return []books.Book{
		{
			Id:     "1",
			Title:  "War and Peace",
			Author: "Leo Tolstoy",
			Read:   false,
		},
		{
			Id:     "2",
			Title:  "The Brothers Karamazov",
			Author: "Fyodor Dostoevsky",
			Read:   false,
		},
		{
			Id:     "3",
			Title:  "The Master and Margarita",
			Author: "Mikhail Bulgakov",
			Read:   false,
		},
	}
}

func (s *mockStorage) Find(_ context.Context) ([]books.Book, error) {
	return s.books, nil
}

func (s *mockStorage) FindOne(_ context.Context, id string) (b books.Book, err error) {
	for _, b := range s.books {
		if b.Id == id {
			return b, nil
		}
	}
	return b, errors.New("book not found")
}

func (s *mockStorage) Remove(_ context.Context, id string) error {
	for i, b := range s.books {
		if b.Id == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			return nil
		}
	}
	return errors.New("book not found")
}
