package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Antonvasilache/Chirpy/internal/helpers"
	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	data := WebhookRequest{}
	err := decoder.Decode(&data)
	if err != nil {
		log.Printf("Error decoding data: %s", err)	
		helpers.ResponseHelper(w, 400, errorResponse{Error: err.Error()})
		return
	}

	if data.Event != "user.upgraded" {
		helpers.ResponseHelper(w, 204, nil)
		return
	}

	userID, err := uuid.Parse(data.Data.UserID)
	if err != nil {
		http.Error(w, "Invalid chirp ID format", http.StatusBadRequest)
		return
	}

	err = cfg.Queries.UpgradeUserToRedByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.ResponseHelper(w, 404, nil)
			return
		}
		log.Printf("Error upgrading user: %s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}

	helpers.ResponseHelper(w, 204, nil)
}