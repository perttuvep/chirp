package main

import (
	"net/http"
	"time"

	"github.com/perttuvep/chirp/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	rtoken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Revoke:Error reading rtoken", err)
		return
	}
	rtokdb, err := cfg.DbQueries.GetRTokenByToken(r.Context(), rtoken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Revoke:Error reading rtoken from db", err)
		return
	}
	if !rtokdb.ExpiresAt.After(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Revoke:Refresh token expired!", err)
	}
	cfg.DbQueries.RevokeToken(r.Context(), rtokdb.Token)

	w.WriteHeader(204)
}
