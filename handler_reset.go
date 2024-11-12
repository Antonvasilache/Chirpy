package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)	

	cfg.fileserverHits.Store(0)
	w.Write([]byte("Counter reset successfully!"))
}