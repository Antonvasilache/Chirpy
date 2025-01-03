package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain")
	if cfg.PLATFORM != "dev" {
		responseHelper(w, 403, errorResponse{Error: "Forbidden"})
		return
	}

	
	cfg.fileserverHits.Store(0)
	
	err := cfg.Queries.DeleteUsers(r.Context())
	if err != nil {
		log.Printf("Could not delete users: %s", err)
		responseHelper(w, 500, errorResponse{Error: "Internal server error"})
		return
	}
	
	w.WriteHeader(http.StatusOK)	
	w.Write([]byte("Server reset successfully!\n"))
}