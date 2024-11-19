package admin

import (
	"crypto/subtle"

	"golang.org/x/crypto/bcrypt"
)

type basicCredentials struct {
	username     string
	passwordHash string
}

func NewBasicCredentials(username, password string) *basicCredentials {
	return &basicCredentials{username: username, passwordHash: password}
}

func (c *basicCredentials) Validate(username, password string) (bool, error) {
	if subtle.ConstantTimeCompare([]byte(username), []byte(c.username)) != 1 {
		return false, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(c.passwordHash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
