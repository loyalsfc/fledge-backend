// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: followers.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const getFollowers = `-- name: GetFollowers :many
SELECT id, name, username, email, password, bio, profession, is_verified, is_active, profile_picture, cover_picture, created_at, updated_at, follower_id, following_id -- Replace with desired user data
FROM users u
INNER JOIN followers f ON f.follower_id = u.id
WHERE f.following_id = $1
`

type GetFollowersRow struct {
	ID             uuid.UUID
	Name           string
	Username       string
	Email          string
	Password       string
	Bio            sql.NullString
	Profession     sql.NullString
	IsVerified     sql.NullBool
	IsActive       sql.NullBool
	ProfilePicture sql.NullString
	CoverPicture   sql.NullString
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
	FollowerID     uuid.UUID
	FollowingID    uuid.UUID
}

func (q *Queries) GetFollowers(ctx context.Context, followingID uuid.UUID) ([]GetFollowersRow, error) {
	rows, err := q.db.QueryContext(ctx, getFollowers, followingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFollowersRow
	for rows.Next() {
		var i GetFollowersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.Bio,
			&i.Profession,
			&i.IsVerified,
			&i.IsActive,
			&i.ProfilePicture,
			&i.CoverPicture,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FollowerID,
			&i.FollowingID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFollowing = `-- name: GetFollowing :many
SELECT id, name, username, email, password, bio, profession, is_verified, is_active, profile_picture, cover_picture, created_at, updated_at, follower_id, following_id -- Replace with desired user data
FROM users u
INNER JOIN followers f ON f.following_id = u.id
WHERE f.follower_id = $1
`

type GetFollowingRow struct {
	ID             uuid.UUID
	Name           string
	Username       string
	Email          string
	Password       string
	Bio            sql.NullString
	Profession     sql.NullString
	IsVerified     sql.NullBool
	IsActive       sql.NullBool
	ProfilePicture sql.NullString
	CoverPicture   sql.NullString
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
	FollowerID     uuid.UUID
	FollowingID    uuid.UUID
}

func (q *Queries) GetFollowing(ctx context.Context, followerID uuid.UUID) ([]GetFollowingRow, error) {
	rows, err := q.db.QueryContext(ctx, getFollowing, followerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFollowingRow
	for rows.Next() {
		var i GetFollowingRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.Bio,
			&i.Profession,
			&i.IsVerified,
			&i.IsActive,
			&i.ProfilePicture,
			&i.CoverPicture,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FollowerID,
			&i.FollowingID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newFollower = `-- name: NewFollower :one

INSERT INTO followers (
    follower_id, following_id
) VALUES (
    $1, $2
)
RETURNING follower_id, following_id
`

type NewFollowerParams struct {
	FollowerID  uuid.UUID
	FollowingID uuid.UUID
}

func (q *Queries) NewFollower(ctx context.Context, arg NewFollowerParams) (Follower, error) {
	row := q.db.QueryRowContext(ctx, newFollower, arg.FollowerID, arg.FollowingID)
	var i Follower
	err := row.Scan(&i.FollowerID, &i.FollowingID)
	return i, err
}
