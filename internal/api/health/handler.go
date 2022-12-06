package health

import (
	"net/http"

	"github.com/Serasmi/home-library/internal/handlers"
	"github.com/julienschmidt/httprouter"
)

const (
	healthURL = "/health"
)

type handler struct {
	apiPath string
}

func NewHandler(apiPath string) handlers.Handler {
	return &handler{apiPath}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, h.apiPath+healthURL, h.Health)
}

func (h *handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
