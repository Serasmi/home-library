package uploader

import (
	"context"
	"io"
)

type Upload struct {
	ID       string `json:"ID"`
	Filename string `json:"filename"`
}

type Uploader interface {
	Upload(ctx context.Context, r io.ReadCloser, upload Upload) error
}
