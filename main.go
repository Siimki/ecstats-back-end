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

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	dbConn := db.ConnectToDB(cfg)
	defer dbConn.Close()

	r:= chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handlers.NewHandler(dbConn)

	r.Get("/api/stats", h.GetSummary)
	r.Get("/api/homepagestats", h.GetHomepageRaceStats)
	r.Get("/api/riderprofile", h.GetRiderProfile)

	log.Println("Server started on: 1337")
	http.ListenAndServe(":1337", r)
}