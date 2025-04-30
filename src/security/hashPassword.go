package security

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, string, error) {
	// Generate salt for password
	salt := uuid.New().String()

	// Hash password with bcrypt
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)

	if err != nil {
		return "", "", err
	}

	// Return hashed password and salt
	return string(passwordHash), salt, nil

}
