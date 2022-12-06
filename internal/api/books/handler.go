package books

import (
	"encoding/json"
	"fmt"
	"github.com/Serasmi/home-library/internal/handlers"
	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

const (
	booksURL = "/books"
	bookURL  = "/books/:id"
)

type handler struct {
	apiPath string
	logger  logging.Logger
	service Service
}

func NewHandler(apiPath string, service Service, logger logging.Logger) handlers.Handler {
	return &handler{apiPath, logger, service}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, h.apiPath+booksURL, h.GetAll)
	router.HandlerFunc(http.MethodGet, h.apiPath+bookURL, h.GetByID)
	router.HandlerFunc(http.MethodPost, h.apiPath+booksURL, h.Create)
	router.HandlerFunc(http.MethodPatch, h.apiPath+bookURL, h.PartiallyUpdate)
	router.HandlerFunc(http.MethodDelete, h.apiPath+bookURL, h.Delete)
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Get all books")
	w.Header().Set("Content-Type", "application/json")

	books, err := h.service.GetAll(r.Context())
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
		fmt.Fprint(w, "id parameter is required in request path")

		return
	}

	book, err := h.service.GetByID(r.Context(), id)
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

	defer r.Body.Close() //nolint:errcheck

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid data")

		return
	}

	id, err := h.service.Create(r.Context(), dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)

		if mongo.IsDuplicateKeyError(err) {
			fmt.Fprint(w, "duplicated key")
			return
		}

		fmt.Fprint(w, "creating entity server error")

		return
	}

	resDto, err := json.Marshal(CreateBookResponseDto{id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "creating entity server error")

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
		fmt.Fprint(w, "id parameter is required in request path")

		return
	}

	h.logger.Debug("Decode update book dto")

	var dto UpdateBookDto

	defer r.Body.Close() //nolint:errcheck

	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid data")

		return
	}

	dto.ID = id

	err = h.service.Update(r.Context(), dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "updating entity server error")

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
		fmt.Fprint(w, "id parameter is required in request path")

		return
	}

	if err = h.service.Delete(r.Context(), id); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
