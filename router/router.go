package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/session"
)

type Routes []Route

type DbWrapperHandler func(*session.TsgSession) http.HandlerFunc

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler DbWrapperHandler
}

func MakeRouter(session *session.TsgSession) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	for _, route := range templateRoutes {
		router.Handle(route.Pattern, route.Handler(session)).Methods(route.Method).Name(route.Name)
	}

	for _, route := range groupRoutes {
		router.Handle(route.Pattern, route.Handler(session)).Methods(route.Method).Name(route.Name)
	}

	return router
}
