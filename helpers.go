package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func responseHelper(w http.ResponseWriter, code int, payload interface{}){	
	data, err := json.Marshal(payload)
		if err != nil {
			log.Printf("error marshalling response: %s", err)
			w.WriteHeader(500)			
			return
		}
		w.WriteHeader(code)
		w.Write(data)
}

func cleanBody(body string) string{
	words := strings.Split(body, " ")
	for index, word := range words {
		word_to_lower := strings.ToLower(word)
		if (word_to_lower == "kerfuffle" || 
			word_to_lower == "sharbert" || 
			word_to_lower == "fornax") {
				words[index] = "****"
		}
	}
	cleaned_body := strings.Join(words, " ")
	return cleaned_body
}