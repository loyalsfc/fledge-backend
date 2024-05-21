-- +goose Up

CREATE TABLE comment_likes(
    username TEXT REFERENCES users(username) ON DELETE CASCADE NOT NULL,
    comment_id UUID REFERENCES comments(id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (username, comment_id)
)


-- +goose Down
DROP TABLE comment_likes