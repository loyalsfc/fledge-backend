-- name: NewPost :one

INSERT INTO posts (
    id, user_id, username, content, media, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id;

-- name: GetUserPosts :many
SELECT p.*, u.name, u.profile_picture, u.is_verified
FROM posts p
INNER JOIN users u ON p.user_id = u.id
WHERE u.username = $1
ORDER BY p.created_at DESC;

-- name: GetPost :one
SELECT p.*, u.name, u.profile_picture, u.is_verified
FROM posts p
INNER JOIN users u ON p.user_id = u.id
WHERE p.id = $1 LIMIT 1;