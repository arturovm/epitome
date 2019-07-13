package storage

import "errors"

// ErrUserNotFound is returned when a repository method attempts to retrieve
// a user that does not exist.
var ErrUserNotFound = errors.New("user not found")
