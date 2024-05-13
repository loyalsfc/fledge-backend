-- name: NewComments :exec

INSERT INTO comments (
    id, comment_text, media, username, post_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateCommentIncrease :one
UPDATE posts
    SET comment_count =  comment_count + 1
WHERE id = $1
RETURNING comment_count;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1
RETURNING *;

-- name: UpdateCommentDecrease :one
UPDATE posts
    SET comment_count =  comment_count - 1
WHERE id = $1
RETURNING comment_count;

-- name: GetComments :many
SELECT c.*, u.name, u.profile_picture, u.is_verified
FROM comments c
INNER JOIN users u On u.username = c.username
WHERE c.post_id = $1;