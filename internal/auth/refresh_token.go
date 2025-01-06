package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error){
	randByte := make([]byte, 32)
	_, err := rand.Read(randByte)
	if err != nil {
		return "", err
	}
	
	refreshToken := hex.EncodeToString(randByte)
	return refreshToken, nil
}