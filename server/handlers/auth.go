package handlers

import (
	"context"
	"net/http"

	"github.com/joyent/triton-service-groups/server/handlers/auth"
)

// authHandler encapsulates the authentication HTTP handler itself. We pipe all
// active HTTP requests through this object's ServeHTTP method.
type authHandler struct {
	handler http.Handler
}

// AuthHandler constructs and returns the HTTP handler object responsible for
// authenticating a request. This accepts a chain of HTTP handlers.
func AuthHandler(handler http.Handler) authHandler {
	return authHandler{
		handler: handler,
	}
}

// GetAuthSession pulls the current active authenticated session out of the
// current request context. This keeps authentication scoped to the active
// request only.
func GetAuthSession(ctx context.Context) auth.Session {
	if session, ok := ctx.Value(authKey).(auth.Session); ok {
		return session
	}
	return auth.Session{}
}

// ServeHTTP serves HTTP requests through the authentication process scoped to
// whatever pre-defined data we need accessible through the authHandler struct.
func (a authHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	session := auth.Session{
		AccountID: "joyent",
	}

	if !session.IsAuthenticated() {
		http.Error(w, ErrFailedAuth.Error(), http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(req.Context(), authKey, session)
	a.handler.ServeHTTP(w, req.WithContext(ctx))
}
