package uploader

import (
	"context"
	"io"

	"github.com/Serasmi/home-library/pkg/logging"
)

type GoogleUploader struct {
	logger *logging.Logger
}

func NewGoogleUploader(logger *logging.Logger) Uploader {
	return &GoogleUploader{logger: logger}
}

func (u GoogleUploader) Upload(ctx context.Context, r io.ReadCloser, upload Upload) error {
	// TODO implement me
	panic("implement me")
}
