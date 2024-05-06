-- name: CreateUser :one

INSERT INTO users (
    id, name, email, username, password, cover_picture, profile_picture
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
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

-- name: ChangeProfilePicture :one
UPDATE users
set profile_picture = $2
WHERE username = $1
RETURNING *;

-- name: ChangeCoverPicture :one
UPDATE users
set cover_picture = $2
WHERE username = $1
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE users
    set name = $2,
    bio = $3,
    profession = $4
WHERE username = $1
RETURNING *;
    