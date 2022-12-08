package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	jwtSigningKey = "fjsadomwoi3475872364895n23t4hf9328n4ytv6c2"
	jwtTTL        = 3 * time.Hour
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

	tokenString, err := token.SignedString([]byte(jwtSigningKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
