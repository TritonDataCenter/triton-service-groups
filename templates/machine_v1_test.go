package templates_v1_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx"
	tsgRouter "github.com/joyent/triton-service-groups/router"
	"github.com/joyent/triton-service-groups/session"
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

func TestGet(t *testing.T) {
	dbpool, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    dbpool,
	}

	router := tsgRouter.MakeRouter(session)

	req := httptest.NewRequest("GET", "http://example.com/v1/tsg/templates/test-template-1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	testBody := "{\"ID\":319209784155176962,\"Name\":\"test-template-1\",\"Package\":\"test-package\",\"ImageId\":\"test-image-updated\",\"AccountId\":\"joyent\",\"FirewallEnabled\":false,\"Networks\":[\"daeb93a2-532e-4bd4-8788-b6b30f10ac17\"],\"UserData\":\"bash script here\",\"MetaData\":null,\"Tags\":null}"
	assert.Equal(t, testBody, string(body))

}
