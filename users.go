package main

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/perttuvep/chirp/internal/auth"
	"github.com/perttuvep/chirp/internal/database"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Email string `json:"email"`
		Pass  string `json:"password"`
	}

	type userParams struct {
		Id            string    `json:"id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Is_Chirpy_red bool      `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	reqparam := reqParams{}

	err := decoder.Decode(&reqparam)

	if err != nil {
		log.Printf("Error decoding json %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(reqparam.Email) < 1 || len(reqparam.Pass) < 1 {
		log.Printf(reqparam.Email, reqparam.Pass)
		respondWithError(w, http.StatusUnauthorized, "Password or email missing!", nil)
		return
	}
	hash, err := auth.HashPass(reqparam.Pass)
	if err != nil {
		log.Printf("Failed to hash password %v", err)
		return
	}

	user, err := cfg.DbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:      reqparam.Email,
		HashedPass: hash,
	})

	if err != nil {
		log.Printf("Error creating user %v", err)
		respondWithError(w, http.StatusUnauthorized, "Error creating user %v", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, userParams{Id: user.ID.String(), CreatedAt: user.CreatedAt, UpdatedAt: user.CreatedAt, Email: user.Email})

}

func (cfg *apiConfig) handlerEditUser(w http.ResponseWriter, r *http.Request) {

	type reqParams struct {
		Email string `json:"email"`
		Pass  string `json:"password"`
	}

	type userParams struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	reqparam := reqParams{}

	err := decoder.Decode(&reqparam)

	if err != nil {
		log.Printf("Error decoding json %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(reqparam.Email) < 1 || len(reqparam.Pass) < 1 {
		log.Printf(reqparam.Email, reqparam.Pass)
		respondWithError(w, http.StatusUnauthorized, "Password or email missing!", nil)
		return
	}
	atoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting token %v", err)
		return
	}
	uuid, err := auth.ValidateJWT(atoken, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token %v", err)
		return
	}
	dbuser, err := cfg.DbQueries.GetUserByID(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting user %v", err)
		return
	}
	hashed, err := auth.HashPass(reqparam.Pass)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", nil)
		return
	}

	editargs := database.EditUserEmailParams{ID: dbuser.ID, Email: reqparam.Email, HashedPass: hashed}
	updatedUser, err := cfg.DbQueries.EditUserEmail(r.Context(), editargs)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error editing user %v", err)
		return
	}

	respondWithJSON(w, http.StatusOK, userParams{ID: updatedUser.ID, Email: updatedUser.Email, CreatedAt: updatedUser.CreatedAt, UpdatedAt: updatedUser.UpdatedAt, IsChirpyRed: updatedUser.IsChirpyRed})
}
