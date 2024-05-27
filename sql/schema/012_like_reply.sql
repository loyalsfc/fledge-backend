-- +goose Up 

CREATE TABLE like_reply (
    username TEXT REFERENCES users(username) ON DELETE CASCADE NOT NULL,
    reply_id UUID REFERENCES comment_reply(id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (username, reply_id)
);

-- +goose Down

DROP TABLE like_reply;