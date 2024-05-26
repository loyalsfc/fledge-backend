-- +goose Up

CREATE TABLE comment_reply (
    id UUID UNIQUE NOT NULL,
    reply_text TEXT NOT NULL,
    media JSONB NOT NULL,
    username VARCHAR(50) REFERENCES users(username) ON DELETE CASCADE NOT NULL,
    comment_id UUID REFERENCES comments(id) ON DELETE CASCADE NOT NULL,
    likes_count INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE comment_reply;