package books

import (
	"context"
	"github.com/Serasmi/home-library/pkg/logging"
)

type Service interface {
	GetAll(ctx context.Context) ([]Book, error)
}

type service struct {
	logger  logging.Logger
	storage Storage
}

func NewService(storage Storage, logger logging.Logger) Service {
	return &service{logger, storage}
}

func (s *service) GetAll(ctx context.Context) ([]Book, error) {
	// TODO: implement
	return s.storage.Find(ctx)
}
