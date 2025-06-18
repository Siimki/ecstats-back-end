package db

import (
	"database/sql"
	"ecstats-back-end/models"
	"ecstats-back-end/utils"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"time"
)

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
	start := time.Now() 

	var homepageStats models.HomePageData

	// Run all fetches in parallel
	racesCh, err1 := utils.RunAsync(func() ([]models.HomePageRace, error) {
		return fetchLastRaces(db)
	})
	
	menCh, err2 := utils.RunAsync(func() ([]models.HomePageRankingRider, error) {
		return fetchTop5(db, "M")
	})
	
	womenCh, err3 := utils.RunAsync(func() ([]models.HomePageRankingRider, error) {
		return fetchTop5(db, "F")
	})
	
	juniorsCh, err4 := utils.RunAsync(func() ([]models.HomePageRankingRider, error) {
		return fetchTopJuniors(db)
	})
	
	// Collect results or return first error
	for i := 0; i < 4; i++ {
		select {
		case err := <-err1:
			return homepageStats, err
		case err := <-err2:
			return homepageStats, err
		case err := <-err3:
			return homepageStats, err
		case err := <-err4:
			return homepageStats, err
	
		case homepageStats.LastRaces = <-racesCh:
		case homepageStats.TopMen = <-menCh:
		case homepageStats.TopWomen = <-womenCh:
		case homepageStats.TopJuniors = <-juniorsCh:
		}
	}
	
	duration := time.Since(start)
	log.Printf("FetchHomepageRaceStats took %s", duration)

	homepageStats.News = fetchFakeNews()
	
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

	var rawDate time.Time
	var races []models.HomePageRace
	for rows.Next() {
		var race models.HomePageRace
		err := rows.Scan(
			&race.RaceID,
			&race.RaceName,
			&rawDate,
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
		race.Date = rawDate.Format("January 2, 2006")
		races = append(races, race)
	}
	return races, nil
}

func fetchFakeNews() []models.NewsArticle {
	return []models.NewsArticle{
		{
			Header:  "Estonian Cycling Progress",
			Content: "In 2010, only 2 riders crossed 45km/h in TTs. By 2025, over 20 Estonian riders have achieved this.",
			Date:    "2025-05-15",
		},
		{
			Header:  "New National Team Selection Criteria",
			Content: "Estonian Cycling Federation now uses a combination of UCI and national points for team selection.",
			Date:    "2025-05-10",
		},
	}
}

