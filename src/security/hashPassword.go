package security

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, string, error) {
	salt := uuid.New().String()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)

	if err != nil {
		return "", "", err
	}

	return string(passwordHash), salt, nil

}
