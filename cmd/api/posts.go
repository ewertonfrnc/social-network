package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/ewertonfrnc/social-network/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string

const postCtxKey postKey = "post"

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

	err = app.jsonResponse(w, http.StatusCreated, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	err = app.jsonResponse(w, http.StatusOK, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	err := app.store.Posts.Delete(r.Context(), post.ID)
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

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,min=3,max=50"`
	Content *string `json:"content" validate:"omitempty,min=3,max=500"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	var payload UpdatePostPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	err = app.store.Posts.Update(r.Context(), post)
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

	err = app.jsonResponse(w, http.StatusOK, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// [TODO] Refactor this middleware to be more generic and reusable for other resources (e.g., users, comments)
func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx = context.WithValue(ctx, postCtxKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtxKey).(*store.Post)
	return post
}
