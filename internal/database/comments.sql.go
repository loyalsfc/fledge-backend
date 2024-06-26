// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: comments.sql

package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const deleteComment = `-- name: DeleteComment :one
DELETE FROM comments
WHERE id = $1
RETURNING post_id
`

func (q *Queries) DeleteComment(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, deleteComment, id)
	var post_id uuid.UUID
	err := row.Scan(&post_id)
	return post_id, err
}

const editComment = `-- name: EditComment :exec
UPDATE comments
    SET comment_text = $1,
    media = $2,
    updated_at = now()
WHERE id = $3
`

type EditCommentParams struct {
	CommentText string
	Media       json.RawMessage
	ID          uuid.UUID
}

func (q *Queries) EditComment(ctx context.Context, arg EditCommentParams) error {
	_, err := q.db.ExecContext(ctx, editComment, arg.CommentText, arg.Media, arg.ID)
	return err
}

const getComment = `-- name: GetComment :one
SELECT comment_text, username 
FROM comments
WHERE id = $1
`

type GetCommentRow struct {
	CommentText string
	Username    string
}

func (q *Queries) GetComment(ctx context.Context, id uuid.UUID) (GetCommentRow, error) {
	row := q.db.QueryRowContext(ctx, getComment, id)
	var i GetCommentRow
	err := row.Scan(&i.CommentText, &i.Username)
	return i, err
}

const getComments = `-- name: GetComments :many
SELECT c.id, c.comment_text, c.media, c.username, c.post_id, c.likes_count, c.reply_count, c.created_at, c.updated_at, u.name, u.profile_picture, u.is_verified,
    array_agg(l.username) AS liked_users_username
FROM comments c
INNER JOIN users u On u.username = c.username
LEFT JOIN comment_likes l ON l.comment_id =c.id
WHERE c.post_id = $1
GROUP BY c.id, c.comment_text, c.media,c.username,c.post_id,c.likes_count,c.reply_count, c.created_at, c.updated_at, u.name, u.profile_picture, u.is_verified
`

type GetCommentsRow struct {
	ID                 uuid.UUID
	CommentText        string
	Media              json.RawMessage
	Username           string
	PostID             uuid.UUID
	LikesCount         int32
	ReplyCount         int32
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Name               string
	ProfilePicture     sql.NullString
	IsVerified         sql.NullBool
	LikedUsersUsername interface{}
}

func (q *Queries) GetComments(ctx context.Context, postID uuid.UUID) ([]GetCommentsRow, error) {
	rows, err := q.db.QueryContext(ctx, getComments, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommentsRow
	for rows.Next() {
		var i GetCommentsRow
		if err := rows.Scan(
			&i.ID,
			&i.CommentText,
			&i.Media,
			&i.Username,
			&i.PostID,
			&i.LikesCount,
			&i.ReplyCount,
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

const newComments = `-- name: NewComments :exec

INSERT INTO comments (
    id, comment_text, media, username, post_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, comment_text, media, username, post_id, likes_count, reply_count, created_at, updated_at
`

type NewCommentsParams struct {
	ID          uuid.UUID
	CommentText string
	Media       json.RawMessage
	Username    string
	PostID      uuid.UUID
}

func (q *Queries) NewComments(ctx context.Context, arg NewCommentsParams) error {
	_, err := q.db.ExecContext(ctx, newComments,
		arg.ID,
		arg.CommentText,
		arg.Media,
		arg.Username,
		arg.PostID,
	)
	return err
}

const updateCommentDecrease = `-- name: UpdateCommentDecrease :one
UPDATE posts
    SET comment_count =  comment_count - 1
WHERE id = $1
RETURNING comment_count
`

func (q *Queries) UpdateCommentDecrease(ctx context.Context, id uuid.UUID) (int32, error) {
	row := q.db.QueryRowContext(ctx, updateCommentDecrease, id)
	var comment_count int32
	err := row.Scan(&comment_count)
	return comment_count, err
}

const updateCommentIncrease = `-- name: UpdateCommentIncrease :one
UPDATE posts
    SET comment_count =  comment_count + 1
WHERE id = $1
RETURNING comment_count
`

func (q *Queries) UpdateCommentIncrease(ctx context.Context, id uuid.UUID) (int32, error) {
	row := q.db.QueryRowContext(ctx, updateCommentIncrease, id)
	var comment_count int32
	err := row.Scan(&comment_count)
	return comment_count, err
}
