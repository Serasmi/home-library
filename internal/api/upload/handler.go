package upload

import (
	"net/http"

	"github.com/Serasmi/home-library/internal/handlers"
	"github.com/Serasmi/home-library/internal/jwt"
	"github.com/julienschmidt/httprouter"

	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/Serasmi/home-library/pkg/uploader"
)

const (
	uploadURL = "/upload"
)

type handler struct {
	apiPath string
	service *Service
	logger  *logging.Logger
}

func NewHandler(apiPath string, service *Service, logger *logging.Logger) handlers.Handler {
	return &handler{apiPath, service, logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, h.apiPath+uploadURL, jwt.Protected(h.Upload, h.logger))
}

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("upload file")
	w.Header().Set("Content-Type", "application/json")

	meta := uploader.FileMeta{Filename: "uploaded.pdf"}

	err := h.service.Upload(r.Context(), r.Body, meta)
	if err != nil {
		h.logger.Error("file uploading error:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
