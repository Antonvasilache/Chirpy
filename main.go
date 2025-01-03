package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Antonvasilache/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)



func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	apiCfg := &apiConfig{
		Queries: dbQueries,
		PLATFORM: os.Getenv("PLATFORM"),
	}
		
	fileServer := http.FileServer(http.Dir("."))
	stripPrefixHandler := http.StripPrefix("/app/", fileServer)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(stripPrefixHandler))
	mux.HandleFunc("GET /api/healthz", readyHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
