package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antonvasilache/Chirpy/internal/auth"
	"github.com/Antonvasilache/Chirpy/internal/database"
	"github.com/Antonvasilache/Chirpy/internal/helpers"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)	
		helpers.ResponseHelper(w, 400, errorResponse{Error: err.Error()})
		return
	}

	//authenticate user
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Could not retrieve bearer token :%s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.JWTSECRET)
	if err != nil {
		log.Printf("Could not validate token :%s", err)
		helpers.ResponseHelper(w, 401, errorResponse{Error: "Unauthorized"})
		return
	}

	if len(params.Body) > 140 {
		helpers.ResponseHelper(w, 400, errorResponse{Error: "Chirp is too long"})
		return
	}

	cleaned_body := helpers.CleanBody(params.Body)

	chirp_id := uuid.New()

	databaseChirp, err := cfg.Queries.CreateChirp(r.Context(), database.CreateChirpParams{
		ID: chirp_id,
		Body: cleaned_body,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Could not create chirp: %s", err)
		helpers.ResponseHelper(w, 400, errorResponse{Error: "Error. Please try again later"})
		return
	}

	response := ChirpResponse{
		ID: databaseChirp.ID,
		CreatedAt: databaseChirp.CreatedAt,
		UpdatedAt: databaseChirp.UpdatedAt,
		Body: databaseChirp.Body,
		UserID: databaseChirp.UserID,
	}

	helpers.ResponseHelper(w, 201, response)
}

