package models

// RiderProfile holds personal and ranking info
type RiderProfile struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Team          string `json:"team"`
	Birthplace    string `json:"birthplace"`
	ActiveSince   int    `json:"activeSince"`
	RankingPointsRoad int    `json:"rankingPointsRoad"`
	RankingPlaceRoad  int    `json:"rankingPlaceRoad"`
	RankingPointsMTB int    `json:"rankingPointsMTB"`
	RankingPlaceMTB  int    `json:"rankingPlaceMTB"`
}

// RiderPartnerships lists gear and sponsor partners
type RiderPartnerships struct {
	Frame     string `json:"frame"`
	Tyres     string `json:"tyres"`
	Clothing  string `json:"clothing"`
	BikeShop  string `json:"bikeShop"`
	Sponsor   string `json:"sponsor"`
}

// RiderTotals is career-wide summary
type RiderTotals struct {
	Participations int `json:"participations"`
	Wins           int `json:"wins"`
	Podiums        int `json:"podiums"`
	CareerPoints   int `json:"careerPoints"`
}

// RiderSeasonStats holds stats for a single year
type RiderSeasonStats struct {
	Year          int `json:"year"`
	Races         int `json:"races"`
	Wins          int `json:"wins"`
	Podiums       int `json:"podiums"`
	SeasonPoints  int `json:"points"`
}

// RiderResult is one row in result table
type RiderResult struct {
	Season   int    `json:"season"`
	Date     string `json:"date"`
	Race     string `json:"race"`
	Category string `json:"category"`
	Position int    `json:"position"`
	Points   int    `json:"points"`
}

// RiderTopResult is simplified for top 3 highlight
type RiderTopResult struct {
	Race     string `json:"race"`
	Date     string `json:"date"`
	Position int    `json:"position"`
	Points   int    `json:"points"`
}

// FullRiderProfile wraps everything sent to frontend
type FullRiderProfile struct {
	Profile      RiderProfile          `json:"profile"`
	Partnerships RiderPartnerships     `json:"partnerships"`
	Totals       RiderTotals           `json:"totals"`
	SeasonStats  []RiderSeasonStats    `json:"seasonStats"`
	Results      []RiderResult         `json:"results"`
	TopResults   []RiderTopResult      `json:"topResults"`
}

type RiderRanking struct {
	TotalPoints int
	Place       int
  }