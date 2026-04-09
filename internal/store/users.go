package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}

type Password struct {
	Hash []byte `json:"-"`
}

type UserStore struct {
	db *sql.DB
}

func (p *Password) SetPassword(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Hash = hash

	return nil
}

func (store *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, username, email, created_at
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.Hash,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key" (23505)`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key" (23505)`:
			return ErrDuplicateUsername
		}

		return err
	}

	return nil
}

func (store *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT id, email, username
	FROM users
	WHERE id = $1
	`

	user := &User{}
	err := store.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (store *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, expiresAt time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return withTx(store.db, ctx, func(tx *sql.Tx) error {
		if err := store.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := store.createUserInvite(ctx, tx, user.ID, token, expiresAt); err != nil {
			return err
		}

		return nil
	})
}

func (store *UserStore) Activate(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return withTx(store.db, ctx, func(tx *sql.Tx) error {
		user, err := store.GetUserFromInvite(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true
		if err := store.updateUser(ctx, tx, user); err != nil {
			return err
		}

		if err := store.deleteUserInvite(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (store *UserStore) GetUserFromInvite(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		u.id,
		u.email,
		u.username,
		u.created_at,
		u.is_active
	FROM
		users u
		JOIN user_invitations ui ON ui.user_id = u.id
	WHERE
		ui.token = $1
		AND ui.expires_at > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (store *UserStore) createUserInvite(ctx context.Context, tx *sql.Tx, userID int64, token string, expiresAt time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO user_invitations (token, user_id, expires_at)
	VALUES ($1, $2, $3)
	`

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(expiresAt))
	if err != nil {
		return err
	}

	return nil
}

func (store *UserStore) updateUser(ctx context.Context, tx *sql.Tx, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	UPDATE users
	SET
		username = $1,
		email = $2,
		is_active = $3
	WHERE
		id = $4
	`

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (store *UserStore) deleteUserInvite(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	DELETE FROM user_invitations
	WHERE
		user_id = $1
	`

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
