-- name: LikePost :exec

INSERT INTO likes (
    username, post_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateLikeIncrease :exec
UPDATE posts
SET likes_count = likes_count + 1
WHERE id = $1;

-- name: UnlikePost :exec
DELETE FROM likes 
WHERE username = $1 AND post_id = $2;

-- name: UpdateLikeDecrease :exec
UPDATE posts
SET likes_count = likes_count - 1
WHERE id = $1;