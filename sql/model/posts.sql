-- name: NewPost :one

INSERT INTO posts (
    id, user_id, username, content, media, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id;

-- name: GetUserPosts :many
SELECT p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, u.name, u.profile_picture, u.is_verified,
    array_agg(CAST(l.username AS VARCHAR)) AS liked_users_username
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
WHERE p.username = $1
GROUP BY p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, u.name, u.profile_picture, u.is_verified
ORDER BY p.created_at DESC;

-- name: GetPost :one
SELECT p.*, u.name, u.profile_picture, u.is_verified
FROM posts p
INNER JOIN users u ON p.user_id = u.id
WHERE p.id = $1 LIMIT 1;

-- name: GetFeedPosts :many
SELECT p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, u.name, u.profile_picture, u.is_verified,
    array_agg(CAST(l.username AS VARCHAR)) AS liked_users_username
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN followers f ON p.user_id = f.following_id
WHERE f.follower_id = $1 OR p.user_id = $1
GROUP BY p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, u.name, u.profile_picture, u.is_verified
ORDER BY p.created_at DESC;