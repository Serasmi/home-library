package upload

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Serasmi/home-library/internal/handlers"
	"github.com/Serasmi/home-library/internal/jwt"
	"github.com/julienschmidt/httprouter"

	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/Serasmi/home-library/pkg/uploader"
)

const (
	uploadURL = "/upload"
	metaURL   = "/upload/meta"
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
	router.HandlerFunc(http.MethodPost, h.apiPath+metaURL, jwt.Protected(h.CreateMeta, h.logger))
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

func (h *handler) CreateMeta(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("create meta")
	w.Header().Set("Content-Type", "application/json")

	var dto CreateMetaDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		h.logger.Error("decode meta:", err)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "invalid data")

		return
	}

	// TODO: create validator
	if dto.Filename == "" {
		h.logger.Error("empty filename")

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "empty filename")

		return
	}

	id, err := h.service.CreateMeta(r.Context(), dto)
	if err != nil {
		h.logger.Error(err)

		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "duplicated entity")

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "creating entity server error")

		return
	}

	resDTO, err := json.Marshal(CreateMetaResponseDTO{id})
	if err != nil {
		h.logger.Error("encode meta response:", err)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "creating entity server error")

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resDTO)
}
