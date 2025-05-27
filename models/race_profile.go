package models

type FullRaceProfile struct {
    Race    RaceDetails     `json:"race"`
    Results []RaceResultRow `json:"results"`
}

type RaceDetails struct {
    ID                int     `json:"id"`
    Date              string  `json:"date"`
    Name              string  `json:"name"`
    Category          string  `json:"category"`
    Location          string  `json:"location"`
    Distance          float64 `json:"distance"`
    Elevation         int     `json:"elevation"`
    Roughness         string  `json:"roughness,omitempty"`
    StartTime         string  `json:"start_time"`
    Temperature       int     `json:"temperature"`
    TotalParticipants int     `json:"totalParticipants"`
}

type RaceResultRow struct {
    Position    int    `json:"position"`
    RiderID     int    `json:"riderId"`
    RiderName   string `json:"riderName"`
    Team        string `json:"team"`
    Time        string `json:"time"`
    Points      int    `json:"points"`
}
