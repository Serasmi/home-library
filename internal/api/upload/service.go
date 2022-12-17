package upload

import (
	"context"
	"io"

	"github.com/Serasmi/home-library/pkg/uploader"

	"github.com/Serasmi/home-library/pkg/logging"
)

type Service struct {
	uploader uploader.Uploader
	logger   *logging.Logger
}

func NewService(uploader uploader.Uploader, logger *logging.Logger) *Service {
	return &Service{uploader: uploader, logger: logger}
}

func (s Service) Upload(ctx context.Context, r io.ReadCloser, meta uploader.FileMeta) error {
	return s.uploader.Upload(ctx, r, meta)
}
