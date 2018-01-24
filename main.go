package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/groups"
	"github.com/joyent/triton-service-groups/templates"
)

type authenticationHandler struct{}

func (handler *authenticationHandler) AuthenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: authenticate
		log.Printf("Authenticating against _something_\n")

		// http.Error(w, "Forbidden", 403)

		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)

	sub := router.PathPrefix("/v1/tsg").Subrouter()

	sub.Methods("POST").Path("/templates").HandlerFunc(templates_v1.Create)
	sub.Methods("GET").Path("/templates/{name}").HandlerFunc(templates_v1.Get)
	sub.Methods("PUT").Path("/templates/{name}").HandlerFunc(templates_v1.Update)
	sub.Methods("DELETE").Path("/templates/{name}").HandlerFunc(templates_v1.Delete)
	sub.Methods("GET").Path("/templates").HandlerFunc(templates_v1.List)

	sub.Methods("POST").Path("/").HandlerFunc(groups_v1.Create)
	sub.Methods("GET").Path("/{name}").HandlerFunc(groups_v1.Get)
	sub.Methods("PUT").Path("/{name}").HandlerFunc(groups_v1.Update)
	sub.Methods("DELETE").Path("/{name}").HandlerFunc(groups_v1.Delete)
	sub.Methods("GET").Path("/").HandlerFunc(groups_v1.List)

	log.Fatal(http.ListenAndServe(":3000", router))
}
