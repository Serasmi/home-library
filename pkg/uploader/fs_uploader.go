package uploader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Serasmi/home-library/pkg/logging"
)

const (
	dirPath = "uploads"
)

type fsUploader struct {
	logger *logging.Logger
}

func NewFSUploader(logger *logging.Logger) Uploader {
	mustInitFS(logger)

	return &fsUploader{
		logger,
	}
}

func filePath(f FileMeta) string {
	return fmt.Sprintf("%s/%s", dirPath, f.Filename)
}

func mustInitFS(logger *logging.Logger) {
	logger.Debugf("create uploads directory: %s", dirPath)

	err := os.Mkdir(dirPath, 0740)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

func (u fsUploader) Upload(ctx context.Context, r io.ReadCloser, meta FileMeta) error {
	if meta.Filename == "" {
		return errors.New("empty filename")
	}

	f, err := os.Create(filePath(meta))
	if err != nil {
		u.logger.Error("file creating error:", err)
		// TODO: custom error is needed
		return errors.New("file creating error")
	}

	defer func() { _ = f.Close() }()

	b, err := io.Copy(f, r)
	if err != nil {
		u.logger.Error("data saving error", err)
		// TODO: custom error is needed
		return errors.New("data saving error")
	}

	defer func() { _ = r.Close() }()

	u.logger.Debugf("written %d bytes", b)

	return nil
}
