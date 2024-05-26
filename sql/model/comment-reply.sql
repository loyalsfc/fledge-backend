-- name: NewReply :exec

INSERT INTO comment_reply (
    id, reply_text, media, username, comment_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateReplyIncrease :one
UPDATE comments
    SET reply_count =  reply_count + 1
WHERE id = $1
RETURNING reply_count;

-- name: DeleteReply :exec
DELETE FROM comment_reply
WHERE id = $1
RETURNING *;

-- name: UpdateReplyDecrease :one
UPDATE comments
    SET reply_count =  reply_count - 1
WHERE id = $1
RETURNING reply_count;

-- name: GetReplies :many
SELECT c.*, u.name, u.profile_picture, u.is_verified,
    array_agg(l.username) AS liked_users_username
FROM comment_reply c
INNER JOIN users u On u.username = c.username
LEFT JOIN comment_likes l ON l.comment_id =c.id
WHERE c.comment_id = $1
GROUP BY c.id, c.reply_text, c.media,c.username,c.comment_id,c.likes_count, c.created_at, c.updated_at, u.name, u.profile_picture, u.is_verified;

-- name: GetReply :one
SELECT *
FROM comment_reply
WHERE id = $1;