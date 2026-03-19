package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := WriteJSON(w, http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
