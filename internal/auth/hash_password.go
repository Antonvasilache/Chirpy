package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error){
	bytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}