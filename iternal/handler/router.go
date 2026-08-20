package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(h CounterHandler) http.Handler {
	router := mux.NewRouter()
	router.Path("/health").Methods(http.MethodGet).HandlerFunc(h.Health)
	router.Path("/counter").Methods(http.MethodGet).HandlerFunc(h.GetCount)
	router.Path("/counter/increment").Methods(http.MethodPost).HandlerFunc(h.IncrCount)
	return router

}
