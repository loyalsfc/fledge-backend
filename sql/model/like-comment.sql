-- name: LikeComment :exec

INSERT INTO comment_likes (
    username, comment_id
) VALUES (
    $1, $2
);

-- name: RemoveCommentLike :exec
DELETE FROM comment_likes
    WHERE username = $1 AND comment_id = $2;

-- name: IncreaseCommentLikeCount :one
UPDATE comments
    SET likes_count = likes_count + 1
WHERE id = $1
RETURNING likes_count, username, comment_text;

-- name: DecreaseCommentLikeCount :one
UPDATE comments
    SET likes_count = likes_count - 1
WHERE id = $1
RETURNING likes_count, username;