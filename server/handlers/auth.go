package handlers

import (
	"context"
	"net/http"
)

type authHandler struct {
	handler http.Handler
}

func AuthHandler(handler http.Handler) authHandler {
	return authHandler{
		handler: handler,
	}
}

type authSession struct {
	AccountID string
}

func (a authSession) IsAuthenticated() bool {
	return a.AccountID != ""
}

func GetAuthSession(ctx context.Context) (authSession, bool) {
	if session, ok := ctx.Value(authKey).(authSession); ok {
		return session, true
	}
	return authSession{}, false
}

func (a authHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO(justinwr): Actually authenticate the request against something.
	session := authSession{
		AccountID: "joyent",
	}

	if !session.IsAuthenticated() {
		http.Error(w, ErrFailedAuth.Error(), http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(req.Context(), authKey, session)
	a.handler.ServeHTTP(w, req.WithContext(ctx))
}
