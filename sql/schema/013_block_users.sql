-- +goose Up 

CREATE TABLE blocks(
    blocker_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    blocked_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (blocker_id, blocked_id)
)

-- +goose Down
DROP TABLE blocks