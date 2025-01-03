package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antonvasilache/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)	
		responseHelper(w, 400, errorResponse{Error: err.Error()})
		return
	}

	if len(params.Body) > 140 {
		responseHelper(w, 400, errorResponse{Error: "Chirp is too long"})
		return
	}

	cleaned_body := cleanBody(params.Body)

	chirp_id := uuid.New()

	databaseChirp, err := cfg.Queries.CreateChirp(r.Context(), database.CreateChirpParams{
		ID: chirp_id,
		Body: cleaned_body,
		UserID: params.UserID,
	})
	if err != nil {
		log.Printf("Could not create chirp: %s", err)
		responseHelper(w, 400, errorResponse{Error: "Error. Please try again later"})
		return
	}

	response := ChirpResponse{
		ID: databaseChirp.ID,
		CreatedAt: databaseChirp.CreatedAt,
		UpdatedAt: databaseChirp.UpdatedAt,
		Body: databaseChirp.Body,
		UserID: databaseChirp.UserID,
	}

	responseHelper(w, 201, response)
}

