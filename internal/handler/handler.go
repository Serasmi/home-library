package handler

import (
	"encoding/json"
	"github.com/Serasmi/home-library/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) HealthCheck() (path string, handler http.Handler) {
	path = "/api/health"

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Errorf("Error encoding healthcheck response: %s", err.Error())
		}
	}

	handler = http.HandlerFunc(handlerFunc)

	return
}
