-- +goose Up

CREATE TABLE notifications(
    id UUID UNIQUE NOT NULL,
    sender_username TEXT REFERENCES users(username) ON DELETE CASCADE NOT NULL,
    receiver_username TEXT REFERENCES users(username) ON DELETE CASCADE NOT NULL,
    content VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_viewed BOOLEAN DEFAULT FALSE NOT NULL,
    notifications_source TEXT NOT NULL,
    reference VARCHAR(225) NOT NULL
);

-- +goose Down

DROP TABLE notifications;