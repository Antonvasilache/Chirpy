package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antonvasilache/Chirpy/internal/auth"
	"github.com/Antonvasilache/Chirpy/internal/database"
	"github.com/Antonvasilache/Chirpy/internal/helpers"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Could not retrieve bearer token :%s", err)
		helpers.ResponseHelper(w, 401, errorResponse{Error: "Unauthorized"})
		return
	}
	
	userID, err := auth.ValidateJWT(tokenString, cfg.JWTSECRET)
	if err != nil {
		log.Printf("Could not validate token :%s", err)
		helpers.ResponseHelper(w, 401, errorResponse{Error: "Unauthorized"})
		return
	}
	
	decoder := json.NewDecoder(r.Body)
	updateRequest := userCredentials{}
	err = decoder.Decode(&updateRequest)
	if err != nil {
		log.Printf("Error decoding user data: %s", err)		
		helpers.ResponseHelper(w, 400, errorResponse{Error: err.Error()})
		return
	}

	newHashedPassword, err := auth.HashPassword(updateRequest.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}

	databaseUser, err := cfg.Queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: updateRequest.Email,
		HashedPassword: newHashedPassword,
		ID: userID,
	})
	if err != nil {
		log.Printf("Error updating user: %s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}

	user := User{
		ID: databaseUser.ID,
		CreatedAt: databaseUser.CreatedAt,
		UpdatedAt: databaseUser.UpdatedAt,
		Email: databaseUser.Email,
	}

	helpers.ResponseHelper(w, 200, user)	
}