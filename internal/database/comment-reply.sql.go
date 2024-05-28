// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: comment-reply.sql

package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const deleteReply = `-- name: DeleteReply :exec
DELETE FROM comment_reply
WHERE id = $1
RETURNING id, reply_text, media, username, comment_id, likes_count, created_at, updated_at
`

func (q *Queries) DeleteReply(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteReply, id)
	return err
}

const getReplies = `-- name: GetReplies :many
SELECT c.id, c.reply_text, c.media, c.username, c.comment_id, c.likes_count, c.created_at, c.updated_at, u.name, u.profile_picture, u.is_verified,
    array_agg(l.username) AS liked_users_username
FROM comment_reply c
INNER JOIN users u On u.username = c.username
LEFT JOIN like_reply l ON l.reply_id =c.id
WHERE c.comment_id = $1
GROUP BY c.id, c.reply_text, c.media,c.username,c.comment_id,c.likes_count, c.created_at, c.updated_at, u.name, u.profile_picture, u.is_verified
`

type GetRepliesRow struct {
	ID                 uuid.UUID
	ReplyText          string
	Media              json.RawMessage
	Username           string
	CommentID          uuid.UUID
	LikesCount         int32
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Name               string
	ProfilePicture     sql.NullString
	IsVerified         sql.NullBool
	LikedUsersUsername interface{}
}

func (q *Queries) GetReplies(ctx context.Context, commentID uuid.UUID) ([]GetRepliesRow, error) {
	rows, err := q.db.QueryContext(ctx, getReplies, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRepliesRow
	for rows.Next() {
		var i GetRepliesRow
		if err := rows.Scan(
			&i.ID,
			&i.ReplyText,
			&i.Media,
			&i.Username,
			&i.CommentID,
			&i.LikesCount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.ProfilePicture,
			&i.IsVerified,
			&i.LikedUsersUsername,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReply = `-- name: GetReply :one
SELECT id, reply_text, media, username, comment_id, likes_count, created_at, updated_at
FROM comment_reply
WHERE id = $1
`

func (q *Queries) GetReply(ctx context.Context, id uuid.UUID) (CommentReply, error) {
	row := q.db.QueryRowContext(ctx, getReply, id)
	var i CommentReply
	err := row.Scan(
		&i.ID,
		&i.ReplyText,
		&i.Media,
		&i.Username,
		&i.CommentID,
		&i.LikesCount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const likeReply = `-- name: LikeReply :exec
INSERT INTO like_reply(
    username, reply_id
) VALUES (
    $1, $2
)
`

type LikeReplyParams struct {
	Username string
	ReplyID  uuid.UUID
}

func (q *Queries) LikeReply(ctx context.Context, arg LikeReplyParams) error {
	_, err := q.db.ExecContext(ctx, likeReply, arg.Username, arg.ReplyID)
	return err
}

const newReply = `-- name: NewReply :exec

INSERT INTO comment_reply (
    id, reply_text, media, username, comment_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, reply_text, media, username, comment_id, likes_count, created_at, updated_at
`

type NewReplyParams struct {
	ID        uuid.UUID
	ReplyText string
	Media     json.RawMessage
	Username  string
	CommentID uuid.UUID
}

func (q *Queries) NewReply(ctx context.Context, arg NewReplyParams) error {
	_, err := q.db.ExecContext(ctx, newReply,
		arg.ID,
		arg.ReplyText,
		arg.Media,
		arg.Username,
		arg.CommentID,
	)
	return err
}

const removeReplyLike = `-- name: RemoveReplyLike :exec
DELETE FROM like_reply
WHERE username = $1 AND reply_id = $2
`

type RemoveReplyLikeParams struct {
	Username string
	ReplyID  uuid.UUID
}

func (q *Queries) RemoveReplyLike(ctx context.Context, arg RemoveReplyLikeParams) error {
	_, err := q.db.ExecContext(ctx, removeReplyLike, arg.Username, arg.ReplyID)
	return err
}

const updateReplyDecrease = `-- name: UpdateReplyDecrease :one
UPDATE comments
    SET reply_count =  reply_count - 1
WHERE id = $1
RETURNING reply_count
`

func (q *Queries) UpdateReplyDecrease(ctx context.Context, id uuid.UUID) (int32, error) {
	row := q.db.QueryRowContext(ctx, updateReplyDecrease, id)
	var reply_count int32
	err := row.Scan(&reply_count)
	return reply_count, err
}

const updateReplyIncrease = `-- name: UpdateReplyIncrease :one
UPDATE comments
    SET reply_count =  reply_count + 1
WHERE id = $1
RETURNING reply_count
`

func (q *Queries) UpdateReplyIncrease(ctx context.Context, id uuid.UUID) (int32, error) {
	row := q.db.QueryRowContext(ctx, updateReplyIncrease, id)
	var reply_count int32
	err := row.Scan(&reply_count)
	return reply_count, err
}

const updateReplyLikesCountDecrease = `-- name: UpdateReplyLikesCountDecrease :one
UPDATE comment_reply
    SET likes_count = likes_count + 1
WHERE id = $1
RETURNING likes_count, username, reply_text
`

type UpdateReplyLikesCountDecreaseRow struct {
	LikesCount int32
	Username   string
	ReplyText  string
}

func (q *Queries) UpdateReplyLikesCountDecrease(ctx context.Context, id uuid.UUID) (UpdateReplyLikesCountDecreaseRow, error) {
	row := q.db.QueryRowContext(ctx, updateReplyLikesCountDecrease, id)
	var i UpdateReplyLikesCountDecreaseRow
	err := row.Scan(&i.LikesCount, &i.Username, &i.ReplyText)
	return i, err
}

const updateReplyLikesCountIncrease = `-- name: UpdateReplyLikesCountIncrease :one
UPDATE comment_reply
    SET likes_count = likes_count + 1
WHERE id = $1    
RETURNING likes_count, username, reply_text
`

type UpdateReplyLikesCountIncreaseRow struct {
	LikesCount int32
	Username   string
	ReplyText  string
}

func (q *Queries) UpdateReplyLikesCountIncrease(ctx context.Context, id uuid.UUID) (UpdateReplyLikesCountIncreaseRow, error) {
	row := q.db.QueryRowContext(ctx, updateReplyLikesCountIncrease, id)
	var i UpdateReplyLikesCountIncreaseRow
	err := row.Scan(&i.LikesCount, &i.Username, &i.ReplyText)
	return i, err
}