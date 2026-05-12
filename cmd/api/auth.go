package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ewertonfrnc/social-network/internal/mailer"
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

	isProdEnv := app.config.env == "production"
	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// Send Email
	err = app.mailer.Send(mailer.UserWelcomeInviteTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("Error sending Welcome email", "error", err)

		// Rollback user creation if email sending fails
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("Error rolling back user creation after email sending failure", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {}
