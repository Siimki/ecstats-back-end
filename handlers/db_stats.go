package handlers

import (
	"encoding/json"
	"net/http"
	"ecstats-back-end/db"
	"database/sql"
	"fmt"
	"strconv"
)

type Handler struct {
	DB *sql.DB
}


func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	dbStats, err := db.FetchDbStats(h.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching riders", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dbStats)
}

func (h *Handler) GetHomepageRaceStats(w http.ResponseWriter, r *http.Request) {
	homepageStats, err := db.FetchHomepageRaceStats(h.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching homepage stats", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(homepageStats)
}

func (h *Handler) GetRiderProfile(w http.ResponseWriter, r *http.Request) {
	riderIdStr := r.URL.Query().Get("riderId")
	yearStr := r.URL.Query().Get("year")

	riderId, err := strconv.Atoi(riderIdStr)
	if err != nil {
		http.Error(w, "Invalid riderId", http.StatusBadRequest)
		return
	}

	var year int
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "Invalid year", http.StatusBadRequest)
			return
		}
	}

	fullRiderProfile, err := db.GetRiderProfile(h.DB, riderId, year)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching rider profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(fullRiderProfile)
}

func (h *Handler) GetRaceProfile(w http.ResponseWriter, r *http.Request) {
	raceIdStr := r.URL.Query().Get("raceId")

	raceId, err := strconv.Atoi(raceIdStr)
	if err != nil {
		http.Error(w, "Invalid riderId", http.StatusBadRequest)
		return
	}

	homepageStats, err := db.GetRaceProfile(h.DB, raceId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching homepage stats", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(homepageStats)
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
	}
}