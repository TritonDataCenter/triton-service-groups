//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package templates_v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/session"
)

type InstanceTemplate struct {
	ID                 int64
	TemplateName       string
	AccountId          string
	Package            string
	ImageId            string
	InstanceNamePrefix string
	FirewallEnabled    bool
	Networks           []string
	UserData           string
	MetaData           map[string]string
	Tags               map[string]string
}

func Get(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		t, ok := FindByName(name, session.DbPool, session.AccountId)
		if !ok {
			http.NotFound(w, r)
			return
		}

		jsonBytes, err := json.Marshal(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		writeJsonResponse(w, jsonBytes)
	}
}

func Create(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var t *InstanceTemplate
		err = json.Unmarshal(body, &t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		t.Save(session.DbPool, session.AccountId)

		w.Header().Set("Location", path.Join(r.URL.Path, t.TemplateName))
		w.WriteHeader(http.StatusCreated)
	}
}

func Delete(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		t, ok := FindByName(name, session.DbPool, session.AccountId)
		if !ok {
			http.NotFound(w, r)
			return
		}

		t.Archive(session.DbPool, session.AccountId)
		w.WriteHeader(http.StatusGone)
	}
}

func List(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates, err := FindAll(session.DbPool, session.AccountId)
		if err != nil {
			log.Fatal(err)
			http.NotFound(w, r)
			return
		}

		bytes, err := json.Marshal(templates)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		writeJsonResponse(w, bytes)
	}
}

func writeJsonResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if n, err := w.Write(bytes); err != nil {
		log.Printf("%v", err)
	} else if n != len(bytes) {
		log.Printf("short write: %d/%d", n, len(bytes))
	}
}
