package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ewertonfrnc/social-network/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq, err := parsePaginatedFeedQuery(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(60), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusOK, feed)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func parsePaginatedFeedQuery(r *http.Request) (store.PaginatedFeedQuery, error) {
	fq := store.NewPaginatedFeedQuery()

	queryString := r.URL.Query()

	if limitString := queryString.Get("limit"); limitString != "" {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			return fq, err
		}
		fq.Limit = limit
	}

	if offsetString := queryString.Get("offset"); offsetString != "" {
		offset, err := strconv.Atoi(offsetString)
		if err != nil {
			return fq, err
		}
		fq.Offset = offset
	}

	if sortDirection := queryString.Get("sort"); sortDirection != "" {
		fq.SortDirection = sortDirection
	}

	if tags := queryString.Get("tags"); tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	if search := queryString.Get("search"); search != "" {
		fq.Search = search
	}

	return fq, nil
}
