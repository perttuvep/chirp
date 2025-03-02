package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	// "github.com/google/uuid"
	"github.com/google/uuid"
	"github.com/perttuvep/chirp/internal/auth"
	// "github.com/perttuvep/chirp/internal/auth"
)

func (cfg *apiConfig) handlerWebHooks(w http.ResponseWriter, r *http.Request) {

	type Data struct {
		UserID string `json:"user_id,omitempty"`
	}
	type requestParameters struct {
		Event string `json:"event,omitempty"`
		Data  Data   `json:"data"`
	}
	req := requestParameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding json", err)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting token v", err)
		return
	}
	if apiKey != cfg.Polka_Key {
		respondWithError(w, http.StatusUnauthorized, "Invalid api key", err)
		return
	}
	// token, err := auth.GetBearerToken(r.Header.Clone())
	// uid, err := auth.ValidateJWT(token, cfg.Secret)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Error validating token ", err)
	// 	return
	// }

	if req.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	uid, err := uuid.Parse(req.Data.UserID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing uuid", nil)
		return
	}
	err = cfg.DbQueries.ChirpyRedEnableByID(r.Context(), uid)
	if err == sql.ErrNoRows {
		respondWithJSON(w, 404, nil)
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error enablind red %v", err)
	}
	respondWithJSON(w, 204, nil)
	return
}
