package books

import (
	"context"
	"errors"

	"github.com/Serasmi/home-library/pkg/logging"
)

type mockStorage struct {
	books  []Book
	logger *logging.Logger
}

func NewMockStorage(logger *logging.Logger) Storage {
	return &mockStorage{
		books:  initBooks(),
		logger: logger,
	}
}

func initBooks() []Book {
	return []Book{
		{
			ID:     "1",
			Title:  "War and Peace",
			Author: "Leo Tolstoy",
			Read:   false,
		},
		{
			ID:     "2",
			Title:  "The Brothers Karamazov",
			Author: "Fyodor Dostoevsky",
			Read:   false,
		},
		{
			ID:     "3",
			Title:  "The Master and Margarita",
			Author: "Mikhail Bulgakov",
			Read:   false,
		},
	}
}

func (s *mockStorage) Find(_ context.Context) ([]Book, error) {
	return s.books, nil
}

func (s *mockStorage) FindOne(_ context.Context, id string) (b Book, err error) {
	for _, b := range s.books {
		if b.ID == id {
			return b, nil
		}
	}

	return b, errors.New("book not found")
}

func (s *mockStorage) Insert(_ context.Context, book Book) (string, error) {
	s.books = append(s.books, book)
	return book.ID, nil
}

func (s *mockStorage) Update(_ context.Context, book UpdateBookDto) error {
	// TODO: implement
	s.logger.Debug("Update method not implemented")
	return nil
}

func (s *mockStorage) Remove(_ context.Context, id string) error {
	for i, b := range s.books {
		if b.ID == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			return nil
		}
	}

	return errors.New("book not found")
}
