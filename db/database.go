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

func FetchHomepageRaceStats(db *sql.DB) (models.HomePageData, error) {
	var homepageStats models.HomePageData
	var err error

	homepageStats.LastRaces, err = fetchLastRaces(db)
	if err != nil {
		return homepageStats, err
	}

	// homepageStats.Upcoming, err = fetchUpcomingRaces(db)
	// if err != nil {
	// 	return homepageStats, err
	// }

	homepageStats.TopMen, err = fetchTop5(db, "M")
	if err != nil {
		return homepageStats, err
	}

	homepageStats.TopWomen, err = fetchTop5(db, "F")
	if err != nil {
		return homepageStats, err
	}

	homepageStats.TopJuniors, err = fetchTopJuniors(db)
	if err != nil {
		return homepageStats, err
	}

	homepageStats.News = fetchFakeNews() // or load from DB if you want

	return homepageStats, nil
}

func fetchLastRaces(db *sql.DB) ([]models.HomePageRace, error) {
	query := `
	SELECT r.id, r.name, r.date,
		(SELECT rider_id FROM results WHERE race_id = r.id AND position = 1 LIMIT 1) AS first_place,
		(SELECT CONCAT(first_name, ' ', last_name) FROM riders WHERE id = (SELECT rider_id FROM results WHERE race_id = r.id AND position = 1 LIMIT 1)) AS first_name,
		(SELECT rider_id FROM results WHERE race_id = r.id AND position = 2 LIMIT 1) AS second_place,
		(SELECT CONCAT(first_name, ' ', last_name) FROM riders WHERE id = (SELECT rider_id FROM results WHERE race_id = r.id AND position = 2 LIMIT 1)) AS second_name,
		(SELECT rider_id FROM results WHERE race_id = r.id AND position = 3 LIMIT 1) AS third_place,
		(SELECT CONCAT(first_name, ' ', last_name) FROM riders WHERE id = (SELECT rider_id FROM results WHERE race_id = r.id AND position = 3 LIMIT 1)) AS third_name
	FROM races r
	ORDER BY r.date DESC
	LIMIT 3`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var races []models.HomePageRace
	for rows.Next() {
		var race models.HomePageRace
		err := rows.Scan(
			&race.RaceID,
			&race.RaceName,
			&race.Date,
			&race.FirstPlace,
			&race.FirstPlaceName,
			&race.SecondPlace,
			&race.SecondPlaceName,
			&race.ThirdPlace,
			&race.ThirdPlaceName,
		)
		if err != nil {
			return nil, err
		}
		races = append(races, race)
	}
	return races, nil
}

