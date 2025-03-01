package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform == "dev" {
		w.WriteHeader(http.StatusOK)
		err := cfg.DbQueries.ResetUsers(r.Context())
		if err != nil {
			log.Printf("reset users %v", err)
		}
		w.Write([]byte("Users&chirps reset"))
	} else {
		w.WriteHeader(http.StatusForbidden)
	}

}
