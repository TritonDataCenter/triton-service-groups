package handlers

import (
	"context"
	"net/http"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/jackc/pgx"
)

type contextKey int

const (
	dbKeyName contextKey = iota
	authKey
	nomadKeyName
)

type dbValue struct {
	pool *pgx.ConnPool
}

type nomadValue struct {
	client *nomad.Client
}

// GetDBPool pulls a configured database client out of the current request
// context.
func GetDBPool(ctx context.Context) (*pgx.ConnPool, bool) {
	if db, ok := ctx.Value(dbKeyName).(dbValue); ok {
		return db.pool, true
	}
	return nil, false
}

// GetNomadClient pulls a configured nomad client out of the current request
// context.
func GetNomadClient(ctx context.Context) (*nomad.Client, bool) {
	if nomad, ok := ctx.Value(nomadKeyName).(nomadValue); ok {
		return nomad.client, true
	}
	return nil, false
}

type contextHandler struct {
	pool    *pgx.ConnPool
	nomad   *nomad.Client
	handler http.Handler
}

func ContextHandler(pool *pgx.ConnPool, nomad *nomad.Client, h http.Handler) *contextHandler {
	return &contextHandler{
		pool:    pool,
		nomad:   nomad,
		handler: h,
	}
}

func (h *contextHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := context.WithValue(req.Context(), dbKeyName, dbValue{h.pool})
	h.handler.ServeHTTP(w, req.WithContext(ctx))
}
