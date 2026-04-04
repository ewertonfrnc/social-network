package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (store *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	sortDirection := "DESC"
	if strings.EqualFold(fq.SortDirection, "asc") {
		sortDirection = "ASC"
	}

	query := fmt.Sprintf(`
	SELECT
		p.id,
		p.title,
		p.content,
		p.created_at,
		p.tags,
		p.version,
		u.username,
		COUNT(c.id) AS comments_count
	FROM
		posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN followers f ON f.follower_id = p.user_id
		OR p.user_id = $1
	WHERE
		f.user_id = $1
		AND (
			p.title ILIKE '%%' || $4 || '%%'
			OR p.content ILIKE '%%' || $4 || '%%'
		)
		AND (
			p.tags && $5
			OR $5 IS NULL
		)
	GROUP BY
		p.id,
		u.username
	ORDER BY
		p.created_at %s,
		p.id %s
	LIMIT
		$2
	OFFSET
		$3
	`, sortDirection, sortDirection)

	rows, err := store.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []PostWithMetadata{}
	for rows.Next() {
		var post PostWithMetadata

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			pq.Array(&post.Tags),
			&post.Version,
			&post.User.Username,
			&post.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (store *PostStore) Create(ctx context.Context, post *Post) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO posts (title, content, tags, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, title, content, tags, created_at, updated_at, version
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
		&post.Version,
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostStore) GetByID(ctx context.Context, postId int64) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		SELECT id, title, content, tags, created_at, updated_at, user_id, version
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
		&post.Version,
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

func (store *PostStore) Delete(ctx context.Context, postId int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	DELETE FROM posts
	WHERE id = $1
	`
	result, err := store.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (store *PostStore) Update(ctx context.Context, post *Post) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	UPDATE posts
	SET title = $1, content = $2, updated_at = NOW(), version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING id, title, content, tags, created_at, updated_at, version
	`

	err := store.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
