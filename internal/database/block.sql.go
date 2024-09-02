// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: block.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const block = `-- name: Block :exec

INSERT INTO blocks(
    blocker_id, blocked_id
) VALUES (
    $1, $2
)
RETURNING blocker_id, blocked_id
`

type BlockParams struct {
	BlockerID uuid.UUID
	BlockedID uuid.UUID
}

func (q *Queries) Block(ctx context.Context, arg BlockParams) error {
	_, err := q.db.ExecContext(ctx, block, arg.BlockerID, arg.BlockedID)
	return err
}

const getBlocked = `-- name: GetBlocked :many
SELECT blocked_id FROM blocks
WHERE blocker_id = $1
`

func (q *Queries) GetBlocked(ctx context.Context, blockerID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getBlocked, blockerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var blocked_id uuid.UUID
		if err := rows.Scan(&blocked_id); err != nil {
			return nil, err
		}
		items = append(items, blocked_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unblock = `-- name: Unblock :exec
DELETE FROM blocks
WHERE blocker_id = $1 AND blocked_id = $2
`

type UnblockParams struct {
	BlockerID uuid.UUID
	BlockedID uuid.UUID
}

func (q *Queries) Unblock(ctx context.Context, arg UnblockParams) error {
	_, err := q.db.ExecContext(ctx, unblock, arg.BlockerID, arg.BlockedID)
	return err
}
