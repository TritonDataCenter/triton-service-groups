package templates_v1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	tsgRouter "github.com/joyent/triton-service-groups/router"
	"github.com/joyent/triton-service-groups/session"
	"github.com/joyent/triton-service-groups/templates"
	"github.com/joyent/triton-service-groups/testutils"
	"github.com/stretchr/testify/assert"
)

func TestAcc_Get(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	db, err := testutils.NewTestDB()
	if err != nil {
		t.Error(err)
	}
	db.Clear(t)

	tmplName := "test-template-1"
	testTmpl := &templates_v1.InstanceTemplate{
		TemplateName:       tmplName,
		Package:            "test-package",
		ImageId:            "123456",
		InstanceNamePrefix: "sample-",
		FirewallEnabled:    false,
		Networks:           []string{"123456"},
		UserData:           "bash script here",
		MetaData:           nil,
		Tags:               nil,
	}
	testTmpl.Save(db.Conn, "joyent")

	tmpl, ok := templates_v1.FindByName(tmplName, db.Conn, "joyent")
	if !ok {
		t.Error("failed to find test template")
	}

	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    db.Conn,
	}

	router := tsgRouter.MakeRouter(session)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://example.com/v1/tsg/templates/%s", tmplName), nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	var respTmpl *templates_v1.InstanceTemplate

	if err := json.Unmarshal(body, &respTmpl); err != nil {
		t.Error(err)
	}

	assert.Equal(t, tmpl.TemplateName, respTmpl.TemplateName)
	assert.Equal(t, tmpl.Package, respTmpl.Package)
	assert.Equal(t, tmpl.ImageId, respTmpl.ImageId)
	assert.Equal(t, tmpl.InstanceNamePrefix, respTmpl.InstanceNamePrefix)
	assert.Equal(t, tmpl.FirewallEnabled, respTmpl.FirewallEnabled)
	assert.Equal(t, tmpl.Networks, respTmpl.Networks)
	assert.Equal(t, tmpl.UserData, respTmpl.UserData)
	assert.Equal(t, tmpl.MetaData, respTmpl.MetaData)
	assert.Equal(t, tmpl.Tags, respTmpl.Tags)
}

func TestAcc_GetIncorrectTemplateName(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	db, err := testutils.NewTestDB()
	if err != nil {
		t.Error(err)
	}
	db.Clear(t)

	tmplName := "test-template-1"
	testTmpl := &templates_v1.InstanceTemplate{
		TemplateName:       tmplName,
		Package:            "test-package",
		ImageId:            "123456",
		InstanceNamePrefix: "sample-",
		FirewallEnabled:    false,
		Networks:           []string{"123456"},
		UserData:           "bash script here",
		MetaData:           nil,
		Tags:               nil,
	}
	testTmpl.Save(db.Conn, "joyent")

	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    db.Conn,
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

	db, err := testutils.NewTestDB()
	if err != nil {
		t.Error(err)
	}
	db.Clear(t)

	names := []string{"test-template-1", "another-template-2"}
	tmpls := make([]*templates_v1.InstanceTemplate, len(names))
	for n, name := range names {
		testTmpl := &templates_v1.InstanceTemplate{
			TemplateName:       name,
			Package:            fmt.Sprintf("test-package-%d", n),
			ImageId:            fmt.Sprintf("12345%d", n),
			InstanceNamePrefix: "sample-",
			FirewallEnabled:    false,
			Networks:           []string{"123456"},
			UserData:           "bash script here",
			MetaData:           nil,
			Tags:               nil,
		}
		testTmpl.Save(db.Conn, "joyent")
		tmpls[n] = testTmpl
	}

	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    db.Conn,
	}

	router := tsgRouter.MakeRouter(session)

	req := httptest.NewRequest("GET", "http://example.com/v1/tsg/templates", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	var respTmpls []templates_v1.InstanceTemplate
	if err := json.Unmarshal(body, &respTmpls); err != nil {
		t.Error(err)
	}

	assert.Equal(t, 2, len(respTmpls))

	for n, respTmpl := range respTmpls {
		tmpl := tmpls[n]
		assert.Equal(t, tmpl.TemplateName, respTmpl.TemplateName)
		assert.Equal(t, tmpl.Package, respTmpl.Package)
		assert.Equal(t, tmpl.ImageId, respTmpl.ImageId)
		assert.Equal(t, tmpl.InstanceNamePrefix, respTmpl.InstanceNamePrefix)
		assert.Equal(t, tmpl.FirewallEnabled, respTmpl.FirewallEnabled)
		assert.Equal(t, tmpl.Networks, respTmpl.Networks)
		assert.Equal(t, tmpl.UserData, respTmpl.UserData)
		assert.Equal(t, tmpl.MetaData, respTmpl.MetaData)
		assert.Equal(t, tmpl.Tags, respTmpl.Tags)
	}
}

