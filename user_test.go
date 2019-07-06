package epitome_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arturovm/epitome"
)

func TestNewUser(t *testing.T) {
	name := "testuser"
	u := epitome.NewUser(name)
	require.Equal(t, u.Username, name)
}

func TestNewUserMixedCase(t *testing.T) {
	name := "TestUser"
	u := epitome.NewUser(name)
	require.Equal(t, u.Username, "testuser")
}
