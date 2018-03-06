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

	"strconv"

	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/session"
)

type InstanceTemplate struct {
	ID                 int64             `json:"id"`
	TemplateName       string            `json:"template_name"`
	AccountId          string            `json:"account_id"`
	Package            string            `json:"package"`
	ImageId            string            `json:"image_id"`
	InstanceNamePrefix string            `json:"instance_name_prefix"`
	FirewallEnabled    bool              `json:"firewall_enabled"`
	Networks           []string          `json:"networks"`
	UserData           string            `json:"userdata"`
	MetaData           map[string]string `json:"metadata"`
	Tags               map[string]string `json:"tags"`
}

func Get(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		var template *InstanceTemplate

		id, err := strconv.Atoi(name)
		if err != nil {
			//At this point we have an actual name so we need to find by name
			t, ok := FindTemplateByName(session.DbPool, name, session.AccountId)
			if !ok {
				http.NotFound(w, r)
				return
			}

			template = t
		} else {
			t, ok := FindTemplateByID(session.DbPool, int64(id), session.AccountId)
			if !ok {
				http.NotFound(w, r)
				return
			}

			template = t
		}

		bytes, err := json.Marshal(template)
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

		var template *InstanceTemplate
		err = json.Unmarshal(body, &template)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}

		SaveTemplate(session.DbPool, session.AccountId, template)

		w.Header().Set("Location", r.URL.Path+"/"+template.TemplateName)

		com, ok := FindTemplateByName(session.DbPool, template.TemplateName, session.AccountId)
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

		_, ok := FindTemplateByName(session.DbPool, name, session.AccountId)
		if !ok {
			http.NotFound(w, r)
			return
		}

		RemoveTemplate(session.DbPool, name, session.AccountId)
		w.WriteHeader(http.StatusNoContent)
	}
}

func List(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := FindTemplates(session.DbPool, session.AccountId)
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
