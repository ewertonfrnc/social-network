package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	FollowedUserID int64  `json:"user_id"`     // Who is being followed
	FollowerID     int64  `json:"follower_id"` // Who is following
	CreatedAt      string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (store *FollowerStore) Follow(ctx context.Context, followedID, followerID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO followers (user_id, follower_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`

	_, err := store.db.ExecContext(ctx, query, followedID, followerID)
	return err
}

func (store *FollowerStore) Unfollow(ctx context.Context, followedID, followerID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
	DELETE FROM followers
	WHERE user_id = $1 AND follower_id = $2
	`

	_, err := store.db.ExecContext(ctx, query, followedID, followerID)
	return err
}
