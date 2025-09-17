package main

import (
	"errors"
	"net/http"
)

func (me *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": "development",
			"version":     "1.0.0",
		},
	}

	err := me.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		me.serverErrorResponse(w, r, err)
		return
	}
}

func (me *application) searchByTag(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	tag := queryValues.Get("tag")
	if tag == "" {
		me.badRequestResponse(w, r, errors.New("tag params is required"))
	}

	me.db.Exec("")

}
