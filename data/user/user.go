package user

import (
	"crypto/rand"
	"strings"

	"github.com/arturovm/epitome/data"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User represents a system user
type User struct {
	ID       int64
	Username string
	Password []byte
	Salt     []byte
}

const (
	saltLen                  = 16
	bcryptCost               = 14
	duplicateUserErrorString = "UNIQUE constraint failed: users.username"
)

// ErrInvalidPassword is returned when the given password doesn't satisfy
// the minimum criteria
var ErrInvalidPassword = errors.New("password is invalid")

// ErrPasswordHashingFailed is returned when an error occurred during the
// password hashing process
var ErrPasswordHashingFailed = errors.New("failed to hash password")

// ErrUserExists is returned when a user with the provided username already
// exists
var ErrUserExists = errors.New("a user with that username already exists")

// New initializes a new user with the given parameters
func New(username, password string) (*User, error) {
	// validate password
	valid := validatePassword(password)
	if !valid {
		return nil, ErrInvalidPassword
	}
	// hash password
	hash, salt, err := hashPassword(password)
	if err != nil {
		log.WithField("error", err).Error("error hashing password")
		return nil, ErrPasswordHashingFailed
	}

	// build user
	return &User{
		Username: strings.ToLower(username),
		Password: hash,
		Salt:     salt,
	}, nil
}

func validatePassword(p string) bool {
	// validate that password isn't empty
	if p == "" {
		return false
	}
	return true
}

func hashPassword(p string) (hash []byte, salt []byte, err error) {
	// make salt
	salt = make([]byte, saltLen)
	_, err = rand.Read(salt)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error generating salt")
	}

	// prepend salt to password
	data := make([]byte, len(salt)+len(p))
	for i := range salt {
		data[i] = salt[i]
	}
	for i := 0; i < len(p); i++ {
		data[saltLen+i] = p[i]
	}

	// hash salt and password
	hash, err = bcrypt.GenerateFromPassword(data, bcryptCost)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error hashing data")
	}

	return hash, salt, err
}

// Create initializes a user with the given parameters and writes it to the
// database
func Create(username, password string) (*User, error) {
	u, err := New(username, password)
	if err != nil {
		return nil, err
	}
	err = u.Save()
	if err != nil {
		if err.Error() == duplicateUserErrorString {
			return nil, ErrUserExists
		}
		return nil, err
	}
	return u, nil
}

// Save writes a user to the database
func (u *User) Save() error {
	sess := data.GetSession()
	_, err := sess.InsertInto("users").
		Columns("username", "password", "salt").Record(u).Exec()
	return err
}

// FindByID searches the database for a user with the provided ID
func FindByID(id int) (*User, error) {
	return nil, nil
}

// FindByUsername searches the database for a user with the provided username
func FindByUsername(username string) (*User, error) {
	return nil, nil
}

// Update writes an updated user to the database
func (u *User) Update() error {
	return nil
}

// Delete deletes the user from the database
func (u *User) Delete() error {
	return nil
}
