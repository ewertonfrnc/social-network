package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ewertonfrnc/social-network/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	err := ReadJSON(w, r, &payload)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
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
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = WriteJSON(w, http.StatusCreated, post)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Failed to write response")
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	rawPostId := chi.URLParam(r, "postId")
	postId, err := strconv.ParseInt(rawPostId, 10, 64)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			WriteJSONError(w, http.StatusNotFound, "Post not found")
			return
		default:
			WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = WriteJSON(w, http.StatusOK, post)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Failed to write GET response")
		return
	}
}
