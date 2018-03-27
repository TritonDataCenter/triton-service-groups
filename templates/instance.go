//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package templates_v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/rs/zerolog/log"
)

type InstanceTemplate struct {
	ID                 string            `json:"id"`
	TemplateName       string            `json:"template_name"`
	AccountID          string            `json:"account_id"`
	Package            string            `json:"package"`
	ImageID            string            `json:"image_id"`
	InstanceNamePrefix string            `json:"instance_name_prefix"`
	FirewallEnabled    bool              `json:"firewall_enabled"`
	Networks           []string          `json:"networks"`
	UserData           string            `json:"userdata"`
	MetaData           map[string]string `json:"metadata"`
	Tags               map[string]string `json:"tags"`
}

func (t *InstanceTemplate) ShortID() string {
	if t.ID == "" {
		return ""
	}
	return t.ID[:8]
}

func Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	uuid := vars["identifier"]

	var template *InstanceTemplate

	template, ok := FindTemplateByID(ctx, uuid, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(template)
	if err != nil {
		log.Fatal().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	writeJsonResponse(w, bytes)
}

func Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var template *InstanceTemplate
	err = json.Unmarshal(body, &template)
	if err != nil {
		log.Fatal().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	SaveTemplate(ctx, session.AccountID, template)

	w.Header().Set("Location", path.Join(r.URL.Path, template.TemplateName))

	com, ok := FindTemplateByName(ctx, template.TemplateName, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(com)
	if err != nil {
		log.Fatal().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)

	writeJsonResponse(w, bytes)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	uuid := vars["identifier"]

	var template *InstanceTemplate

	template, ok := FindTemplateByID(ctx, uuid, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	RemoveTemplate(ctx, template.ID, session.AccountID)
	w.WriteHeader(http.StatusNoContent)
}

func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	rows, err := FindTemplates(ctx, session.AccountID)
	if err != nil {
		log.Fatal().Err(err)
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(rows)
	if err != nil {
		log.Fatal().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	writeJsonResponse(w, bytes)
}

func writeJsonResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if n, err := w.Write(bytes); err != nil {
		log.Printf("%v", err)
	} else if n != len(bytes) {
		log.Printf("short write: %d/%d", n, len(bytes))
	}
}
