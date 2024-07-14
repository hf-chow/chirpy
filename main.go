package main

import (
	"log"
	"net/http"
    "os"

	"github.com/hf-chow/chirpy/internal/database"
    "github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits      int
	DB                  *database.DB
    jwtSecret           string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }

    jwtSecret := os.Getenv("JWT_SECRET")

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits:     0,
		DB:                 db,
        jwtSecret:          jwtSecret,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
    mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerChirpsRetrieveById)
    mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
    mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
    mux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)
   

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
