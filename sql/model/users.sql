-- name: CreateUser :one

INSERT INTO users (
    id, name, email, username, password
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUser :many
SELECT * FROM users;