func TestAcc_Delete(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	db, err := testutils.NewTestDB()
	if err != nil {
		t.Error(err)
	}
	db.Clear(t)

	tmplName := "test-template-1"

	testTmpl := &templates_v1.InstanceTemplate{
		TemplateName: tmplName,
	}
	testTmpl.Save(db.Conn, "joyent")

	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    db.Conn,
	}

	router := tsgRouter.MakeRouter(session)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("http://example.com/v1/tsg/templates/%s", tmplName), nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusGone, resp.StatusCode)

	if string(body) != "" {
		t.Fatal()
	}
}

func TestAcc_DeleteNonExistantTemplate(t *testing.T) {
	if os.Getenv("TRITON_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	db, err := testutils.NewTestDB()
	if err != nil {
		t.Error(err)
	}
	db.Clear(t)

	tmplName := "test-template-1"

	testTmpl := &templates_v1.InstanceTemplate{
		TemplateName: tmplName,
	}
	testTmpl.Save(db.Conn, "joyent")

	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    db.Conn,
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

	db, err := testutils.NewTestDB()
	if err != nil {
		t.Error(err)
	}
	db.Clear(t)

	session := &session.TsgSession{
		AccountId: "joyent",
		DbPool:    db.Conn,
	}

	tmplName := "test-template-7"
	tmpl := &templates_v1.InstanceTemplate{
		TemplateName:       tmplName,
		Package:            "test-package",
		ImageId:            "123456",
		InstanceNamePrefix: "sample-",
		FirewallEnabled:    false,
		Networks:           []string{"123456"},
		UserData:           "bash script here",
		MetaData:           nil,
		Tags:               nil,
	}
	jsonTmpl, err := json.Marshal(tmpl)
	if err != nil {
		t.Error(err)
	}

	r := bytes.NewReader(jsonTmpl)
	router := tsgRouter.MakeRouter(session)
	req := httptest.NewRequest("POST", "http://example.com/v1/tsg/templates", r)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	resp := recorder.Result()
	_, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respTmpl, ok := templates_v1.FindByName("test-template-7", db.Conn, "joyent")
	if !ok {
		t.Error("failed to find created template by name test-template-7")
	}

	assert.Equal(t, tmpl.TemplateName, respTmpl.TemplateName)
	assert.Equal(t, tmpl.Package, respTmpl.Package)
	assert.Equal(t, tmpl.ImageId, respTmpl.ImageId)
	assert.Equal(t, tmpl.InstanceNamePrefix, respTmpl.InstanceNamePrefix)
	assert.Equal(t, tmpl.FirewallEnabled, respTmpl.FirewallEnabled)
	assert.Equal(t, tmpl.Networks, respTmpl.Networks)
	assert.Equal(t, tmpl.UserData, respTmpl.UserData)
	assert.Equal(t, tmpl.MetaData, respTmpl.MetaData)
	assert.Equal(t, tmpl.Tags, respTmpl.Tags)
}
