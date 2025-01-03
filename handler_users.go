package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	email := createUserRequest{}
	err := decoder.Decode(&email)
	if err != nil {
		log.Printf("Error decoding user data: %s", err)
		responseHelper(w, 400, errorResponse{Error: err.Error()})
		return
	}

	databaseUser, err := cfg.Queries.CreateUser(r.Context(), email.Email)
	if err != nil {
		log.Printf("Could not create user: %s", err)
		switch {		
		case strings.Contains(err.Error(), "duplicate key"):
			responseHelper(w, 409, errorResponse{Error: "Email already exists"})
		default:
			responseHelper(w, 500, errorResponse{Error: "Internal server error"})
		}		
		return
	}

	user := User{
		ID: databaseUser.ID,
		CreatedAt: databaseUser.CreatedAt,
		UpdatedAt: databaseUser.UpdatedAt,
		Email: databaseUser.Email,
	}

	responseHelper(w, 201, user)	
}