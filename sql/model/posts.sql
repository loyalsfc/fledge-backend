-- name: NewPost :one

INSERT INTO posts (
    id, user_id, username, content, media, created_at, updated_at, is_shared_post, shared_post_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id;

-- name: GetUserPosts :many
SELECT p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, p.is_shared_post, p.shared_post_id, u.name, u.profile_picture, u.is_verified,
    array_agg(l.username) AS liked_users_username,
    array_agg(b.username) AS bookmarked_users_username
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN bookmarks b ON p.id = b.post_id
WHERE p.username = $1
GROUP BY p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, p.is_shared_post, p.shared_post_id, u.name, u.profile_picture, u.is_verified
ORDER BY p.created_at DESC;

-- name: GetPost :one
SELECT p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, p.is_shared_post, p.shared_post_id, u.name, u.profile_picture, u.is_verified,
    array_agg(l.username) AS liked_users_username,
	array_agg(b.username) AS bookmarked_users_username
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN bookmarks b ON p.id = b.post_id
WHERE p.id = $1
GROUP BY p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, p.is_shared_post, p.shared_post_id, u.name, u.profile_picture, u.is_verified
ORDER BY p.created_at DESC;

-- name: GetFeedPosts :many
SELECT p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, p.is_shared_post, p.shared_post_id, u.name, u.profile_picture, u.is_verified,
    array_agg(l.username) AS liked_users_username,
    array_agg(b.username) AS bookmarked_users_username
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN bookmarks b ON p.id = b.post_id
LEFT JOIN followers f ON p.user_id = f.following_id
WHERE f.follower_id = $1 OR p.user_id = $1
GROUP BY p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, p.is_shared_post, p.shared_post_id, u.name, u.profile_picture, u.is_verified
ORDER BY p.created_at DESC;

-- name: GetBookmarkedPosts :many
SELECT p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, p.is_shared_post, p.shared_post_id, u.name, u.profile_picture, u.is_verified,
    array_agg(l.username) AS liked_users_username,
    array_agg(b.username) AS bookmarked_users_username
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
INNER JOIN bookmarks b ON p.id = b.post_id
WHERE b.username = $1
GROUP BY p.id, p.user_id, p.content, p.media, p.username, p.created_at, p.updated_at, p.likes_count, p.comment_count, p.bookmarks_count, p.share_count, p.is_shared_post, p.shared_post_id, u.name, u.profile_picture, u.is_verified;

-- name: IncreaseShareCount :one
UPDATE posts
    SET share_count = share_count + 1
WHERE id = $1
RETURNING share_count;

-- name: DecreaseShareCount :one
UPDATE posts
    SET share_count = share_count - 1
WHERE id = $1
RETURNING share_count;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

-- name: EditPost :exec
UPDATE posts
    SET content = $1,
    media = $2,
    updated_at = now()
WHERE id = $3;