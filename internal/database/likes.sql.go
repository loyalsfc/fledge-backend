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