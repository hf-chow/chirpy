package main

import (
	"encoding/json"
	"net/http"

    "github.com/hf-chow/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerWebhooks (w http.ResponseWriter, r *http.Request) {
    type data struct {
        UserID  int    `json:"user_id"`
    }
	type parameters struct {
        Data    data
		Event   string `json:"event"`
	}

    apiKey, err := auth.GetAPIKey(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't parse API Key")
        return 
    }

    err = auth.ValidateAPIKey(apiKey, cfg.apiKey)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid API Key")
        return 
    }

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

    if params.Event == "user.upgraded" {
        err = cfg.DB.UpgradeUser(params.Data.UserID)
        if err != nil {
    		respondWithError(w, http.StatusNotFound, "Couldn't find user")
        }
        respondWithJSON(w, http.StatusNoContent, nil)
    } else {
        respondWithJSON(w, http.StatusNoContent, nil)
    }
}
