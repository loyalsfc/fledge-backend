-- +goose Up
ALTER TABLE posts
ADD likes_count INT DEFAULT 0 NOT NULL,
ADD comment_count INT DEFAULT 0 NOT NULL,
ADD bookmarks_count INT DEFAULT 0 NOT NULL,
ADD share_count INT DEFAULT 0 NOT NULL;

-- +goose Down
ALTER TABLE posts
DROP likes_count,
DROP comment_count,
DROP bookmarks_count,
DROP share_count;