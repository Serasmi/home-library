package books

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Serasmi/home-library/internal/jwt"

	"github.com/Serasmi/home-library/internal/handlers"
	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	booksURL    = "/books"
	bookURL     = "/books/:id"
	downloadURL = "/download/:id"
)

type handler struct {
	apiPath string
	useCase UseCase
	logger  *logging.Logger
}

func NewHandler(apiPath string, useCase UseCase, logger *logging.Logger) handlers.Handler {
	return &handler{apiPath, useCase, logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, h.apiPath+booksURL, jwt.Protected(h.GetAll, h.logger))
	router.HandlerFunc(http.MethodGet, h.apiPath+bookURL, jwt.Protected(h.GetByID, h.logger))
	router.HandlerFunc(http.MethodPost, h.apiPath+booksURL, jwt.Protected(h.Create, h.logger))
	router.HandlerFunc(http.MethodPatch, h.apiPath+bookURL, jwt.Protected(h.PartiallyUpdate, h.logger))
	router.HandlerFunc(http.MethodDelete, h.apiPath+bookURL, jwt.Protected(h.Delete, h.logger))
	router.HandlerFunc(http.MethodGet, h.apiPath+downloadURL, jwt.Protected(h.Download, h.logger))
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Get all books")
	w.Header().Set("Content-Type", "application/json")

	books, err := h.useCase.GetAll(r.Context())
	if err != nil {
		h.logger.Error(err)
	}

	booksBytes, err := json.Marshal(books)
	if err != nil {
		h.logger.Error(err)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(booksBytes)
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Get book by id")
	w.Header().Set("Content-Type", "application/json")

	id, err := handlers.RequestID(r, h.logger)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "id parameter is required in request path")

		return
	}

	book, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error(err)
	}

	bookBytes, err := json.Marshal(book)
	if err != nil {
		h.logger.Error(err)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bookBytes)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Create new book")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Debug("Decode create book dto")

	var dto CreateBookDto

	defer func() { _ = r.Body.Close() }()

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "invalid data")

		return
	}

	id, err := h.useCase.Create(r.Context(), dto)
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

	resDto, err := json.Marshal(CreateBookResponseDto{id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "creating entity server error")

		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resDto)
}

func (h *handler) PartiallyUpdate(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Partially update book")
	w.Header().Set("Content-Type", "application/json")

	id, err := handlers.RequestID(r, h.logger)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "id parameter is required in request path")

		return
	}

	h.logger.Debug("Decode update book dto")

	var dto UpdateBookDto

	defer func() { _ = r.Body.Close() }()

	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "invalid data")

		return
	}

	dto.ID = id

	err = h.useCase.Update(r.Context(), dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "updating entity server error")

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Delete book")
	w.Header().Set("Content-Type", "application/json")

	id, err := handlers.RequestID(r, h.logger)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "id parameter is required in request path")

		return
	}

	if err = h.useCase.Delete(r.Context(), id); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) Download(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Download book")

	id, err := handlers.RequestID(r, h.logger)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "id parameter is required in request path")

		return
	}

	bytes, err := h.useCase.Download(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "download file error")

		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bytes)
}
