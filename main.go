package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "strings"
    "github.com/hf-chow/chirpy/internal/database"
)

const template = `
<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>
</html>
`

type apiConfig struct {
    fileserverHits  int
    DB              *DBB
}

func main() {
    const filepathRoot = "."
    const port = "8080"

    apiCfg := apiConfig {
        fileserverHits: 0,
    }

    mux := http.NewServeMux()
    mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

    mux.HandleFunc("GET /api/healthz", handlerReadiness)
    mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
    mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
    mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)

    mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

    srv := &http.Server{
        Addr:       ":" + port,
        Handler:    mux,
    }

    //    mux.Handle("/app/assets", http.FileServer(http.Dir("/assets/logo.png")))

    log.Printf("Serving on port: %s\n", port)
    log.Fatal(srv.ListenAndServe())
}

func (cfg* apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "application/json; charset=utf-8")

    NewDB("/database/database.json")

    type parameters struct {
        Body string `json:"body"`
    }
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
        w.WriteHeader(500)
        return 
    }

    if len(params.Body) > 140 {
        type returnVals struct {
            Error string `json:"error"`
        }
        respBody := returnVals{
            Error: "Chirp is too long",
        }
        dat, err := json.Marshal(respBody)
        if err != nil {
            log.Printf("Error marshalling JSON: %s", err)
            w.WriteHeader(500)
            return 
        }
        w.WriteHeader(400)
        w.Write(dat)
        return 
    }

    var profanities = [3]string{"kerfuffle", "sharbert", "fornax"}
    cleaned_body := params.Body

    for i := 0; i < len(profanities); i++ {
        cleaned_body = strings.ReplaceAll(cleaned_body, profanities[i], "****")
        cleaned_body = strings.ReplaceAll(cleaned_body, strings.Title(profanities[i]), "****")
    }

    type returnVals struct {
        CleanedBody string `json:"cleaned_body"`
    }
    respBody := returnVals{
        CleanedBody: cleaned_body,
    }
    dat, err := json.Marshal(respBody)
    if err != nil {
        log.Printf("Error marshalling JSON: %s", err)
        w.WriteHeader(500)
        return 
    }
    w.WriteHeader(200)
    w.Write(dat)
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet{
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    w.Header().Add("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf(template, cfg.fileserverHits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits++
        next.ServeHTTP(w, r)
    })
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet{
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits = 0
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hits reset to 0"))
}
