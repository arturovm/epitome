package epitome_test

import (
	"testing"

	"github.com/arturovm/epitome"
	"github.com/stretchr/testify/require"
)

func TestNewSession(t *testing.T) {
	name := "testuser"
	sess, err := epitome.NewSession(name)
	require.NoError(t, err)
	require.NotNil(t, sess)
	require.Equal(t, name, sess.Username)
}
