package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) validateHandler(w http.ResponseWriter, r *http.Request){
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

	responseHelper(w, 200, validResponse{Cleaned_body: cleaned_body})
}

