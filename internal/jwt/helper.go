package jwt

import (
	"time"

	"github.com/Serasmi/home-library/internal/user"

	"github.com/golang-jwt/jwt"
)

const (
	jwtKey = "fjsadomwoi3475872364895n23t4hf9328n4ytv6c2"
	jwtTTL = 3 * time.Hour
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(u *user.LoginUserDto) (string, error) {
	expired := time.Now().Add(jwtTTL)

	claims := &Claims{Username: u.Username, StandardClaims: jwt.StandardClaims{
		ExpiresAt: expired.Unix(),
	}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
