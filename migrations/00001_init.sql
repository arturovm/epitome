-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE users (
	username TEXT NOT NULL PRIMARY KEY,
	password BLOB NOT NULL,
	salt     BLOB NOT NULL
);

CREATE TABLE sessions (
	id       BLOB NOT NULL PRIMARY KEY,
	key      BLOB NOT NULL,
	username TEXT REFERENCES users(id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE sessions;
DROP TABLE users;
