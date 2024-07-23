package main

import (
	"net/http"
    "strconv"

    "github.com/hf-chow/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpDelete (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
        return
    }

    subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
        return
    }

    chirpIDString := r.PathValue("chirpID")

    chirpID, err := strconv.Atoi(chirpIDString)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
        return
    }

    author_id, err := cfg.DB.GetChirpAuthorID(chirpID)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
        return
    }

    requester_id, err := strconv.Atoi(subject) 
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't parse author_id")
        return
    }

    if author_id != requester_id {
        respondWithError(w, http.StatusForbidden, "Couldn't delete chirp")
        return
    }

	err = cfg.DB.DeleteChirp(chirpID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

    respondWithJSON(w, http.StatusNoContent, nil)
}

