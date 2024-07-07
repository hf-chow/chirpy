package main

import (
	"encoding/json"
	"net/http"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int     `json:"id"`
	Email       string  `json:"email"`
    Password    string  `json:"password"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email       string      `json:"email"`
        Password    string      `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

    hashed_password, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
    
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't encrypt password")
        return 
    }

	user, err := cfg.DB.CreateUser(params.Email, hashed_password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
        ID:         user.ID,
		Email:      user.Email,
	})
}

