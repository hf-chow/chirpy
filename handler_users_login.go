package main

import (
	"encoding/json"
	"net/http"
    "golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email       string      `json:"email"`
        Password    string      `json:"password"`
	}

    dbUsers, err := cfg.DB.GetUsers()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Counldn't retrieve users")
    }

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

    for _, dbUser := range dbUsers {
        if params.Email == dbUser.Email {
            err = bcrypt.CompareHashAndPassword(dbUser.Password, []byte(params.Password))
            if err != nil {
                respondWithError(w, http.StatusUnauthorized, "Incorrect password")
                return
            }
            respondWithJSON(w, http.StatusOK, User {
                ID:         dbUser.ID,
                Email:      dbUser.Email,
            })
            return
        }
    }
    respondWithError(w, http.StatusUnauthorized, "User does not exist")
    return
}
