package files

import (
	"context"
	"io"

	"github.com/Serasmi/home-library/pkg/logging"
)

var _ FileProvider = (*googleProvider)(nil)

type googleProvider struct {
	logger *logging.Logger
}

func NewGoogleProvider(logger *logging.Logger) FileProvider {
	return &googleProvider{logger: logger}
}

func (p googleProvider) Download(ctx context.Context, filename string) (io.ReadCloser, error) {
	// TODO implement me
	panic("implement me")
}

func (p googleProvider) Upload(ctx context.Context, r io.ReadCloser, file StoredFile) error {
	// TODO implement me
	panic("implement me")
}
