package db

import (
	"database/sql"
	"ecstats-back-end/models"
	"ecstats-back-end/utils"
	"log"
	"fmt"
	_ "github.com/lib/pq"
)

func GetRiderProfile(db *sql.DB, riderID int, year int) (models.FullRiderProfile, error) {
	var profile models.FullRiderProfile

	profileCh, errCh1 := utils.RunAsync(func() (models.RiderProfile, error) {
		return fetchRiderBaseInfo(db, riderID, year)
	})

	partnershipsCh, errCh2 := utils.RunAsync(func() (models.RiderPartnerships, error) {
		return fetchRiderPartnerships(db, riderID)
	})

	totalsCh, errCh3 := utils.RunAsync(func() (models.RiderTotals, error) {
		return fetchRiderTotals(db, riderID)
	})

	seasonStatsCh, errCh4 := utils.RunAsync(func() ([]models.RiderSeasonStats, error) {
		return fetchRiderSeasonStats(db, riderID, year)
	})

	resultsCh, errCh5 := utils.RunAsync(func() ([]models.RiderResult, error) {
		return fetchRiderResults(db, riderID)
	})

	topResultsCh, errCh6 := utils.RunAsync(func() ([]models.RiderTopResult, error) {
		return fetchRiderTopResults(db, riderID)
	})

	// rankingCh, errCh7 := utils.RunAsync(func() (models.RiderRanking, error) {
	// 	return fetchRiderRanking(db, riderID, year)
	// })

	for i := 0; i < 6; i++ {
		select {
		case err := <-errCh1:
			return profile, err
		case err := <-errCh2:
			return profile, err
		case err := <-errCh3:
			return profile, err
		case err := <-errCh4:
			return profile, err
		case err := <-errCh5:
			return profile, err
		case err := <-errCh6:
			return profile, err
		// case err := <-errCh7:
		// 	return profile, err

		case p := <-profileCh:
			profile.Profile = p
		case ps := <-partnershipsCh:
			profile.Partnerships = ps
		case t := <-totalsCh:
			profile.Totals = t
		case ss := <-seasonStatsCh:
			profile.SeasonStats = ss
		case rs := <-resultsCh:
			profile.Results = rs
		case tr := <-topResultsCh:
			profile.TopResults = tr
		// case r := <-rankingCh:
		// 	fmt.Println(r.TotalPoints, "hei")
		// 	fmt.Println(r.Place, "hei")
		// 	profile.Profile.RankingPointsMTB = r.TotalPointsMTB
		// 	profile.Profile.RankingPlaceMTB = r.PlaceMTB
		// 	profile.Profile.RankingPointsMTB = r.TotalPointsRoad
		// 	profile.Profile.RankingPlaceMTB = r.PlaceRoad
		// 	fmt.Println(profile.Profile)
		}
	}
	// fmt.Println(profile.Profile.RankingPoints, "Points")
	fmt.Println(profile.Profile)
	fmt.Println(profile)
	return profile, nil
}

func fetchRiderBaseInfo(db *sql.DB, riderID int, year int) (models.RiderProfile, error) {
	var rp models.RiderProfile

	query := `
	SELECT 
	r.id,
	CONCAT(r.first_name, ' ', r.last_name) AS name,
	COALESCE(t.name, '') AS team,
	r.nationality AS birthplace,
	COALESCE(MIN(ra.date), CURRENT_DATE) AS active_since,
	COALESCE(rr_road.total_points, 0) AS ranking_points_road,
	COALESCE(rr_road.ranking_place, 0) AS ranking_place_road,
	COALESCE(rr_mtb.total_points, 0) AS ranking_points_mtb,
	COALESCE(rr_mtb.ranking_place, 0) AS ranking_place_mtb
	FROM riders r
	LEFT JOIN rider_teams rt ON rt.rider_id = r.id AND rt.year = $2
	LEFT JOIN teams t ON t.id = rt.team_id
	LEFT JOIN results res ON res.rider_id = r.id
	LEFT JOIN races ra ON ra.id = res.race_id
	LEFT JOIN rider_rankings rr_road ON rr_road.rider_id = r.id AND rr_road.year = $2 AND rr_road.discipline = 'road'
	LEFT JOIN rider_rankings rr_mtb ON rr_mtb.rider_id = r.id AND rr_mtb.year = $2 AND rr_mtb.discipline = 'mtb'
	WHERE r.id = $1
	GROUP BY r.id, CONCAT(r.first_name, ' ', r.last_name), t.name, r.nationality, rr_road.total_points, rr_road.ranking_place, rr_mtb.total_points, rr_mtb.ranking_place
	LIMIT 1;
	`

	var activeSinceRaw sql.NullTime

	err := db.QueryRow(query, riderID, year).Scan(
		&rp.ID,
		&rp.Name,
		&rp.Team,
		&rp.Birthplace,
		&activeSinceRaw,
		&rp.RankingPointsRoad,
		&rp.RankingPlaceRoad,
		&rp.RankingPointsMTB,
		&rp.RankingPlaceMTB,
	)
	if err != nil {
		log.Printf("[FetchRiderBaseInfo]Error executing query %q with riderID=%d: %v", query, riderID, err)
		return rp, err
	}

	if activeSinceRaw.Valid {
		rp.ActiveSince = activeSinceRaw.Time.Year()
	} else {
		rp.ActiveSince = 0
	}

	return rp, nil
}

