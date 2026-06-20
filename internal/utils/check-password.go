package utils

import "golang.org/x/crypto/bcrypt"

func CheckPassword(hashedPassword, rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	if err != nil {
		return false
	}
	return true
}
