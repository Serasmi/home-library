package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) GetAllBooks() (path string, handler http.Handler) {
	path = "/api/books"

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		books, err := h.services.Books.GetAllBooks()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Error(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(books)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Error(err.Error())
			return
		}
	}

	handler = http.HandlerFunc(handlerFunc)

	return
}
