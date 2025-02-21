package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: profanity(params.Body)})

}

func profanity(s string) string {
	out := strings.Split(s, " ")
	for i, v := range out {
		if strings.ToLower(v) == "kerfuffle" || strings.ToLower(v) == "sharbert" || strings.ToLower(v) == "fornax" {
			out[i] = "****"
		}
	}
	return strings.Join(out, " ")
}
func ContainsAt[T comparable](list []T, item T) (bool, int) {
	for i, v := range list {
		if v == item {
			return true, i
		}
	}
	return false, -1
}
