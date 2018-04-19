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
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-service-groups/accounts"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ServiceGroup struct {
	ID         string    `json:"id"`
	GroupName  string    `json:"group_name"`
	TemplateID string    `json:"template_id"`
	Capacity   int       `json:"capacity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	uuid := vars["identifier"]

	var group *ServiceGroup

	group, ok := FindGroupByID(ctx, uuid, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}

func Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	group, err := decodeGroupResponseBodyAndValidate(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = SaveGroup(ctx, session.AccountID, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
		return
	}

	writeJSONResponse(w, bytes, http.StatusCreated)
}

func Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	identifier := vars["identifier"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, err := decodeGroupResponseBodyAndValidate(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = UpdateGroup(ctx, identifier, session.AccountID, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = UpdateOrchestratorJob(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	com, ok := FindGroupByID(ctx, identifier, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(com)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, bytes, http.StatusOK)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	uuid := vars["identifier"]

	var group *ServiceGroup

	group, ok := FindGroupByID(ctx, uuid, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	err := RemoveGroup(ctx, group.ID, session.AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := DeleteOrchestratorJob(ctx, group); err != nil {
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

type ActionableInput struct {
	InstanceCount int `json:"instance_count"` //Number of instances to decrement by
	MaxInstance   int `json:"max_instance"`   //Maximum number of instances allowed in group
	MinInstance   int `json:"min_instance"`   //Minimum number of instances allowed in group
}

func Increment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	uuid := vars["identifier"]

	// Get the Current Group Config
	group, ok := FindGroupByID(ctx, uuid, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	input, err := buildActionableInput(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
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
	err = UpdateGroup(ctx, uuid, session.AccountID, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	uuid := vars["identifier"]

	// Get the Current Group Config
	group, ok := FindGroupByID(ctx, uuid, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	input, err := buildActionableInput(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
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
	err = UpdateGroup(ctx, uuid, session.AccountID, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := UpdateOrchestratorJob(ctx, group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Return a 202 to suggest accepted
	w.WriteHeader(http.StatusAccepted)
}

func ListInstances(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := handlers.GetAuthSession(ctx)

	vars := mux.Vars(r)
	uuid := vars["identifier"]

	group, ok := FindGroupByID(ctx, uuid, session.AccountID)
	if !ok {
		http.NotFound(w, r)
		return
	}

	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		http.Error(w, handlers.ErrNoConnPool.Error(), http.StatusInternalServerError)
		return
	}
	store := accounts.NewStore(db)
	account, err := store.FindByID(ctx, session.AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	credential, err := account.GetTritonCredential(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input := authentication.PrivateKeySignerInput{
		KeyID:              credential.KeyID,
		PrivateKeyMaterial: []byte(credential.KeyMaterial),
		AccountName:        credential.AccountName,
	}
	signer, err := authentication.NewPrivateKeySigner(input)
	if err != nil {
		returnError := errors.Wrapf(err, "error Creating SSH Private Key Signer")
		http.Error(w, returnError.Error(), http.StatusInternalServerError)
		return
	}

	config := &triton.ClientConfig{
		TritonURL:   session.TritonURL,
		AccountName: credential.AccountName,
		Signers:     []authentication.Signer{signer},
	}

	c, err := compute.NewClient(config)
	if err != nil {
		returnError := errors.Wrapf(err, "error constructing ComputeClient")
		http.Error(w, returnError.Error(), http.StatusInternalServerError)
		return
	}

	params := &compute.ListInstancesInput{}
	t := make(map[string]interface{}, 0)
	t["tsg.name"] = group.GroupName
	params.Tags = t

	instances, err := c.Instances().List(ctx, params)
	if err != nil {
		returnError := errors.Wrapf(err, "error listing instances in TSG")
		http.Error(w, returnError.Error(), http.StatusInternalServerError)
		return
	}

	if instances != nil && len(instances) == 0 {
		writeJSONResponse(w, []byte("[]"), http.StatusOK)
		return
	}

	bytes, err := json.Marshal(instances)
	if err != nil {
		returnError := errors.Wrapf(err, "error marshalling TSG instance list")
		http.Error(w, returnError.Error(), http.StatusInternalServerError)
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

	err = input.Validate()
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (i *ActionableInput) Validate() error {

	if i.MinInstance < 0 || i.MaxInstance < 0 || i.InstanceCount < 0 {
		return errors.New("only positive integers are allowed for instance count & max and min range")
	}

	return nil
}

func decodeGroupResponseBodyAndValidate(body []byte) (*ServiceGroup, error) {
	var group *ServiceGroup
	err := json.Unmarshal(body, &group)
	if err != nil {
		return nil, errors.New("error in unmarshal request body")
	}

	if len(group.GroupName) > 182 {
		return nil, errors.New("group name cannot be more than 182 characters")
	}

	if !isValidUUID(group.TemplateID) {
		return nil, errors.New("templateID must be a valid UUID")
	}

	if group.Capacity < 0 {
		return nil, errors.New("group capacity cannot be a negative number")
	}

	if group.Capacity > 100 {
		return nil, errors.New("group capacity cannot be more than 100 compute instances")
	}

	return group, nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
