-- +goose Up
CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL, 
    pub_date INTEGER DEFAULT 0,
    link TEXT NOT NULL UNIQUE
);
-- +goose Down
DROP TABLE feeds;
