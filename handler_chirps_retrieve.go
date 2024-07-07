package main

import (
	"net/http"
	"sort"
    "strconv"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsRetrieveById(w http.ResponseWriter, r *http.Request) {

    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
    }

    chirp, err := cfg.DB.GetChirpById(id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chrip ID not found")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
