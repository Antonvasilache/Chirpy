package auth

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(password, hash string) error {
	passwordBytes := []byte(password)
	hashBytes := []byte(hash)
	return bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)	
}