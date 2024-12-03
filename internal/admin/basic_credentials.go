package admin

import (
	"crypto/subtle"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type basicCredentials struct {
	username     string
	passwordHash string
	logger       *log.Logger
}

func NewBasicCredentials(username, password string, logger *log.Logger) *basicCredentials {
	return &basicCredentials{username: username, passwordHash: password, logger: logger}
}

func (c *basicCredentials) Validate(username, password string) (bool, error) {
	if subtle.ConstantTimeCompare([]byte(username), []byte(c.username)) != 1 {
		c.logger.Println("username is incorrect")
		return false, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(c.passwordHash), []byte(password))
	if err != nil {
		c.logger.Println("password is incorrect")
		return false, nil
	}
	return true, nil
}
