package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/ewertonfrnc/social-network/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtxKey userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	err := app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromContext(r)

	// [TODO]: Revert back to auth userID once we have auth implemented
	var payload struct {
		FollowerID int64 `json:"follower_id"`
	}
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if followedUser.ID == payload.FollowerID {
		app.badRequest(w, r, store.ErrSelfFollow)
		return
	}

	err = app.store.Followers.Follow(r.Context(), followedUser.ID, payload.FollowerID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromContext(r)

	// [TODO]: Revert back to auth userID once we have auth implemented
	var payload struct {
		FollowerID int64 `json:"follower_id"`
	}
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if followedUser.ID == payload.FollowerID {
		app.badRequest(w, r, store.ErrSelfUnfollow)
		return
	}

	err = app.store.Followers.Unfollow(r.Context(), followedUser.ID, payload.FollowerID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFound(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// Middleware
func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawUserId := chi.URLParam(r, "userId")
		userId, err := strconv.ParseInt(rawUserId, 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userId)
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

		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtxKey).(*store.User)
	return user
}
