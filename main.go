package main

import (
	"log"
	"net/http"
)

func main(){
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))	
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}