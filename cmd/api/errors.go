package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[ERROR]: Internal Server Error. Method: %v. Path: %v. Error: %v", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[ERROR]: Bad Request. Method: %v. Path: %v. Error: %v", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusBadRequest, "Bad Request")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[ERROR]: Not Found. Method: %v. Path: %v. Error: %v", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusNotFound, "Not Found")
}
