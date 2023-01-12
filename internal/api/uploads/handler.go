package uploads

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
	uploadsURL = "/uploads"
	uploadURL  = "/uploads/:id"
)

type handler struct {
	apiPath string
	useCase *UseCase
	logger  *logging.Logger
}

func NewHandler(apiPath string, useCase *UseCase, logger *logging.Logger) handlers.Handler {
	return &handler{apiPath, useCase, logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPut, h.apiPath+uploadURL, jwt.Protected(h.Upload, h.logger))
	router.HandlerFunc(http.MethodPost, h.apiPath+uploadsURL, jwt.Protected(h.CreateUpload, h.logger))
	router.HandlerFunc(http.MethodDelete, h.apiPath+uploadURL, jwt.Protected(h.DeleteUpload, h.logger))
}

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("uploads file")
	w.Header().Set("Content-Type", "application/json")

	id, err := handlers.RequestID(r, h.logger)
	if err != nil {
		h.logger.Error("id parameter is required in request path:", err)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "uploads id is required")

		return
	}

	upload, err := h.useCase.GetUploadByID(r.Context(), id)
	if err != nil {
		h.logger.Error("finding uploads error:", err)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "uploads not found")

		return
	}

	if upload.Status != Created {
		h.logger.Error("file has already uploaded. Status:", upload.Status)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "file has already uploaded")

		return
	}

	filename, err := h.useCase.Upload(r.Context(), r.Body, upload)
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

func (h *handler) CreateUpload(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("create uploads")
	w.Header().Set("Content-Type", "application/json")

	var dto CreateUploadDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		h.logger.Error("decode uploads:", err)

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

	id, err := h.useCase.CreateUpload(r.Context(), dto)
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

	resDTO, err := json.Marshal(CreateUploadResponseDTO{id})
	if err != nil {
		h.logger.Error("encode uploads response:", err)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "creating entity server error")

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resDTO)
}

func (h *handler) DeleteUpload(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("delete uploads")
	w.Header().Set("Content-Type", "application/json")

	id, err := handlers.RequestID(r, h.logger)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "id parameter is required in request path")

		return
	}

	err = h.useCase.DeleteUpload(r.Context(), id)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
