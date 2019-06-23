package epitome

import (
	"crypto/rand"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Session represents a user's session on the server.
type Session struct {
	ID       uuid.UUID
	Key      []byte
	Username string
}

const keyLen = 256

// NewSession generates a new session for the given username.
func NewSession(username string) (*Session, error) {
	key := make([]byte, keyLen/8)
	_, err := rand.Read(key)
	if err != nil {
		return nil, errors.Wrap(err, "error reading random bytes")
	}

	return &Session{
		ID:       uuid.New(),
		Key:      key,
		Username: username,
	}, nil
}
