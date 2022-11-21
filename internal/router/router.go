package router

import (
	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(logger logging.Logger) (r *httprouter.Router) {
	logger.Infof("Create router")

	r = httprouter.New()

	return
}
