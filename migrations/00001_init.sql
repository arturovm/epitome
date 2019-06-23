-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE users (
	username TEXT PRIMARY KEY,
	password BLOB,
	salt     BLOB
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;
