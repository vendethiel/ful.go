package main

/**
 * JSON-related helpers
 */

import (
	"encoding/json"
	"net/http"
)

type JsonError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func MarkAsJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SendJSON(w http.ResponseWriter, o interface{}) error {
	if json, err := json.Marshal(o); err == nil {
		w.Write(json)
		return nil
	} else {
		return err
	}
}

func SendError(w http.ResponseWriter, err error) {
	w.WriteHeader(400)                       // invalid request
	SendJSON(w, JsonError{400, err.Error()}) // XXX check error here?
}
