package epitome

import (
	"crypto/rand"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User represents a system user.
type User struct {
	Username string `json:"username"`
	Password []byte `json:"-"`
	Salt     []byte `json:"-"`
}

const (
	saltLen    = 16
	bcryptCost = 14
)

// ErrPasswordHashingFailed is returned when an error occurrs during the
// password hashing process.
var ErrPasswordHashingFailed = errors.New("failed to hash password")

// NewUser initializes a new user with the given parameters.
func NewUser(username, password string) (*User, error) {
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

func hashPassword(p string) (hash []byte, salt []byte, err error) {
	// make salt
	salt = make([]byte, saltLen)
	_, err = rand.Read(salt)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error generating salt")
	}

	// prepend salt to password
	data := make([]byte, len(salt)+len(p))
	copy(data, salt)
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

// PasswordMatch compares a user's hashed password agsinst the given password
// and returns whether they match or not.
func (u *User) PasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	return err == nil
}
