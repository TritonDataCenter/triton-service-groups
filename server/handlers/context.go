package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx"
)

type contextKey int

const (
	dbKeyName contextKey = iota
	authKey
)

type dbValue struct {
	pool *pgx.ConnPool
}

func GetDBPool(ctx context.Context) (*pgx.ConnPool, bool) {
	if db, ok := ctx.Value(dbKeyName).(dbValue); ok {
		return db.pool, true
	}
	return nil, false
}

type contextHandler struct {
	pool    *pgx.ConnPool
	handler http.Handler
}

func ContextHandler(pool *pgx.ConnPool, h http.Handler) *contextHandler {
	return &contextHandler{
		pool:    pool,
		handler: h,
	}
}

func (h *contextHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := context.WithValue(req.Context(), dbKeyName, dbValue{h.pool})
	h.handler.ServeHTTP(w, req.WithContext(ctx))
}
