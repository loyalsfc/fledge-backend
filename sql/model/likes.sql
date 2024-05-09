-- name: LikePost :exec

INSERT INTO likes (
    username, post_id
) VALUES (
    $1, $2
);

-- name: UnlikePost :exec
DELETE FROM likes 
WHERE username = $1 AND post_id = $2;