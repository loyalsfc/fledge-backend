-- name: LikePost :exec

INSERT INTO likes (
    username, post_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateLikeIncrease :one
UPDATE posts
SET likes_count = likes_count + 1
WHERE id = $1
RETURNING likes_count;

-- name: UnlikePost :exec
DELETE FROM likes 
WHERE username = $1 AND post_id = $2;

-- name: UpdateLikeDecrease :one
UPDATE posts
SET likes_count = likes_count - 1
WHERE id = $1
RETURNING likes_count;

-- name: GetPostLikes :many
SELECT * FROM likes
WHERE post_id = $1;