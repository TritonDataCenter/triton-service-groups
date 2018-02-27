package templates_v1_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bytes"

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

func TestAcc_Get(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

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

	testBody := "{\"ID\":319209784155176962,\"TemplateName\":\"test-template-1\",\"AccountId\":\"joyent\",\"Package\":\"test-package\",\"ImageId\":\"49b22aec-0c8a-11e6-8807-a3eb4db576ba\",\"MachineNamePrefix\":\"sample-\",\"FirewallEnabled\":false,\"Networks\":[\"f7ed95d3-faaf-43ef-9346-15644403b963\"],\"UserData\":\"bash script here\",\"MetaData\":null,\"Tags\":null}"
	assert.Equal(t, testBody, string(body))
}

func TestAcc_GetIncorrectTemplateName(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	dbpool, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    dbpool,
	}

	router := tsgRouter.MakeRouter(session)

	req := httptest.NewRequest("GET", "http://example.com/v1/tsg/templates/test-template-200", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	_, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAcc_List(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	dbpool, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    dbpool,
	}

	router := tsgRouter.MakeRouter(session)

	req := httptest.NewRequest("GET", "http://example.com/v1/tsg/templates", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	if string(body) == "" {
		t.Fatal()
	}
}

func TestAcc_Delete(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	dbpool, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    dbpool,
	}

	router := tsgRouter.MakeRouter(session)

	req := httptest.NewRequest("DELETE", "http://example.com/v1/tsg/templates/test-template-6", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	if string(body) != "" {
		t.Fatal()
	}
}

func TestAcc_DeleteNonExistantTemplate(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	dbpool, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    dbpool,
	}

	router := tsgRouter.MakeRouter(session)

	req := httptest.NewRequest("DELETE", "http://example.com/v1/tsg/templates/test-template-200", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	_, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAcc_CreateTemplate(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	dbpool, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    dbpool,
	}

	testBody := `{
	"TemplateName": "test-template-7",
		"AccountId": "joyent",
		"Package": "test-package",
		"ImageId": "49b22aec-0c8a-11e6-8807-a3eb4db576ba",
		"MachineNamePrefix": "sample-",
		"FirewallEnabled": false,
		"Networks": [
	"f7ed95d3-faaf-43ef-9346-15644403b963"
	],
	"UserData": "bash script here",
		"Tags": {
	"foo": "bar",
	"owner": "stack72"
	},
	"MetaData": null
}`

	r := bytes.NewReader([]byte(testBody))
	router := tsgRouter.MakeRouter(session)
	req := httptest.NewRequest("POST", "http://example.com/v1/tsg/templates", r)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	_, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
