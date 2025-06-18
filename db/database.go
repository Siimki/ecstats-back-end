package db

import (
	"database/sql"
	"ecstats-back-end/config"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

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

