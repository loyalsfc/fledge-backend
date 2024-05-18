-- +goose Up 

ALTER TABLE posts
ADD is_shared_post BOOLEAN DEFAULT FALSE NOT NULL,
ADD shared_post_id UUID REFERENCES posts(id);

-- +goose Down

ALTER TABLE posts
DROP is_shared_post,
DROP shared_post_id;