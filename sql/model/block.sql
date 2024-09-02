-- name: Block :exec

INSERT INTO blocks(
    blocker_id, blocked_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: Unblock :exec
DELETE FROM blocks
WHERE blocker_id = $1 AND blocked_id = $2;

-- name: GetBlocked :many
SELECT blocked_id FROM blocks
WHERE blocker_id = $1;