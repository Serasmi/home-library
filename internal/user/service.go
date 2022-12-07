package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/Serasmi/home-library/pkg/logging"
)

type Service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(storage Storage, logger logging.Logger) *Service {
	return &Service{storage, logger}
}

func (s Service) CheckUser(ctx context.Context, username string, password string) error {
	s.logger.Infof("check user: %s", username)

	user, err := s.storage.FindByName(ctx, username)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
