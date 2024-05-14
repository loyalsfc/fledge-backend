-- +goose up

CREATE TABLE bookmarks (
    username TEXT REFERENCES users(username) ON DELETE CASCADE,
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    PRIMARY KEY (username, post_id)
);

-- +goose Down
DROP TABLE bookmarks;