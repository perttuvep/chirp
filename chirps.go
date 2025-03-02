package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
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
	respondWithJSON(w, http.StatusCreated, Chirp{chirp.ID, chirp.CreatedAt, chirp.UpdatedAt, chirp.Body, chirp.UserID})

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
	var err error
	s := r.URL.Query().Get("author_id")
	sortQuery := r.URL.Query().Get("sort")

	if s == "" {
		chirpsdb, err = cfg.DbQueries.GetChirps(r.Context())
		if err != nil {
			log.Printf("Error getting chirps %v", err)
			respondWithError(w, http.StatusNotFound, "Chirps not found!", nil)
			return
		}

	} else {
		uid, err := uuid.Parse(s)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error parsing uid", err)
			return
		}
		chirpsdb, err = cfg.DbQueries.GetChirpsByUID(r.Context(), uid)
	}

	var chirpsapi []Chirp
	for _, v := range chirpsdb {
		chirpsapi = append(chirpsapi, Chirp{ID: v.ID, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt, Body: v.Body, UserID: v.UserID})
	}
	if strings.ToLower(sortQuery) == "desc" {
		sort.Slice(chirpsapi, func(i, j int) bool {
			return chirpsapi[j].CreatedAt.Before(chirpsapi[i].CreatedAt)
		})
	}
	respondWithJSON(w, 200, chirpsapi)

}
func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	atoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization token", err)
		return
	}
	userid, err := auth.ValidateJWT(atoken, cfg.Secret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization token", err)
		return
	}
	if r.PathValue("id") == "" {
		respondWithError(w, http.StatusBadRequest, "No id given!", nil)
	}
	chirp_id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Printf("Error getting uuid %v", err)
		respondWithError(w, http.StatusBadRequest, "Chirp not found", nil)
		return
	}
	chirp, err := cfg.DbQueries.GetChirpById(r.Context(), chirp_id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
		return
	}
	if chirp.UserID != userid {
		respondWithError(w, http.StatusForbidden, "", nil)
		return
	}
	cfg.DbQueries.DeleteChirpById(r.Context(), chirp_id)
	respondWithJSON(w, http.StatusNoContent, nil)

}
