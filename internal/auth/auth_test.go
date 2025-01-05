package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Antonvasilache/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestMakeJWTAndValidate(t *testing.T){
	tokenSecret := "secret"
	userID := uuid.New()

	//Create JWT with a 1-hour expiration
	jwt, err := auth.MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("Unexpected error creating JWT: %v", err)
	}

	//Validate JWT
	validatedID, err := auth.ValidateJWT(jwt, tokenSecret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if validatedID != userID {
		t.Errorf("Expected userID %v, got %v", userID, validatedID)
	}
}

func TestValidateJWTExpired(t *testing.T){
	tokenSecret := "secret"
	userID := uuid.New()

	//Create JWT with a past expiration date
	jwt, err := auth.MakeJWT(userID, tokenSecret, -time.Hour)
	if err != nil {
		t.Fatalf("Unexpected error creating JWT: %v", err)
	}

	//Attempt to validate the expired JWT
	_, err = auth.ValidateJWT(jwt, tokenSecret)
	if err == nil {
		t.Error("Expected error for expired JWT, got nil")
	}
}

func TestValidateJWTWithWrongSecret(t *testing.T){
	tokenSecret := "secret"
	wrongSecret := "wrongsecret"
	userID := uuid.New()

	//Create JWT
	jwt, err := auth.MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("Unexpected error creating JWT: %v", err)
	}

	//Attempt to validate JWT with wrong secret
	_, err = auth.ValidateJWT(jwt, wrongSecret)
	if err == nil {
		t.Error("Expected error for wrong secret, got nil")
	}
}

func TestGetBearerToken(t *testing.T){
	tests := []struct {
		name			string
		headers 		http.Header
		expectedToken	string
		expectError		bool
	}{
		{
			name: "No authorization header",
			headers: http.Header{},
			expectedToken: "",
			expectError: true,
		},
		{
			name: "Invalid format",
			headers: http.Header{
				"Authorization": []string{"Invalid format"},
			},
			expectedToken: "",
			expectError: true,
		},
		{
			name: "Valid Bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer token123"},
			},
			expectedToken: "token123",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.GetBearerToken(tt.headers)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
			if token != tt.expectedToken {
				t.Errorf("Expected token: %s, got: %s", tt.expectedToken, token)
			}
		})
	}
}