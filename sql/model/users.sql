-- name: CreateUser :one

INSERT INTO users (
    id, name, email, username, password
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: SignIn :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;