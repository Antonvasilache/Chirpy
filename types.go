package main

import (
	"sync/atomic"
	"time"

	"github.com/Antonvasilache/Chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type userCredentials struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ExpiresInSeconds *int `json:"expires_in_seconds"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	Queries *database.Queries
	PLATFORM string
	JWTSECRET string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token	  string	`json:"token"`
}

type ChirpResponse struct{
	ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body      string    `json:"body"`
    UserID    uuid.UUID `json:"user_id"`
}