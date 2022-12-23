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
)

const (
	uploadURL = "/upload/:id"
	// TODO: rename meta to uploads
	metaURL    = "/meta"
	oneMetaURL = "/meta/:id"
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
	router.HandlerFunc(http.MethodPut, h.apiPath+uploadURL, jwt.Protected(h.Upload, h.logger))

	router.HandlerFunc(http.MethodPost, h.apiPath+metaURL, jwt.Protected(h.CreateMeta, h.logger))
	router.HandlerFunc(http.MethodDelete, h.apiPath+oneMetaURL, jwt.Protected(h.DeleteMeta, h.logger))
}

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("upload file")
	w.Header().Set("Content-Type", "application/json")

	id, err := handlers.RequestID(r, h.logger)
	if err != nil {
		h.logger.Error("id parameter is required in request path:", err)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "upload id is required")

		return
	}

	meta, err := h.service.GetMetaById(r.Context(), id)
	if err != nil {
		h.logger.Error("finding meta error:", err)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "upload not found")

		return
	}

	if meta.Status != Created {
		h.logger.Error("file has already uploaded. Status:", meta.Status)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "file has already uploaded")

		return
	}

	filename, err := h.service.Upload(r.Context(), r.Body, meta)
	if err != nil {
		h.logger.Error("file uploading error:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)

	resDto, err := json.Marshal(ResponseDTO{filename})
	if err != nil {
		h.logger.Error(w, "marshaling error")
		return
	}

	_, _ = w.Write(resDto)
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

func (h *handler) DeleteMeta(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("delete meta")
	w.Header().Set("Content-Type", "application/json")

	id, err := handlers.RequestID(r, h.logger)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "id parameter is required in request path")

		return
	}

	err = h.service.DeleteMeta(r.Context(), id)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
