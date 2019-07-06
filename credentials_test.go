package epitome_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arturovm/epitome"
)

func TestNewCredentials(t *testing.T) {
	password := "testpassword"
	creds, err := epitome.NewCredentials(password)
	require.NoError(t, err)
	require.NotNil(t, creds)
	require.NotEqual(t, password, string(creds.Password))
}

func TestMatchPassword(t *testing.T) {
	password := "this is a test password"
	creds, _ := epitome.NewCredentials(password)

	match := creds.MatchPassword(password)
	require.True(t, match)
}

func TestSetCredentials(t *testing.T) {
	username, password := "testuser", "test password"
	user := epitome.NewUser(username)
	creds, _ := epitome.NewCredentials(password)
	user.SetCredentials(creds)

	require.NotNil(t, user.Credentials())
	require.Equal(t, user.Credentials(), creds)
}
