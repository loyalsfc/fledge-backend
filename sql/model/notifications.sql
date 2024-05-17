-- name: InsertNotification :exec
INSERT INTO notifications (
    id, sender_username, receiver_username, content, notifications_source, reference, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: RemoveNotification :exec
DELETE FROM notifications
WHERE sender_username = $1 AND reference = $2 AND notifications_source = $3;

-- name: GetUserNotifications :many
SELECT n.*, u.name, u.profile_picture, u.is_verified
FROM notifications n
INNER JOIN users u ON n.sender_username = u.username
WHERE n.receiver_username = $1
ORDER BY n.created_at DESC;

-- name: MarkNotificationAsRead :exec
UPDATE notifications
    SET is_viewed = true
WHERE sender_username = $1 AND reference = $2 AND notifications_source = $3;