package files

import (
	"context"
	"io"
)

type StoredFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FileProvider interface {
	Download(ctx context.Context, filename string) (io.ReadCloser, error)
	Upload(ctx context.Context, r io.ReadCloser, file StoredFile) error
}
