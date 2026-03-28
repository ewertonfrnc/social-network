package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ewertonfrnc/social-network/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	rawUserId := chi.URLParam(r, "userId")
	userId, err := strconv.ParseInt(rawUserId, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFound(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	err = app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
