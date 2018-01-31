package router

import (
	"net/http"

	"github.com/joyent/triton-service-groups/session"
)

func isAuthenticated(session *session.TsgSession, r *http.Request) bool {
	session.AccountId = "joyent"
	return true
}

func AuthenticationHandler(session *session.TsgSession, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(session, r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
