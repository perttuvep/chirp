package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/perttuvep/chirp/internal/auth"
	"github.com/perttuvep/chirp/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerNewChirp(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	var request reqParams

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("Error decoding json  %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
	}

	if !ChirpsValidate(request.Body) {
		respondWithError(w, http.StatusBadRequest, "Chirp too long!", nil)
		return
	}
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Print("Error getting token")
		respondWithError(w, http.StatusUnauthorized, "Error in reading token", err)

		return
	}
	log.Print(token)
	valid, err := auth.ValidateJWT(token, cfg.Secret)
	log.Print(valid)
	if valid == uuid.Nil {
		log.Print("Validate fail")
		log.Print(token)
		respondWithError(w, http.StatusUnauthorized, "Error validating JWT", err)
		return

	}
	if err != nil {
		return
	}
	arg := database.CreateChirpParams{Body: profanity(request.Body), UserID: valid}
	chirp, err := cfg.DbQueries.CreateChirp(r.Context(), arg)
	if err != nil {
		log.Printf("Attempting to create chirp for user: %s", valid)
		log.Printf("Error creating chirp %v", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{valid, chirp.CreatedAt, chirp.UpdatedAt, chirp.Body, chirp.UserID})

}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	if r.PathValue("id") == "" {
		respondWithError(w, http.StatusBadRequest, "No id given!", nil)
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Printf("Error getting uuid %v", err)
		respondWithError(w, http.StatusBadRequest, "Chirp not found", nil)
		return
	}

	chirp, err := cfg.DbQueries.GetChirpById(r.Context(), id)
	if err != nil {
		log.Printf("Error getting chirp %v", err)
		respondWithError(w, http.StatusNotFound, "Chirp not found!", nil)
		return
	}
	apichirp := Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID}

	respondWithJSON(w, http.StatusOK, apichirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	var chirpsdb []database.Chirp
	chirpsdb, err := cfg.DbQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps %v", err)
		respondWithError(w, http.StatusNotFound, "Chirps not found!", nil)
		return
	}
	var chirpsapi []Chirp
	for _, v := range chirpsdb {
		chirpsapi = append(chirpsapi, Chirp{ID: v.ID, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt, Body: v.Body, UserID: v.UserID})
	}
	log.Printf("Wat")
	respondWithJSON(w, 200, chirpsapi)

}
