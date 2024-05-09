-- name: LikePost :exec

INSERT INTO likes (
    user_id, post_id
) VALUES (
    $1, $2
);

-- name: UnlikePost :exec
DELETE FROM likes 
WHERE user_id = $1 AND post_id = $2;