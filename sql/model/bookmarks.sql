-- name: AddBookmarks :exec
INSERT INTO bookmarks(
    username, post_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateBookmarksIncrease :one
UPDATE posts
    SET bookmarks_count = bookmarks_count + 1
WHERE id = $1
RETURNING bookmarks_count;

-- name: RemoveBookmarks :exec
DELETE FROM bookmarks
WHERE username = $1 AND post_id = $2;

-- name: UpdateBookmarksDecrease :one
UPDATE posts
    SET bookmarks_count = bookmarks_count - 1
WHERE id = $1
RETURNING bookmarks_count;