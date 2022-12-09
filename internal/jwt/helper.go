package jwt

import (
	"time"

	"github.com/Serasmi/home-library/internal/config"

	"github.com/golang-jwt/jwt"
)

const (
	jwtTTL = 3 * time.Hour
)

type Claims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(id string, username string) (string, error) {
	expired := time.Now().Add(jwtTTL)

	claims := &Claims{
		UserId:   id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expired.UnixMilli(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.GetConfig().JWT.Secret

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
