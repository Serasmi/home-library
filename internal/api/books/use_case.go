package books

import (
	"context"
	"io"

	"github.com/Serasmi/home-library/pkg/logging"
)

var _ UseCase = (*useCase)(nil)

type UseCase interface {
	GetAll(ctx context.Context) ([]Book, error)
	GetByID(ctx context.Context, id string) (Book, error)
	Create(ctx context.Context, dto CreateBookDto) (string, error)
	Update(ctx context.Context, dto UpdateBookDto) error
	Delete(ctx context.Context, id string) error
	Download(ctx context.Context, bookId string) (io.ReadCloser, error)
}

type useCase struct {
	logger        *logging.Logger
	storage       Storage
	uploadStorage UploadStorage
	fileProvider  FileProvider
}

type FileProvider interface {
	Download(ctx context.Context, fileName string) (io.ReadCloser, error)
}

type UploadStorage interface {
	GetUploadNameByBookID(ctx context.Context, bookId string) (string, error)
}

func NewUseCase(storage Storage, uploadStorage UploadStorage, fileProvider FileProvider, logger *logging.Logger) UseCase {
	return &useCase{logger, storage, uploadStorage, fileProvider}
}

func (uc *useCase) GetAll(ctx context.Context) ([]Book, error) {
	// TODO: implement
	return uc.storage.Find(ctx)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (Book, error) {
	// TODO implement me
	return uc.storage.FindOne(ctx, id)
}

func (uc *useCase) Create(ctx context.Context, dto CreateBookDto) (string, error) {
	// TODO implement me
	return uc.storage.Insert(ctx, newBook(dto))
}

func (uc *useCase) Update(ctx context.Context, dto UpdateBookDto) error {
	// TODO implement me
	return uc.storage.Update(ctx, dto)
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	// TODO implement me
	return uc.storage.Remove(ctx, id)
}

func (uc *useCase) Download(ctx context.Context, bookID string) (io.ReadCloser, error) {
	fileName, err := uc.uploadStorage.GetUploadNameByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	return uc.fileProvider.Download(ctx, fileName)
}
