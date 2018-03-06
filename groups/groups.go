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

	"strconv"

	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/session"
	"github.com/y0ssar1an/q"
)

type ServiceGroup struct {
	ID                  int64  `json:"id"`
	GroupName           string `json:"group_name"`
	TemplateId          int64  `json:"template_id"`
	AccountId           string `json:"account_id"`
	Capacity            int    `json:"capacity"`
	HealthCheckInterval int    `json:"health_check_interval"`
}

func Get(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		var group *ServiceGroup

		id, err := strconv.Atoi(name)
		if err != nil {
			//At this point we have an actual name so we need to find by name
			g, ok := FindGroupBy(session.DbPool, name, session.AccountId)
			if !ok {
				http.NotFound(w, r)
				return
			}

			group = g
		} else {
			g, ok := FindGroupByID(session.DbPool, int64(id), session.AccountId)
			if !ok {
				http.NotFound(w, r)
				return
			}

			group = g
		}

		bytes, err := json.Marshal(group)
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

		q.Q(session.AccountId)
		var group *ServiceGroup
		err = json.Unmarshal(body, &group)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		SaveGroup(session.DbPool, session.AccountId, group)

		err = SubmitOrchestratorJob(session, group)
		if err != nil {
			panic(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Location", r.URL.Path+"/"+group.GroupName)

		com, ok := FindGroupBy(session.DbPool, group.GroupName, session.AccountId)
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

func Update(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		q.Q("1")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var group *ServiceGroup
		err = json.Unmarshal(body, &group)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		q.Q("2", group)
		UpdateGroup(session.DbPool, name, session.AccountId, group)

		q.Q("3")
		err = UpdateOrchestratorJob(session, group)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		q.Q("4")
		com, ok := FindGroupBy(session.DbPool, group.GroupName, session.AccountId)
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

func Delete(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		var group *ServiceGroup

		id, err := strconv.Atoi(name)
		if err != nil {
			//At this point we have an actual name so we need to find by name
			g, ok := FindGroupBy(session.DbPool, name, session.AccountId)
			if !ok {
				http.NotFound(w, r)
				return
			}

			group = g
		} else {
			g, ok := FindGroupByID(session.DbPool, int64(id), session.AccountId)
			if !ok {
				http.NotFound(w, r)
				return
			}

			group = g
		}

		RemoveGroup(session.DbPool, name, session.AccountId)

		err = DeleteOrchestratorJob(session, group)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
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
