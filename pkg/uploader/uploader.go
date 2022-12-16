package uploader

import (
	"context"
	"io"
)

type FileMeta struct {
	Filename string
}

type Uploader interface {
	Upload(ctx context.Context, r io.ReadCloser, meta FileMeta) error
}