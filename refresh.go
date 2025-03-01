package main

import (
	"log"
	"net/http"
	"time"

	"github.com/perttuvep/chirp/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	rtoken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error reading rtoken", err)
		return
	}
	rtokdb, err := cfg.DbQueries.GetRTokenByToken(r.Context(), rtoken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error reading rtoken from db", err)
		return
	}
	if !rtokdb.RevokedAt.Time.IsZero() {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired!", err)
	}

	expiresInSeconds := time.Duration(EXPIRE_TIME_SEC) * time.Second

	token, err := auth.MakeJWT(rtokdb.UserID, cfg.Secret, expiresInSeconds)
	if err != nil {
		log.Print("error creating token")
		respondWithError(w, http.StatusBadRequest, "Error creating token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Response{Token: token})

}
