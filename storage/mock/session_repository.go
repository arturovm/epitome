package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/arturovm/epitome"
)

// SessionRepository is a mock implementation of storage.SessionRepository.
type SessionRepository struct {
	mock.Mock
}

// Add implements SessionRepository.Add.
func (r *SessionRepository) Add(session epitome.Session) error {
	args := r.Called(session)
	return args.Error(0)
}
