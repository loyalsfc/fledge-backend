-- name: NewPost :one

INSERT INTO posts (
    id, user_id, username, content, media, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;