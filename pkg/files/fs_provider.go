package files

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Serasmi/home-library/pkg/logging"
)

const (
	dirPath = "files"
)

type fsProvider struct {
	logger *logging.Logger
}

func NewFSProvider(logger *logging.Logger) FileProvider {
	mustInitFS(logger)

	return &fsProvider{
		logger,
	}
}

func filePath(filename string) string {
	return fmt.Sprintf("%s/%s", dirPath, filename)
}

func mustInitFS(logger *logging.Logger) {
	logger.Debugf("create files directory: %s", dirPath)

	err := os.Mkdir(dirPath, 0740)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

func (u fsProvider) Upload(_ context.Context, r io.ReadCloser, file StoredFile) error {
	if file.Name == "" {
		return errors.New("empty filename")
	}

	f, err := os.Create(filePath(file.ID))
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
	//  4. Update filename in upload collection in DB
	//  5. Return new upload in response to client

	err = os.Rename(filePath(file.ID), filePath(file.Name))
	if err != nil {
		u.logger.Error("renaming file error", err)
		// TODO: custom error is needed
		return errors.New("renaming file error")
	}

	return nil
}
