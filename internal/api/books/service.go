package books

import (
	"context"

	"github.com/Serasmi/home-library/pkg/logging"
)

var _ Service = &service{}

type Service interface {
	GetAll(ctx context.Context) ([]Book, error)
	GetByID(ctx context.Context, id string) (Book, error)
	Create(ctx context.Context, dto CreateBookDto) (string, error)
	Update(ctx context.Context, dto UpdateBookDto) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	logger  *logging.Logger
	storage Storage
}

func NewService(storage Storage, logger *logging.Logger) Service {
	return &service{logger, storage}
}

func (s *service) GetAll(ctx context.Context) ([]Book, error) {
	// TODO: implement
	return s.storage.Find(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (Book, error) {
	// TODO implement me
	return s.storage.FindOne(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateBookDto) (string, error) {
	// TODO implement me
	return s.storage.Insert(ctx, newBook(dto))
}

func (s *service) Update(ctx context.Context, dto UpdateBookDto) error {
	// TODO implement me
	return s.storage.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	// TODO implement me
	return s.storage.Remove(ctx, id)
}
