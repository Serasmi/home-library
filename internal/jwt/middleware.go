package jwt

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/Serasmi/home-library/pkg/logging"
)

type username string

func Protected(next http.HandlerFunc, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			logger.Error("Invalid authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Invalid authorization header"))

			return
		}

		tokenString := authHeader[1]

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSigningKey), nil
		})
		if err != nil || !token.Valid {
			logger.Error("Invalid access token", err)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Invalid access token"))

			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || claims.Username == "" {
			logger.Error("Invalid token claims")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Invalid token claims"))

			return
		}

		ctx := context.WithValue(r.Context(), username("username"), claims.Username)

		next(w, r.WithContext(ctx))
	}
}
