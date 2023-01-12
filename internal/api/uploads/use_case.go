package uploads

import (
	"context"
	"io"

	"github.com/Serasmi/home-library/pkg/files"
	"github.com/Serasmi/home-library/pkg/logging"
)

type UseCase struct {
	storage      Storage
	fileProvider files.FileProvider
	logger       *logging.Logger
}

func NewUseCase(storage Storage, fileProvider files.FileProvider, logger *logging.Logger) *UseCase {
	return &UseCase{storage: storage, fileProvider: fileProvider, logger: logger}
}

func (uc UseCase) Upload(ctx context.Context, r io.ReadCloser, upload Upload) (filename string, err error) {
	err = uc.storage.UpdateUploadStatus(ctx, upload.ID, InProgress)
	if err != nil {
		return filename, err
	}

	err = uc.fileProvider.Upload(ctx, r, files.StoredFile{ID: upload.ID, Name: upload.Filename})
	if err != nil {
		statusErr := uc.storage.UpdateUploadStatus(ctx, upload.ID, Created)
		if statusErr != nil {
			return filename, statusErr
		}

		return filename, err
	}

	filename = upload.Filename

	err = uc.storage.UpdateUploadStatus(ctx, upload.ID, Done)
	if err != nil {
		return filename, err
	}

	return filename, nil
}

func (uc UseCase) CreateUpload(ctx context.Context, dto CreateUploadDTO) (string, error) {
	return uc.storage.CreateUpload(ctx, newUpload(dto))
}

func (uc UseCase) GetUploadByID(ctx context.Context, id string) (Upload, error) {
	return uc.storage.GetUploadByID(ctx, id)
}

func (uc UseCase) DeleteUpload(ctx context.Context, id string) error {
	return uc.storage.DeleteUpload(ctx, id)
}
