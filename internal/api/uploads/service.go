package uploads

import (
	"context"
	"io"

	"github.com/Serasmi/home-library/pkg/files"

	"github.com/Serasmi/home-library/pkg/logging"
)

type Service struct {
	storage      Storage
	fileProvider files.FileProvider
	logger       *logging.Logger
}

func NewService(storage Storage, fileProvider files.FileProvider, logger *logging.Logger) *Service {
	return &Service{storage: storage, fileProvider: fileProvider, logger: logger}
}

func (s Service) Upload(ctx context.Context, r io.ReadCloser, upload Upload) (filename string, err error) {
	err = s.storage.UpdateUploadStatus(ctx, upload.ID, InProgress)
	if err != nil {
		return filename, err
	}

	err = s.fileProvider.Upload(ctx, r, files.StoredFile{ID: upload.ID, Name: upload.Filename})
	if err != nil {
		statusErr := s.storage.UpdateUploadStatus(ctx, upload.ID, Created)
		if statusErr != nil {
			return filename, statusErr
		}

		return filename, err
	}

	filename = upload.Filename

	err = s.storage.UpdateUploadStatus(ctx, upload.ID, Done)
	if err != nil {
		return filename, err
	}

	return filename, nil
}

func (s Service) CreateUpload(ctx context.Context, dto CreateUploadDTO) (string, error) {
	return s.storage.CreateUpload(ctx, newUpload(dto))
}

func (s Service) GetUploadByID(ctx context.Context, id string) (Upload, error) {
	return s.storage.GetUploadByID(ctx, id)
}

func (s Service) DeleteUpload(ctx context.Context, id string) error {
	return s.storage.DeleteUpload(ctx, id)
}
