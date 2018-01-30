package groups_v1

import (
	"net/http"

	"github.com/jackc/pgx"
)

type ServiceGroup struct {
	GroupName           string
	TemplateName        string
	Capacity            int
	DataCenter          []string
	HealthCheckInterval int //default will be 300
	InstanceTags        map[string]interface{}
}

func Get(dbPool *pgx.ConnPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func Create(dbPool *pgx.ConnPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func Update(dbPool *pgx.ConnPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func Delete(dbPool *pgx.ConnPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func List(dbPool *pgx.ConnPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
