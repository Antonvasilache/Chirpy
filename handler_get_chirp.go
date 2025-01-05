package main

import (
	"log"
	"net/http"

	"github.com/Antonvasilache/Chirpy/internal/helpers"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	chirpIDstr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		http.Error(w, "Invalid chirp ID format", http.StatusBadRequest)
	}

	databaseChirp, err := cfg.Queries.GetChirpById(r.Context(), chirpID)
	if err != nil {
		log.Printf("Could not retrieve chirp: %s", err)
		helpers.ResponseHelper(w, 404, errorResponse{Error: "Error. Chirp was not found"})
		return
	}

	response := ChirpResponse{
		ID: databaseChirp.ID,
		CreatedAt: databaseChirp.CreatedAt,
		UpdatedAt: databaseChirp.UpdatedAt,
		Body: databaseChirp.Body,
		UserID: databaseChirp.UserID,
	}

	helpers.ResponseHelper(w, 200, response)
}