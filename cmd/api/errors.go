package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal Server Error", "Method", r.Method, "Path", r.URL.Path, "Error", err)

	WriteJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Bad Request", "Method", r.Method, "Path", r.URL.Path, "Error", err)
	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Not Found", "Method", r.Method, "Path", r.URL.Path, "Error", err)
	WriteJSONError(w, http.StatusNotFound, "Not Found")
}
