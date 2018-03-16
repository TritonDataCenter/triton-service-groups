//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package groups_v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/rs/zerolog/log"
)

type ServiceGroup struct {
	ID                  int64  `json:"id"`
	GroupName           string `json:"group_name"`
	TemplateID          int64  `json:"template_id"`
	AccountID           string `json:"account_id"`
	Capacity            int    `json:"capacity"`
	HealthCheckInterval int    `json:"health_check_interval"`
}

func Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	var group *ServiceGroup

	id, err := strconv.Atoi(identifier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, ok := FindGroupByID(ctx, int64(id), session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	writeJsonResponse(w, bytes)
}

func Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var group *ServiceGroup
	err = json.Unmarshal(body, &group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	SaveGroup(ctx, session.AccountID, group)

	err = SubmitOrchestratorJob(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", path.Join(r.URL.Path, group.GroupName))

	com, ok := FindGroupByName(ctx, group.GroupName, session.AccountID)
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

func Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var group *ServiceGroup
	err = json.Unmarshal(body, &group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	UpdateGroup(ctx, identifier, session.AccountID, group)

	err = UpdateOrchestratorJob(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	com, ok := FindGroupByID(ctx, group.ID, session.AccountID)
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

func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	var group *ServiceGroup

	id, err := strconv.Atoi(identifier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, ok := FindGroupByID(ctx, int64(id), session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	RemoveGroup(ctx, group.ID, session.AccountID)

	err = DeleteOrchestratorJob(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	rows, err := FindGroups(ctx, session.AccountID)
	if err != nil {
		log.Fatal().Err(err)
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	writeJsonResponse(w, bytes)
}

type ActionableInput struct {
	InstanceCount int
	MaxInstance   int
	MinInstance   int
}

func Increment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	// Get the Current Group Config
	id, err := strconv.Atoi(identifier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, ok := FindGroupByID(ctx, int64(id), session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	input, err := buildActionableInput(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if group.Capacity >= input.MaxInstance {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Set the new Capacity based on rules
	if group.Capacity+input.InstanceCount > input.MaxInstance {
		group.Capacity = input.MaxInstance
	} else {
		group.Capacity = group.Capacity + input.InstanceCount
	}

	//Update the Database and the orchestration job
	UpdateGroup(ctx, identifier, session.AccountID, group)

	err = UpdateOrchestratorJob(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Return a 202 to suggest accepted
	w.WriteHeader(http.StatusAccepted)
}

func Decrement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	// Get the Current Group Config
	id, err := strconv.Atoi(identifier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, ok := FindGroupByID(ctx, int64(id), session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	input, err := buildActionableInput(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if group.Capacity <= input.MinInstance {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Set the new Capacity based on rules
	if group.Capacity-input.InstanceCount < input.MinInstance {
		group.Capacity = input.MinInstance
	} else {
		group.Capacity = group.Capacity - input.InstanceCount
	}

	//Update the Database and the orchestration job
	UpdateGroup(ctx, identifier, session.AccountID, group)

	err = UpdateOrchestratorJob(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Return a 202 to suggest accepted
	w.WriteHeader(http.StatusAccepted)
}

func writeJsonResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if n, err := w.Write(bytes); err != nil {
		log.Printf("%v", err)
	} else if n != len(bytes) {
		log.Printf("short write: %d/%d", n, len(bytes))
	}
}

func buildActionableInput(r *http.Request) (*ActionableInput, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var input *ActionableInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		return nil, err
	}

	return input, nil
}
