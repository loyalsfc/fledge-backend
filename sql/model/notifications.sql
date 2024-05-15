-- name: InsertNotification :exec
INSERT INTO notifications (
    id, sender_username, receiver_username, content, notifications_source, reference
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: RemoveNotification :exec
DELETE FROM notifications
WHERE sender_username = $1 AND reference = $2;

-- name: GetUserNotifications :many
SELECT n.*, u.name, u.profile_picture, u.is_verified
FROM notifications n
INNER JOIN users u ON n.sender_username = u.username
WHERE n.receiver_username = $1
ORDER BY n.created_at DESC;