package main

import (
	"log"
	"net/http"

	"github.com/Antonvasilache/Chirpy/internal/helpers"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")	

	databaseChirps, err := cfg.Queries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Could not retrieve users: %s", err)
		helpers.ResponseHelper(w, 400, errorResponse{Error: "Error. Please try again later"})
		return
	}

	response := make([]ChirpResponse, len(databaseChirps))
	for index, dbChirp := range databaseChirps {
		chirp := ChirpResponse{
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			UserID: dbChirp.UserID,
		}
		response[index] = chirp
	}

	helpers.ResponseHelper(w, 200, response)	
}