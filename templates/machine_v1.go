package templates_v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joyent/triton-service-groups/session"
)

type MachineTemplate struct {
	Name            string
	Package         string
	ImageID         string
	AccountName     string
	FirewallEnabled bool
	Networks        []string
	UserData        string
	MetaData        map[string]interface{}
	Tags            map[string]string
}

func Get(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		com, ok := FindTemplateBy(session.DbPool, name, session.AccountId)
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

		var template *MachineTemplate
		err = json.Unmarshal(body, &template)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		SaveTemplate(session.DbPool, session.AccountId, template)

		w.Header().Set("Location", r.URL.Path+"/"+template.Name)
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

		var template *MachineTemplate
		err = json.Unmarshal(body, &template)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		UpdateTemplate(session.DbPool, name, session.AccountId, template)
	}
}

func Delete(session *session.TsgSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		_, ok := FindTemplateBy(session.DbPool, name, session.AccountId)
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
		log.Printf("Initializing route")
		log.Printf("Account Id in HTTP Request %v", session.AccountId)

		rows, err := FindTemplates(session.DbPool, session.AccountId)
		if err != nil {
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
