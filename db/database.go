package db

import (
	"database/sql"
	"ecstats-back-end/config"
	"ecstats-back-end/models"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)
var RiderMap = make(map[string]int)

func ConnectToDB(cfg *config.Config) *sql.DB {

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%d",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode,
		cfg.Database.Host, cfg.Database.Port,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("DB ping error:", err)
	}

	log.Println("Connected to database")
	return db

}

func FetchDbStats(db *sql.DB) (models.DbStats, error ) {
	var stats models.DbStats

	query := `
	SELECT
		(SELECT COUNT(*) FROM races) AS race_count,
		(SELECT COUNT(*) FROM results) AS result_count,
		(SELECT COUNT(*) FROM riders) AS rider_count;
	`

	err := db.QueryRow(query).Scan(&stats.RaceCount, &stats.ResultCount, &stats.RiderCount)
	if err != nil { 
		return stats, err
	}

	return stats, nil
}




