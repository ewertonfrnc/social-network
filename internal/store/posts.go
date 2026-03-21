package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	UserID    int64     `json:"user_id"`
	Comments  []Comment `json:"comments"`
}

type PostStore struct {
	db *sql.DB
}

func (store *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (title, content, tags, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, title, content, tags, created_at, updated_at
	`

	err := store.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(&post.Tags),
		post.UserID,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostStore) GetByID(ctx context.Context, postId int64) (*Post, error) {
	query := `
		SELECT id, title, content, tags, created_at, updated_at, user_id
		FROM posts
		WHERE id = $1
		`

	post := &Post{}

	err := store.db.QueryRowContext(ctx, query, postId).Scan(&post.ID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}
