package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/ewertonfrnc/social-network/internal/store"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := ReadJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.SetPassword(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.expiresAt)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequest(w, r, store.ErrDuplicateEmail)
			return
		case store.ErrDuplicateUsername:
			app.badRequest(w, r, store.ErrDuplicateUsername)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	userWithToken := &UserWithToken{
		User:  user,
		Token: plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {}
