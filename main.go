package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/perttuvep/chirp/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DbQueries      *database.Queries
	platform       string
	Secret         string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

const EXPIRE_TIME_SEC int = 3600
const REFRESH_EXPIRE_TIME_DAY int = 60

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	secret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("Error opening db %v", err)
	}

	const filepathRoot = "."
	const port = "8080"

	conf := apiConfig{
		fileserverHits: atomic.Int32{},
		Secret:         secret,
	}

	conf.platform = os.Getenv("PLATFORM")
	conf.DbQueries = database.New(db)

	mux := http.NewServeMux()
	fsHandler := conf.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/chirps", conf.handlerGetChirps)
	mux.HandleFunc("POST /api/chirps", conf.handlerNewChirp)
	mux.HandleFunc("GET /api/chirps/{id}", conf.handlerGetChirpById)
	mux.HandleFunc("POST /admin/reset", conf.handlerReset)
	mux.HandleFunc("GET /admin/metrics", conf.handlerMetrics)

	mux.HandleFunc("POST /api/login", conf.handlerLogin)
	mux.HandleFunc("POST /api/users", conf.handlerNewUser)
	mux.HandleFunc("PUT /api/users", conf.handlerEditUser)
	mux.HandleFunc("POST /api/refresh", conf.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", conf.handlerRevoke)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
