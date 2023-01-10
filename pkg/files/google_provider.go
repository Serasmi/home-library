package files

import (
	"context"
	"io"

	"github.com/Serasmi/home-library/pkg/logging"
)

type googleProvider struct {
	logger *logging.Logger
}

func NewGoogleProvider(logger *logging.Logger) FileProvider {
	return &googleProvider{logger: logger}
}

func (u googleProvider) Upload(ctx context.Context, r io.ReadCloser, file StoredFile) error {
	// TODO implement me
	panic("implement me")
}
