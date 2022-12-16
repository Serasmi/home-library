package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/Serasmi/home-library/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func NewService(storage Storage, logger *logging.Logger) *Service {
	return &Service{storage, logger}
}

func (s Service) FindUser(ctx context.Context, username string, password string) (id string, err error) {
	s.logger.Infof("find user: %s", username)

	user, err := s.storage.FindByName(ctx, username)
	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return id, err
	}

	return user.Id, nil
}
