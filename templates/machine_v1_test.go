package templates_v1_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx"
	"github.com/joyent/triton-service-groups/session"
	templates_v1 "github.com/joyent/triton-service-groups/templates"
	"github.com/stretchr/testify/assert"
)

// TODO: We should refactor how/where our database initializes so we can half
// bootstrap the application from our tests with a simple one-liner.
func initDB() (*pgx.ConnPool, error) {
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		MaxConnections: 5,
		AfterConnect:   nil,
		AcquireTimeout: 0,
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Database: "triton_test",
			Port:     26257,
			User:     "root",
		},
	})
	if err != nil {
		return nil, err
	}

	return connPool, nil
}

func TestList(t *testing.T) {
	dbpool, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	session := &session.TsgSession{
		DbPool: dbpool,
	}

	req := httptest.NewRequest("GET", "http://example.com/tsg/v1/templates", nil)
	recorder := httptest.NewRecorder()
	templates_v1.List(session)(recorder, req)

	resp := recorder.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, "null", string(body))
}
