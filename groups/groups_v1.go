//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package groups_v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/session"
)

type ServiceGroup struct {
	ID                  int64
	GroupName           string
	TemplateId          int64
	AccountId           string
	Capacity            int
	HealthCheckInterval int
}

func Get(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		com, ok := FindGroupBy(session.DbPool, name, session.AccountId)
		if !ok {
			http.NotFound(w, r)
			return
		}

		bytes, err := json.Marshal(com)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		writeJsonResponse(w, bytes)
	}
}

func Create(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var group *ServiceGroup
		err = json.Unmarshal(body, &group)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		SaveGroup(session.DbPool, session.AccountId, group)

		w.Header().Set("Location", r.URL.Path+"/"+group.GroupName)
		w.WriteHeader(http.StatusCreated)
	}
}

func Update(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var group *ServiceGroup
		err = json.Unmarshal(body, &group)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		UpdateGroup(session.DbPool, name, session.AccountId, group)
	}
}

func Delete(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		_, ok := FindGroupBy(session.DbPool, name, session.AccountId)
		if !ok {
			http.NotFound(w, r)
			return
		}

		RemoveGroup(session.DbPool, name, session.AccountId)
		w.WriteHeader(http.StatusGone)
	}
}

func List(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := FindGroups(session.DbPool, session.AccountId)
		if err != nil {
			log.Fatal(err)
			http.NotFound(w, r)
			return
		}

		bytes, err := json.Marshal(rows)
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
