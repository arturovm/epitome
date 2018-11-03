-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE users (
	id       INTEGER PRIMARY KEY ASC AUTOINCREMENT,
	username TEXT,
	password BLOB,
	salt     BLOB
);

CREATE UNIQUE INDEX users_username_index ON users (username);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP INDEX  users_username_index;

DROP TABLE users;
