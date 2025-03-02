package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/perttuvep/chirp/internal/auth"
	"github.com/perttuvep/chirp/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Email string `json:"email"`
		Pass  string `json:"password"`
	}

	type userParams struct {
		Id            string    `json:"id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Token         string    `json:"token"`
		Refresh_token string    `json:"refresh_token"`
		IsChirpyRed   bool      `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	reqparam := reqParams{}

	err := decoder.Decode(&reqparam)

	if err != nil {
		log.Printf("Error decoding json %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.DbQueries.GetUserByEmail(r.Context(), reqparam.Email)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email of password", err)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(reqparam.Pass)); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email of password", err)
		log.Printf(user.HashedPass, reqparam.Pass)
		hashed, err := auth.HashPass(reqparam.Pass)
		log.Println(hashed, err)
		return
	}
	expiresInSeconds := time.Duration(EXPIRE_TIME_SEC) * time.Second
	token, err := auth.MakeJWT(user.ID, cfg.Secret, expiresInSeconds)
	if err != nil {
		log.Print("error creating token")
		respondWithError(w, http.StatusBadRequest, "Error creating token", err)
		return
	}

	refrtoken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error creating refresh token", err)
		return
	}
	arg := database.CreateRTokenParams{Token: refrtoken, UserID: user.ID, ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * time.Duration(REFRESH_EXPIRE_TIME_DAY)), RevokedAt: sql.NullTime{Valid: false}}
	dbrtoken, err := cfg.DbQueries.CreateRToken(r.Context(), arg)

	if err != nil {
		log.Print("Error creating refresh token db entry")
		respondWithError(w, http.StatusBadRequest, "refresh token db err", err)
		return
	}
	respondWithJSON(w, http.StatusOK, userParams{Id: user.ID.String(), CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Token: token, Refresh_token: dbrtoken.Token, IsChirpyRed: user.IsChirpyRed})
}
