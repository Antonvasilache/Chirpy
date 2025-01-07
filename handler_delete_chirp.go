package main

import (
	"log"
	"net/http"

	"github.com/Antonvasilache/Chirpy/internal/auth"
	"github.com/Antonvasilache/Chirpy/internal/helpers"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	//1. Get chirp ID from route
	chirpIDstr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		http.Error(w, "Invalid chirp ID format", http.StatusBadRequest)
		return
	}

	//2. Get user ID from chirp ID
	dbUserID, err := cfg.Queries.GetUserIdByChirpId(r.Context(), chirpID)
	if err != nil {
		log.Printf("Could not retrieve user ID: %s", err)
		helpers.ResponseHelper(w, 404, errorResponse{Error: "Error. User ID was not found"})
		return
	}

	//3. Check token in header, get user ID from token
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

	//4. Compare db user ID with token user ID
	if userID != dbUserID {
		log.Printf("User is not the author of the chirp: %s", err)
		helpers.ResponseHelper(w, 403, errorResponse{Error: "Forbidden"})
		return
	}

	//5. Delete chirp from database
	err = cfg.Queries.DeletChirpById(r.Context(), chirpID)
	if err != nil {
		log.Printf("Could not delete chirp: %s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}

	helpers.ResponseHelper(w, 204, nil)
}