//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package templates_v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"errors"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/rs/zerolog/log"
)

type InstanceTemplate struct {
	ID              string            `json:"id"`
	TemplateName    string            `json:"template_name"`
	Package         string            `json:"package"`
	ImageID         string            `json:"image_id"`
	FirewallEnabled bool              `json:"firewall_enabled"`
	Networks        []string          `json:"networks"`
	UserData        string            `json:"userdata"`
	MetaData        map[string]string `json:"metadata"`
	Tags            map[string]string `json:"tags"`
	CreatedAt       time.Time         `json:"created_at"`
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}

func Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	template, err := decodeResponseBodyAndValidate(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = SaveTemplate(ctx, session.AccountID, template)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	com, ok := FindTemplateByName(ctx, template.TemplateName, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(com)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", path.Join(r.URL.Path, template.TemplateName))
	writeJSONResponse(w, bytes, http.StatusCreated)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	uuid := vars["identifier"]

	var template *InstanceTemplate

	templateAllocated, err := CheckTemplateAllocationByID(ctx, uuid, session.AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if templateAllocated {
		http.Error(w, fmt.Sprintf("Cannot delete template %q while in use, "+
			"must be removed from all groups first.", uuid),
			http.StatusConflict)
		return
	}

	template, ok := FindTemplateByID(ctx, uuid, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	err = RemoveTemplate(ctx, template.ID, session.AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	rows, err := FindTemplates(ctx, session.AccountID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if rows == nil || len(rows) == 0 {
		writeJSONResponse(w, []byte("[]"), http.StatusOK)
		return
	}

	bytes, err := json.Marshal(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}

func writeJSONResponse(w http.ResponseWriter, bytes []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if n, err := w.Write(bytes); err != nil {
		log.Printf("%v", err)
	} else if n != len(bytes) {
		log.Printf("short write: %d/%d", n, len(bytes))
	}
}

func decodeResponseBodyAndValidate(body []byte) (*InstanceTemplate, error) {
	var template *InstanceTemplate
	err := json.Unmarshal(body, &template)
	if err != nil {
		return nil, errors.New("error in unmarshal request body")
	}

	if !isValidUUID(template.Package) {
		return nil, errors.New("package must be a valid UUID")
	}

	if !isValidUUID(template.ImageID) {
		return nil, errors.New("imageID must be a valid UUID")
	}

	return template, nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
