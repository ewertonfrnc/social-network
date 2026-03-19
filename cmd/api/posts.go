package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ewertonfrnc/social-network/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,min=3,max=50"`
	Content string   `json:"content" validate:"required,min=3,max=500"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	ctx := r.Context()

	err = app.store.Posts.Create(ctx, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = WriteJSON(w, http.StatusCreated, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	rawPostId := chi.URLParam(r, "postId")
	postId, err := strconv.ParseInt(rawPostId, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, postId)
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

	err = WriteJSON(w, http.StatusOK, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
