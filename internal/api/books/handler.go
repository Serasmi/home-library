package books

import (
	"github.com/Serasmi/home-library/internal/handlers"
	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	booksUrl = "/books"
	bookUrl  = "/books/:uuid"
)

type handler struct {
	logger  logging.Logger
	service Service
}

func NewHandler(service Service, logger logging.Logger) handlers.Handler {
	return &handler{logger, service}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, booksUrl, h.GetAll)
	router.HandlerFunc(http.MethodGet, bookUrl, h.GetById)
	router.HandlerFunc(http.MethodPost, bookUrl, h.Create)
	router.HandlerFunc(http.MethodPatch, bookUrl, h.PartiallyUpdate)
	router.HandlerFunc(http.MethodDelete, bookUrl, h.Delete)
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Need to implement method GetAll"))
}

func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Need to implement method GetById"))
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Need to implement method Create"))
}

func (h *handler) PartiallyUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Need to implement method PartiallyUpdate"))
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Need to implement method Delete"))
}
