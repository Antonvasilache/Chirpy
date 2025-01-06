package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Antonvasilache/Chirpy/internal/auth"
	"github.com/Antonvasilache/Chirpy/internal/helpers"
)

func (cfg *apiConfig) loginHandler (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	loginRequest := userCredentials{}
	err := decoder.Decode(&loginRequest)
	if err != nil {
		log.Printf("Error decoding login: %s", err)
		helpers.ResponseHelper(w, 400, errorResponse{Error: err.Error()})
		return
	}

	databaseUser, err := cfg.Queries.GetUserByEmail(r.Context(), loginRequest.Email)
	if err != nil {
		log.Printf("Could not retrieve user: %s", err)
		helpers.ResponseHelper(w, 401, errorResponse{Error: "Incorrect email or password"})
		return
	}

	err = auth.CheckPasswordHash(loginRequest.Password, databaseUser.HashedPassword)
	if err != nil{
		log.Printf("Incorrect password: %s", err)
		helpers.ResponseHelper(w, 401, errorResponse{Error: "Incorrect email or password"})
		return
	}

	expirationTime := time.Hour
	if loginRequest.ExpiresInSeconds != nil && *loginRequest.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(*loginRequest.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(databaseUser.ID, cfg.JWTSECRET, expirationTime)
	if err != nil {
		log.Printf("Could not create JWT Token: %s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}

	user := User{
		ID: databaseUser.ID,
		CreatedAt: databaseUser.CreatedAt,
		UpdatedAt: databaseUser.UpdatedAt,
		Email: databaseUser.Email,
		Token: token,
	}

	helpers.ResponseHelper(w, 200, user)	
}