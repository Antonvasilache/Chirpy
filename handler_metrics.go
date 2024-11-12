package main

import (
	"fmt"
	"net/http"
)


func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	hitCountMessage := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(hitCountMessage))
}