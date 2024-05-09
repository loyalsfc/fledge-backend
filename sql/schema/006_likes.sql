-- +goose Up

CREATE TABLE likes (
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE NOT NULL,
    username TEXT REFERENCES users(username) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (post_id, username)
);

-- +goose Down
DROP TABLE likes;