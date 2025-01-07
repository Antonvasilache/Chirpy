package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/Antonvasilache/Chirpy/internal/database"
	"github.com/Antonvasilache/Chirpy/internal/helpers"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")	

	authorIDStr := r.URL.Query().Get("author_id")
	authorID, err := uuid.Parse(authorIDStr)
	if authorIDStr != "" && err != nil {
		http.Error(w, "Invalid author ID format", http.StatusBadRequest)
		return
	}

	var databaseChirps []database.Chirp
	if authorIDStr != "" {
		databaseChirps, err = cfg.Queries.GetChirpsByUserId(r.Context(), authorID)		
	} else {
		databaseChirps, err = cfg.Queries.GetChirps(r.Context())
	}
	if err != nil {
		log.Printf("Could not retrieve users: %s", err)
		helpers.ResponseHelper(w, 400, errorResponse{Error: "Error. Please try again later"})
		return
	}

	//sorting chirps
	sortParam := r.URL.Query().Get("sort")
	if sortParam == "desc" {
		sort.Slice(databaseChirps, func(i, j int) bool {
			return databaseChirps[i].CreatedAt.After(databaseChirps[j].CreatedAt)
		})
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