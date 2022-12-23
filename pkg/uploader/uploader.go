package uploader

import (
	"context"
	"io"
)

type FileMeta struct {
	ID       string `json:"ID"`
	Filename string `json:"filename"`
}

type Uploader interface {
	Upload(ctx context.Context, r io.ReadCloser, meta FileMeta) error
}
