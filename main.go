package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ecstats-back-end/config"
	"ecstats-back-end/db"
	"ecstats-back-end/handlers"
)

func main() {

	// cfg, err := config.LoadConfig("config/config.yaml")
	// if err != nil {
	// 	log.Fatal("Failed to load config:", err)
	// }
	cfg := config.LoadConfigFromEnv()

	dbConn := db.ConnectToDB(cfg)
	defer dbConn.Close()

	r:= chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	h := handlers.NewHandler(dbConn)

	r.Get("/api/stats", h.GetSummary)
	r.Get("/api/homepagestats", h.GetHomepageRaceStats)
	r.Get("/api/riderprofile", h.GetRiderProfile)
	r.Get("/api/raceprofile", h.GetRaceProfile)
	r.Get("/api/top100riders", h.GetTop100Riders)
	r.Get("/api/riders", h.SearchRiders)

	

	log.Println("Server started on: 1337")
	http.ListenAndServe(":1337", r)
}