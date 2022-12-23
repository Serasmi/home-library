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

func filePath(filename string) string {
	return fmt.Sprintf("%s/%s", dirPath, filename)
}

func mustInitFS(logger *logging.Logger) {
	logger.Debugf("create uploads directory: %s", dirPath)

	err := os.Mkdir(dirPath, 0740)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

func (u fsUploader) Upload(_ context.Context, r io.ReadCloser, meta FileMeta) error {
	if meta.Filename == "" {
		return errors.New("empty filename")
	}

	f, err := os.Create(filePath(meta.ID))
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

	// TODO: add logic if file exist:
	//  1. Check that file exist in storage
	//  2. If exist return error
	//  3. Parse error in service and save file with index strategy
	//  4. Update filename in meta collection in DB
	//  5. Return new meta in response to client

	err = os.Rename(filePath(meta.ID), filePath(meta.Filename))
	if err != nil {
		u.logger.Error("renaming file error", err)
		// TODO: custom error is needed
		return errors.New("renaming file error")
	}

	return nil
}
