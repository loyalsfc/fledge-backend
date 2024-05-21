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
SELECT c.*, u.name, u.profile_picture, u.is_verified,
    array_agg(l.username) AS liked_users_username
FROM comments c
INNER JOIN users u On u.username = c.username
LEFT JOIN comment_likes l ON l.comment_id =c.id
WHERE c.post_id = $1
GROUP BY c.id, c.comment_text, c.media,c.username,c.post_id,c.likes_count,c.reply_count, c.created_at, c.updated_at, u.name, u.profile_picture, u.is_verified;