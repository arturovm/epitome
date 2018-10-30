-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE users (
	id       INTEGER PRIMARY KEY ASC,
	email    TEXT,
	password BLOB
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;
