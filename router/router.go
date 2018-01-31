package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

type Routes []Route

type DbWrapperHandler func(*pgx.ConnPool) http.HandlerFunc

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler DbWrapperHandler
}

func MakeRouter(dbPool *pgx.ConnPool) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	for _, route := range templateRoutes {
		router.Handle(route.Pattern, route.Handler(dbPool)).Methods(route.Method).Name(route.Name)
	}

	for _, route := range groupRoutes {
		router.Handle(route.Pattern, route.Handler(dbPool)).Methods(route.Method).Name(route.Name)
	}

	return router
}