func fetchRiderPartnerships(db *sql.DB, riderID int) (models.RiderPartnerships, error) {
	var rp models.RiderPartnerships

	query := `
	SELECT frame, tyres, clothing, bike_shop, COALESCE(sponsor_1, '') || CASE WHEN sponsor_2 IS NOT NULL THEN ', ' || sponsor_2 ELSE '' END AS sponsor
	FROM rider_partnerships
	WHERE rider_id = $1 AND year = 2025
	LIMIT 1;
	`

	err := db.QueryRow(query, riderID).Scan(
		&rp.Frame,
		&rp.Tyres,
		&rp.Clothing,
		&rp.BikeShop,
		&rp.Sponsor,
	)

	if err == sql.ErrNoRows {
		// Return empty partnerships if none found
		return rp, nil
	}

	return rp, err
}

func fetchRiderTotals(db *sql.DB, riderID int) (models.RiderTotals, error) {
	var totals models.RiderTotals

	query := `
	SELECT
		COUNT(*) AS participations,
		COUNT(*) FILTER (WHERE position = 1) AS wins,
		COUNT(*) FILTER (WHERE position <= 3) AS podiums,
		COALESCE(SUM(points), 0) AS career_points
	FROM results
	WHERE rider_id = $1;
	`

	err := db.QueryRow(query, riderID).Scan(
		&totals.Participations,
		&totals.Wins,
		&totals.Podiums,
		&totals.CareerPoints,
	)

	return totals, err
}

func fetchRiderSeasonStats(db *sql.DB, riderID int, year int) ([]models.RiderSeasonStats, error) {
	var stats []models.RiderSeasonStats

	query := `
	SELECT
	  res.year,
	  COUNT(*) AS races,
	  COUNT(*) FILTER (WHERE res.position = 1) AS wins,
	  COUNT(*) FILTER (WHERE res.position <= 3) AS podiums,
	  COALESCE(rr.total_points, 0) AS points
	FROM (
	  SELECT
	    results.*,
	    EXTRACT(YEAR FROM races.date) AS year
	  FROM results
	  JOIN races ON results.race_id = races.id
	  WHERE results.rider_id = $1
	) res
	LEFT JOIN rider_rankings rr
	  ON rr.rider_id = $1 AND rr.year = res.year
	GROUP BY res.year, rr.total_points
	ORDER BY res.year DESC;
	`


	rows, err := db.Query(query, riderID)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.RiderSeasonStats
		err := rows.Scan(&s.Year, &s.Races, &s.Wins, &s.Podiums, &s.SeasonPoints)
		if err != nil {
			log.Printf("[FetchRiderSeasonStats]Error executing query %q with riderID=%d: %v", query, riderID, err)
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}

func fetchRiderResults(db *sql.DB, riderID int) ([]models.RiderResult, error) {
	var results []models.RiderResult

	query := `
	SELECT 
		EXTRACT(YEAR FROM races.date) AS season,
		race_id,
		TO_CHAR(races.date, 'YYYY-MM-DD') AS date,
		races.name AS race,
		races.category,
		races.race_type,
		results.position,
		COALESCE(results.points, 0) AS points
	FROM results
	JOIN races ON results.race_id = races.id
	WHERE results.rider_id = $1
	ORDER BY races.date DESC;
	`
	
	rows, err := db.Query(query, riderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RiderResult
		err := rows.Scan(&r.Season, &r.RaceId, &r.Date, &r.Race, &r.Category, &r.Type, &r.Position, &r.Points)
		if err != nil {
			log.Printf("[FetchRiderResults]Error executing query %q with riderID=%d: %v", query, riderID, err)
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}


func fetchRiderTopResults(db *sql.DB, riderID int) ([]models.RiderTopResult, error) {
	var topResults []models.RiderTopResult

	query := `
	SELECT 
		races.name AS race,
		TO_CHAR(races.date, 'YYYY-MM-DD') AS date,
		results.position,
		COALESCE(results.points, 0) AS points
	FROM results
	JOIN races ON results.race_id = races.id
	WHERE results.rider_id = $1
	ORDER BY results.points DESC, races.date DESC
	LIMIT 3;
	`

	rows, err := db.Query(query, riderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tr models.RiderTopResult
		err := rows.Scan(&tr.Race, &tr.Date, &tr.Position, &tr.Points)
		if err != nil {
			log.Printf("[FetchRiderTopResults]Error executing query %q with riderID=%d: %v", query, riderID, err)
			return nil, err
		}
		topResults = append(topResults, tr)
	}

	return topResults, nil
}

func fetchRiderRanking(db *sql.DB, riderID int, year int) (models.RiderRanking, error) {
	var ranking models.RiderRanking
	
	query := `
	SELECT total_points, ranking_place
	FROM rider_rankings
	WHERE rider_id = $1 AND year = $2
	LIMIT 1;
	`

	err := db.QueryRow(query, riderID, year).Scan(&ranking.TotalPoints, &ranking.Place)
	if err == sql.ErrNoRows {
		// No ranking found, return zeros
		return ranking, nil
	}
	return ranking, err
}
