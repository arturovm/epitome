package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/arturovm/epitome"
)

type SessionRepository struct {
	mock.Mock
}

func (r *SessionRepository) Add(session epitome.Session) error {
	args := r.Called(session)
	return args.Error(0)
}
