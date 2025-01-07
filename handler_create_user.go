package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Antonvasilache/Chirpy/internal/auth"
	"github.com/Antonvasilache/Chirpy/internal/database"
	"github.com/Antonvasilache/Chirpy/internal/helpers"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	userRequest := userCredentials{}
	err := decoder.Decode(&userRequest)
	if err != nil {
		log.Printf("Error decoding user data: %s", err)		
		helpers.ResponseHelper(w, 400, errorResponse{Error: err.Error()})
		return
	}

	hashedPassword, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}

	databaseUser, err := cfg.Queries.CreateUser(r.Context(), database.CreateUserParams{
		Email: userRequest.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Could not create user: %s", err)
		switch {		
		case strings.Contains(err.Error(), "duplicate key"):
			helpers.ResponseHelper(w, 409, errorResponse{Error: "Email already exists"})
		default:
			helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		}		
		return
	}

	user := User{
		ID: databaseUser.ID,
		CreatedAt: databaseUser.CreatedAt,
		UpdatedAt: databaseUser.UpdatedAt,
		Email: databaseUser.Email,
		IsChirpyRed: databaseUser.IsChirpyRed,
	}

	helpers.ResponseHelper(w, 201, user)	
}