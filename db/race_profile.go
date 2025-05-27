package db

import (
    "database/sql"
    "ecstats-back-end/models"
    "log"
    _ "github.com/lib/pq"
)

func GetRaceProfile(db *sql.DB, raceID int) (models.FullRaceProfile, error) {
    var profile models.FullRaceProfile

    // Fetch race details (including weather/participants)
    race, err := fetchRaceDetails(db, raceID)
    if err != nil {
        return profile, err
    }
    profile.Race = race

    // Fetch race results (all finishers, teams, times, etc)
    results, err := fetchRaceResults(db, raceID)
    if err != nil {
        return profile, err
    }
    profile.Results = results

    return profile, nil
}

func fetchRaceDetails(db *sql.DB, raceID int) (models.RaceDetails, error) {
    var r models.RaceDetails
    // Update query as your schema grows
	query := `
    SELECT
        races.id,
		races.date,
        races.name,
        races.category,
        races.location,
        COALESCE(races.distance, 0) as distance,
        COALESCE(race_data.elevation, 0) as elevation,
        COALESCE(race_data.start_time::text, '') as start_time,
        COALESCE(race_data.temperature, 0) as temperature,
        (
          SELECT COUNT(*) FROM results WHERE race_id = $1
        ) as total_participants
    FROM races
    LEFT JOIN race_data ON races.id = race_data.race_id
    WHERE races.id = $1
    LIMIT 1;
`



    err := db.QueryRow(query, raceID).Scan(
        &r.ID,
		&r.Date,
        &r.Name,
        &r.Category,
        &r.Location,
        &r.Distance,
        &r.Elevation,
        &r.StartTime,
        &r.Temperature,
        &r.TotalParticipants,
    )
    return r, err
}

func fetchRaceResults(db *sql.DB, raceID int) ([]models.RaceResultRow, error) {
    var results []models.RaceResultRow

    query := `
    SELECT
        results.position,
        riders.id,
        CONCAT(riders.first_name, ' ', riders.last_name) as rider_name,
        COALESCE(teams.name, '') as team,
        COALESCE(results.time::text, '') as time,
        COALESCE(results.points, 0)
    FROM results
    JOIN riders ON results.rider_id = riders.id
    LEFT JOIN teams ON results.team_id = teams.id
    WHERE results.race_id = $1
    ORDER BY results.position ASC;
    `

    rows, err := db.Query(query, raceID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var r models.RaceResultRow
        err := rows.Scan(
            &r.Position,
            &r.RiderID,
            &r.RiderName,
            &r.Team,
            &r.Time,
            &r.Points,
        )
        if err != nil {
            log.Printf("[FetchRaceResults]Error executing query %q with raceID=%d: %v", query, raceID, err)
            return nil, err
        }
        results = append(results, r)
    }

    return results, nil
}