func fetchTop5(db *sql.DB, gender string) ([]models.HomePageRankingRider, error) {
	query := `
	SELECT riders.id, CONCAT(riders.first_name, ' ', riders.last_name) AS name,
	  SUM(
	    CASE
	      WHEN results.position = 1 THEN 25
	      WHEN results.position = 2 THEN 20
	      WHEN results.position = 3 THEN 16
	      WHEN results.position = 4 THEN 13
	      WHEN results.position = 5 THEN 11
	      WHEN results.position = 6 THEN 10
	      WHEN results.position = 7 THEN 9
	      WHEN results.position = 8 THEN 8
	      WHEN results.position = 9 THEN 7
	      WHEN results.position = 10 THEN 6
	      WHEN results.position = 11 THEN 5
	      WHEN results.position = 12 THEN 4
	      WHEN results.position = 13 THEN 3
	      WHEN results.position = 14 THEN 2
	      WHEN results.position = 15 THEN 1
	      WHEN results.position BETWEEN 16 AND 20 THEN 1
	      ELSE 0
	    END *
	    CASE
	      WHEN races.category = 'A' THEN 1.0
	      WHEN races.category = 'B' THEN 0.75
	      WHEN races.category = 'C' THEN 0.5
	      ELSE 0.5
	    END
	  ) AS points
	FROM results
	JOIN riders ON results.rider_id = riders.id
	JOIN races ON results.race_id = races.id
	WHERE riders.gender = $1 AND EXTRACT(YEAR FROM races.date) = 2024
	GROUP BY riders.id, CONCAT(riders.first_name, ' ', riders.last_name)
	ORDER BY points DESC
	LIMIT 5`

	rows, err := db.Query(query, gender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var riders []models.HomePageRankingRider
	for rows.Next() {
		var r models.HomePageRankingRider
		err := rows.Scan(&r.RiderID, &r.Name, &r.Points)
		if err != nil {
			return nil, err
		}
		riders = append(riders, r)
	}
	return riders, nil
}

func fetchTopJuniors(db *sql.DB) ([]models.HomePageRankingRider, error) {
	query := `
	SELECT riders.id, CONCAT(riders.first_name, ' ', riders.last_name) AS name,
	  SUM(
	    CASE
	      WHEN results.position = 1 THEN 25
	      WHEN results.position = 2 THEN 20
	      WHEN results.position = 3 THEN 16
	      WHEN results.position = 4 THEN 13
	      WHEN results.position = 5 THEN 11
	      WHEN results.position = 6 THEN 10
	      WHEN results.position = 7 THEN 9
	      WHEN results.position = 8 THEN 8
	      WHEN results.position = 9 THEN 7
	      WHEN results.position = 10 THEN 6
	      WHEN results.position = 11 THEN 5
	      WHEN results.position = 12 THEN 4
	      WHEN results.position = 13 THEN 3
	      WHEN results.position = 14 THEN 2
	      WHEN results.position = 15 THEN 1
	      WHEN results.position BETWEEN 16 AND 20 THEN 1
	      ELSE 0
	    END *
	    CASE
	      WHEN races.category = 'A' THEN 1.0
	      WHEN races.category = 'B' THEN 0.75
	      WHEN races.category = 'C' THEN 0.5
	      ELSE 0.5
	    END
	  ) AS points
	FROM results
	JOIN riders ON results.rider_id = riders.id
	JOIN races ON results.race_id = races.id
	WHERE riders.birth_year >= EXTRACT(YEAR FROM CURRENT_DATE) - 18
	  AND EXTRACT(YEAR FROM races.date) = 2024
	GROUP BY riders.id, CONCAT(riders.first_name, ' ', riders.last_name)
	ORDER BY points DESC
	LIMIT 5`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var juniors []models.HomePageRankingRider
	for rows.Next() {
		var r models.HomePageRankingRider
		err := rows.Scan(&r.RiderID, &r.Name, &r.Points)
		if err != nil {
			return nil, err
		}
		juniors = append(juniors, r)
	}
	return juniors, nil
}

func FetchTop100Riders(db *sql.DB) ([]models.RidersTop100, error) {

	//It looks dirty
	
	query := `
	WITH top_riders AS (
		SELECT rr.rider_id, r.first_name, r.last_name, COALESCE(t.name, '') AS team, rr.total_points
		FROM rider_rankings rr
		JOIN riders r ON rr.rider_id = r.id
		LEFT JOIN rider_teams rt ON rt.rider_id = r.id AND rt.year = rr.year
		LEFT JOIN teams t ON rt.team_id = t.id
		WHERE rr.discipline = 'road' AND rr.year = 2024
		ORDER BY rr.total_points DESC
		LIMIT 100
	),
	race_results_dedup AS (
		SELECT
			res.rider_id,
			races.id AS race_id,
			races.name AS race_name,
			res.position,
			races.date,
			ROW_NUMBER() OVER (
				PARTITION BY res.rider_id, res.race_id
				ORDER BY races.date DESC
			) AS dup_rank
		FROM results res
		JOIN races ON races.id = res.race_id
		WHERE res.rider_id IN (SELECT rider_id FROM top_riders)
	),
	race_results_top3 AS (
		SELECT *,
			ROW_NUMBER() OVER (
				PARTITION BY rider_id
				ORDER BY date DESC
			) AS race_rank
		FROM race_results_dedup
		WHERE dup_rank = 1
	)
	SELECT
		tr.rider_id,
		tr.first_name,
		tr.last_name,
		tr.team,
		tr.total_points,
		rr.race_name,
		rr.position,
		rr.race_rank
	FROM top_riders tr
	LEFT JOIN race_results_top3 rr ON tr.rider_id = rr.rider_id AND rr.race_rank <= 3
	ORDER BY tr.total_points DESC, rr.race_rank ASC;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ridersMap := make(map[int]*models.RidersTop100)
	riderRaceSeen := make(map[int]map[string]bool)

	for rows.Next() {
		var (
			riderID                          int
			firstName, lastName, team       string
			points                           float64
			raceName, position               sql.NullString
			raceRank                         sql.NullInt32
		)

		if err := rows.Scan(&riderID, &firstName, &lastName, &team, &points, &raceName, &position, &raceRank); err != nil {
			return nil, err
		}

		if _, exists := ridersMap[riderID]; !exists {
			ridersMap[riderID] = &models.RidersTop100{
				RiderID:   riderID,
				FirstName: firstName,
				LastName:  lastName,
				Team:      team,
				Points:    points,
				LastRaces: []models.RiderRaceResult{},
			}
			riderRaceSeen[riderID] = make(map[string]bool)
		}

		if raceName.Valid && position.Valid {
			key := raceName.String + "|" + position.String
			if !riderRaceSeen[riderID][key] {
				riderRaceSeen[riderID][key] = true
				ridersMap[riderID].LastRaces = append(ridersMap[riderID].LastRaces, models.RiderRaceResult{
					RaceName: raceName.String,
					Position: position.String,
				})
			}
		}
	}

	// Convert to slice
	var result []models.RidersTop100
	for _, rider := range ridersMap {
		result = append(result, *rider)
	}

	return result, nil
}

func SearchRidersByNameOrTeam(db *sql.DB, q string) ([]models.RidersTop100, error) {
	query := `
	WITH matched_riders AS (
		SELECT
			r.id AS rider_id,
			r.first_name,
			r.last_name,
			COALESCE(t.name, '') AS team,
			rr.total_points
		FROM riders r
		LEFT JOIN rider_teams rt ON rt.rider_id = r.id AND rt.year = EXTRACT(YEAR FROM CURRENT_DATE)
		LEFT JOIN teams t ON rt.team_id = t.id
		LEFT JOIN rider_rankings rr ON rr.rider_id = r.id AND rr.year = EXTRACT(YEAR FROM CURRENT_DATE) AND rr.discipline = 'road'
		WHERE
			r.first_name ILIKE '%' || $1 || '%' OR
			r.last_name ILIKE '%' || $1 || '%' OR
			(r.first_name || ' ' || r.last_name) ILIKE '%' || $1 || '%' OR
			COALESCE(t.name, '') ILIKE '%' || $1 || '%'

		LIMIT 50
	),
	race_results_dedup AS (
		SELECT
			res.rider_id,
			races.name AS race_name,
			res.position,
			races.date,
			ROW_NUMBER() OVER (
				PARTITION BY res.rider_id, res.race_id
				ORDER BY races.date DESC
			) AS dup_rank
		FROM results res
		JOIN races ON races.id = res.race_id
		WHERE res.rider_id IN (SELECT rider_id FROM matched_riders)
	),
	race_results_top3 AS (
		SELECT *,
			ROW_NUMBER() OVER (
				PARTITION BY rider_id
				ORDER BY date DESC
			) AS race_rank
		FROM race_results_dedup
		WHERE dup_rank = 1
	)
	SELECT
		mr.rider_id,
		mr.first_name,
		mr.last_name,
		mr.team,
		COALESCE(mr.total_points, 0) AS total_points,
		rr.race_name,
		rr.position,
		rr.race_rank
	FROM matched_riders mr
	LEFT JOIN race_results_top3 rr ON mr.rider_id = rr.rider_id AND rr.race_rank <= 3
	ORDER BY mr.total_points DESC, rr.race_rank ASC;
	`

	rows, err := db.Query(query, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ridersMap := make(map[int]*models.RidersTop100)

	for rows.Next() {
		var (
			riderID                          int
			firstName, lastName, team       string
			points                           float64
			raceName, position               sql.NullString
			raceRank                         sql.NullInt32
		)

		err := rows.Scan(&riderID, &firstName, &lastName, &team, &points, &raceName, &position, &raceRank)
		if err != nil {
			return nil, err
		}

		if _, exists := ridersMap[riderID]; !exists {
			ridersMap[riderID] = &models.RidersTop100{
				RiderID:   riderID,
				FirstName: firstName,
				LastName:  lastName,
				Team:      team,
				Points:    points,
				LastRaces: []models.RiderRaceResult{},
			}
		}

		if raceName.Valid && position.Valid {
			ridersMap[riderID].LastRaces = append(ridersMap[riderID].LastRaces, models.RiderRaceResult{
				RaceName: raceName.String,
				Position: position.String,
			})
		}
	}

	var result []models.RidersTop100
	for _, rider := range ridersMap {
		result = append(result, *rider)
	}

	fmt.Println(result)
	return result, nil
}
