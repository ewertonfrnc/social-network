package main

import (
	"net/http"

	"github.com/ewertonfrnc/social-network/internal/store"
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
