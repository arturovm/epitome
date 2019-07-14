package database_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/arturovm/epitome"
	"github.com/arturovm/epitome/storage"
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

func TestAddExistingUser(t *testing.T) {
	user := epitome.NewUser("conflict", new(epitome.Credentials))

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(user.Username,
			user.Credentials().Password,
			user.Credentials().Salt).
		WillReturnError(sqlite3.Error{
			ExtendedCode: sqlite3.ErrConstraintUnique,
		})

	repo := database.NewUserRepository(db)
	err = repo.Add(user)
	require.EqualError(t, err, storage.ErrUserExists.Error())
}

const getUserQuery = `SELECT password, salt FROM users`

func TestGetUser(t *testing.T) {
	user, _ := epitome.CreateUser("testusername", "testpassword")

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	credentialsRow := sqlmock.NewRows([]string{"password", "salt"}).
		AddRow(user.Credentials().Password, user.Credentials().Salt)
	mock.ExpectQuery(getUserQuery).
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
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)

}

func TestGetNonExistentUser(t *testing.T) {
	badUsername := "userNotExists"

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectQuery(getUserQuery).
		WithArgs(badUsername).
		WillReturnError(sql.ErrNoRows)

	repo := database.NewUserRepository(db)
	u, err := repo.ByUsername(badUsername)
	require.EqualError(t, err, storage.ErrUserNotFound.Error())
	require.Nil(t, u)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
