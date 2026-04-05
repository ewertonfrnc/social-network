package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("Internal Server Error", "Method", r.Method, "Path", r.URL.Path, "Error", err.Error())

	WriteJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Bad Request", "Method", r.Method, "Path", r.URL.Path, "Error", err.Error())
	WriteJSONError(w, http.StatusBadRequest, "Bad Request")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Not Found", "Method", r.Method, "Path", r.URL.Path, "Error", err.Error())
	WriteJSONError(w, http.StatusNotFound, "Not Found")
}
