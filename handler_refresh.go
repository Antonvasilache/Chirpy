package main

import (
	"log"
	"net/http"

	"github.com/Antonvasilache/Chirpy/internal/auth"
	"github.com/Antonvasilache/Chirpy/internal/helpers"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Could not retrieve bearer token :%s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}	

	userID, err := cfg.Queries.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Error retrieving valid token: %s", err)
		helpers.ResponseHelper(w, 401, errorResponse{Error: "Unauthorized"})
		return
	}

	token, err := auth.MakeJWT(userID, cfg.JWTSECRET)
	if err != nil {
		log.Printf("Could not create JWT Token: %s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}
	
	helpers.ResponseHelper(w, 200, TokenResponse{Token: token})	
}