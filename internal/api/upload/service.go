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

func (s Service) Upload(ctx context.Context, r io.ReadCloser, meta Meta) (filename string, err error) {
	err = s.storage.UpdateMetaStatus(ctx, meta.ID, InProgress)
	if err != nil {
		return filename, err
	}

	err = s.uploader.Upload(ctx, r, uploader.FileMeta{ID: meta.ID, Filename: meta.Filename})
	if err != nil {
		statusErr := s.storage.UpdateMetaStatus(ctx, meta.ID, Created)
		if statusErr != nil {
			return filename, statusErr
		}

		return filename, err
	}

	filename = meta.Filename

	err = s.storage.UpdateMetaStatus(ctx, meta.ID, Done)
	if err != nil {
		return filename, err
	}

	return filename, nil
}

func (s Service) CreateMeta(ctx context.Context, dto CreateMetaDTO) (string, error) {
	return s.storage.CreateMeta(ctx, newMeta(dto))
}

func (s Service) GetMetaById(ctx context.Context, id string) (Meta, error) {
	return s.storage.GetMetaById(ctx, id)
}

func (s Service) DeleteMeta(ctx context.Context, id string) error {
	return s.storage.DeleteMeta(ctx, id)
}
