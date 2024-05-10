// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: likes.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const likePost = `-- name: LikePost :exec

INSERT INTO likes (
    username, post_id
) VALUES (
    $1, $2
)
RETURNING post_id, username
`

type LikePostParams struct {
	Username string
	PostID   uuid.UUID
}

func (q *Queries) LikePost(ctx context.Context, arg LikePostParams) error {
	_, err := q.db.ExecContext(ctx, likePost, arg.Username, arg.PostID)
	return err
}

const unlikePost = `-- name: UnlikePost :exec
DELETE FROM likes 
WHERE username = $1 AND post_id = $2
`

type UnlikePostParams struct {
	Username string
	PostID   uuid.UUID
}

func (q *Queries) UnlikePost(ctx context.Context, arg UnlikePostParams) error {
	_, err := q.db.ExecContext(ctx, unlikePost, arg.Username, arg.PostID)
	return err
}

const updateLikeDecrease = `-- name: UpdateLikeDecrease :exec
UPDATE posts
SET likes_count = likes_count - 1
WHERE id = $1
`

func (q *Queries) UpdateLikeDecrease(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, updateLikeDecrease, id)
	return err
}

const updateLikeIncrease = `-- name: UpdateLikeIncrease :exec
UPDATE posts
SET likes_count = likes_count + 1
WHERE id = $1
`

func (q *Queries) UpdateLikeIncrease(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, updateLikeIncrease, id)
	return err
}
