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
	Upload(ctx context.Context, r io.ReadCloser, file StoredFile) error
}
