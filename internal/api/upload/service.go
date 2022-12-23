package upload

import (
	"context"
	"io"

	"github.com/Serasmi/home-library/pkg/uploader"

	"github.com/Serasmi/home-library/pkg/logging"
)

type Service struct {
	storage  Storage
	uploader uploader.Uploader
	logger   *logging.Logger
}

func NewService(storage Storage, uploader uploader.Uploader, logger *logging.Logger) *Service {
	return &Service{storage: storage, uploader: uploader, logger: logger}
}

func (s Service) Upload(ctx context.Context, r io.ReadCloser, meta uploader.FileMeta) error {
	return s.uploader.Upload(ctx, r, meta)
}

func (s Service) CreateMeta(ctx context.Context, dto CreateMetaDTO) (string, error) {
	return s.storage.CreateMeta(ctx, newMeta(dto))
}

func (s Service) DeleteMeta(ctx context.Context, id string) error {
	return s.storage.DeleteMeta(ctx, id)
}
