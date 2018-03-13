package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type RouteTable []Routes

type Routes []Route

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

func WithRoutes(routes RouteTable) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, rs := range routes {
		for _, r := range rs {
			router.Path(r.Pattern).
				Methods(r.Method).
				Name(r.Name).
				HandlerFunc(r.Handler)
		}
	}

	return router
}
