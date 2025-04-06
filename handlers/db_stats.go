package handlers

import (
	"encoding/json"
	"net/http"
	"ecstats-back-end/db"
	"database/sql"
	"fmt"
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

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
	}
}