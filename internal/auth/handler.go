package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Serasmi/home-library/internal/user"

	"github.com/Serasmi/home-library/internal/handlers"
	"github.com/Serasmi/home-library/internal/jwt"
	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

const (
	loginUrl = "/auth/login"
)

var _ handlers.Handler = &authHandler{}

type authHandler struct {
	userService *user.Service
	logger      logging.Logger
}

func NewHandler(userService *user.Service, logger logging.Logger) handlers.Handler {
	return &authHandler{userService, logger}
}

func (h *authHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, loginUrl, h.Login)
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Login user")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Debug("Decode login user dto")

	var dto user.LoginUserDto

	defer func() { _ = r.Body.Close() }()

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid data"))

		h.logger.Error("Decoding error", err)

		return
	}

	err = h.userService.CheckUser(r.Context(), dto.Username, dto.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid username or password"))

		h.logger.Error("Authorization error", err)

		return
	}

	token, err := jwt.CreateToken(&dto)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	respDTO, err := json.Marshal(LoginResponseDto{token})
	if err != nil {
		h.logger.Error("encoding response error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respDTO)
}
