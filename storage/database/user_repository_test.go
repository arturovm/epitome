package database_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/arturovm/epitome"
	"github.com/arturovm/epitome/storage/database"
)

func TestAddUser(t *testing.T) {
	user, _ := epitome.CreateUser("testusername", "testpassword")

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(user.Username,
			user.Credentials().Password,
			user.Credentials().Salt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := database.NewUserRepository(db)
	err = repo.Add(*user)
	require.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	user, _ := epitome.CreateUser("testusername", "testpassword")

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	credentialsRow := sqlmock.NewRows([]string{"password", "salt"}).
		AddRow(user.Credentials().Password, user.Credentials().Salt)
	mock.ExpectQuery(`SELECT password, salt FROM users`).
		WithArgs(user.Username).
		WillReturnRows(credentialsRow)

	repo := database.NewUserRepository(db)
	resp, err := repo.ByUsername(user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, resp.Username)
	require.Equal(t,
		user.Credentials().Password,
		resp.Credentials().Password)
	require.Equal(t,
		user.Credentials().Salt,
		resp.Credentials().Salt)

}
