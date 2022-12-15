package jwt

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Serasmi/home-library/internal/config"

	"github.com/golang-jwt/jwt"

	"github.com/Serasmi/home-library/pkg/logging"
)

type userCtxKeyType string

var UserCtxKey userCtxKeyType = "user"

type UserCtx struct {
	Id       string
	Username string
}

func withUser(ctx context.Context, user *UserCtx) context.Context {
	return context.WithValue(ctx, UserCtxKey, user)
}

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
			secret := config.GetConfig().JWT.Secret
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			logger.Error("Invalid access token", err)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Invalid access token"))

			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			logger.Error("Invalid token claims")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Invalid token claims"))

			return
		}

		ok = claims.VerifyExpiresAt(time.Now().UnixMilli(), true)
		if !ok {
			logger.Error("Unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Expired access token"))

			return
		}

		ctx := withUser(r.Context(), &UserCtx{Id: claims.UserId, Username: claims.Username})

		next(w, r.WithContext(ctx))
	}
}
