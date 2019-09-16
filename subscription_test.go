package epitome_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arturovm/epitome"
)

func TestSubscribe(t *testing.T) {
	user := epitome.User{Username: "johndoe"}
	source := epitome.Source{URL: "google.com"}
	subscription := user.Subscribe(&source)

	require.Equal(t, user.Username, subscription.Username)
	require.Equal(t, source.URL, subscription.SourceURL)
}

/*
To-do list:

- Add user ID?
- Add source ID?
- Separate Source into its own file?
*/
