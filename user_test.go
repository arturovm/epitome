package epitome_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arturovm/epitome"
)

func TestCreateUser(t *testing.T) {
	username, password := "testuser", "testpassword"
	u, err := epitome.CreateUser(username, password)
	require.NoError(t, err)
	require.NotNil(t, u)
	require.NotNil(t, u.Credentials())
}

func TestCreateUserMixedCase(t *testing.T) {
	username, password := "TestUser", "testpassword"
	u, _ := epitome.CreateUser(username, password)
	require.Equal(t, u.Username, "testuser")
}
