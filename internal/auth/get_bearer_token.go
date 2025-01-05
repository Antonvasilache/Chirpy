package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}