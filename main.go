package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	hitCountMessage := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(hitCountMessage))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)	

	cfg.fileserverHits.Store(0)
	w.Write([]byte("Counter reset successfully!"))
}


func main(){
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	apiCfg := &apiConfig{}
		
	fileServer := http.FileServer(http.Dir("."))
	stripPrefixHandler := http.StripPrefix("/app/", fileServer)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(stripPrefixHandler))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/reset", apiCfg.resetHandler)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}