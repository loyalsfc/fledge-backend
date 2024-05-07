-- name: NewFollower :one

INSERT INTO followers (
    follower_id, following_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetFollowing :many
SELECT * -- Replace with desired user data
FROM users u
INNER JOIN followers f ON f.following_id = u.id
WHERE f.follower_id = $1;

-- name: GetFollowers :many
SELECT * -- Replace with desired user data
FROM users u
INNER JOIN followers f ON f.follower_id = u.id
WHERE f.following_id = $1;

