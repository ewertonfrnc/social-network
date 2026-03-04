package main

import "net/http"

func (app *application) healthCheckHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("OK"))
}
