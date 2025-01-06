package main

import (
	"log"
	"net/http"

	"github.com/Antonvasilache/Chirpy/internal/auth"
	"github.com/Antonvasilache/Chirpy/internal/helpers"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Could not retrieve bearer token :%s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}	

	err = cfg.Queries.RevokeToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Could not revoke token :%s", err)
		helpers.ResponseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}

	helpers.ResponseHelper(w, 204, nil)	
}