package templates_v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type MachineTemplate struct {
	Name            string
	Package         string
	ImageID         string
	FirewallEnabled bool
	Networks        []string
	UserData        string
	MetaData        map[string]interface{}
	Tags            map[string]string
}

func Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	com, ok := FindTemplateBy(name)
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

func Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var template *MachineTemplate
	err = json.Unmarshal(body, &template)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	SaveTemplate(template)

	w.Header().Set("Location", r.URL.Path+"/"+template.Name)
	w.WriteHeader(http.StatusCreated)
}

func Update(w http.ResponseWriter, r *http.Request) {
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

	UpdateTemplate(name, template)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	_, ok := FindTemplateBy(name)
	if !ok {
		http.NotFound(w, r)
		return
	}

	RemoveTemplate(name)
	w.WriteHeader(http.StatusNoContent)
}

func List(w http.ResponseWriter, r *http.Request) {
	rows, err := FindTemplates()
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

func writeJsonResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if n, err := w.Write(bytes); err != nil {
		log.Printf("%v", err)
	} else if n != len(bytes) {
		log.Printf("short write: %d/%d", n, len(bytes))
	}
}
