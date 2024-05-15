-- +goose Up

CREATE TABLE comments (
    id UUID UNIQUE NOT NULL,
    comment_text TEXT NOT NULL,
    media JSONB NOT NULL,
    username VARCHAR(50) REFERENCES users(username) ON DELETE CASCADE NOT NULL,
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE NOT NULL,
    likes_count INT DEFAULT 0 NOT NULL,
    reply_count INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE comments;