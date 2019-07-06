package epitome

import (
	"crypto/rand"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Password []byte
	Salt     []byte
}

const (
	saltLen    = 16
	bcryptCost = 14
)

// hash password
// hash, salt, err := hashPassword(password)
// if err != nil {
//         log.WithField("error", err).Error("error hashing password")
//         return nil, ErrPasswordHashingFailed
// }

// ErrPasswordHashingFailed is returned when an error occurrs during the
// password hashing process.
var ErrPasswordHashingFailed = errors.New("failed to hash password")

func NewCredentials(password string) (*Credentials, error) {
	hash, salt, err := hashPassword(password)
	if err != nil {
		return nil, ErrPasswordHashingFailed
	}
	return &Credentials{
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
	data := prependSalt(salt, p)

	// hash salt and password
	hash, err = bcrypt.GenerateFromPassword(data, bcryptCost)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error hashing data")
	}

	return hash, salt, err
}

func prependSalt(salt []byte, password string) []byte {
	data := make([]byte, len(salt)+len(password))
	copy(data, salt)
	for i := 0; i < len(password); i++ {
		data[saltLen+i] = password[i]
	}
	return data
}

// PasswordMatch compares a user's hashed password agsinst the given password
// and returns whether they match or not.
func (c *Credentials) MatchPassword(password string) bool {
	data := prependSalt(c.Salt, password)
	err := bcrypt.CompareHashAndPassword(c.Password, data)
	return err == nil
}
