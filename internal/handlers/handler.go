package handlers

import (
	"errors"
	"net/http"

	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type Handler interface {
	Register(router *httprouter.Router)
}

func RequestID(r *http.Request, logger *logging.Logger) (string, error) {
	logger.Debug("Get id from request params")

	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	if id == "" {
		return "", errors.New("id not found in request")
	}

	return id, nil
}
