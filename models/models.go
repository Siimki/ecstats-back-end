package models

type Rider struct {
	LastName   string
	FirstName  string
	BirthYear  int
	Nationality string
	Gender      string
	Team 		string
}

type Team struct {
	Name string
	Year int
}

type Result struct {
	FirstName string
	LastName string
	BirthYear int
	RaceId int
	RiderId int
	Position int
	Time string
	BibNumber int
	Status string
	Points int
}

type DbStats struct {
	RaceCount   int `json:"race_count"`
	ResultCount int `json:"result_count"`
	RiderCount  int `json:"rider_count"`
}

type HomePageData struct {
	LastRaces   []HomePageRace        `json:"last_races"`
	Upcoming    []HomePageRace        `json:"upcoming_races"`
	TopMen      []HomePageRankingRider       `json:"top_men"`
	TopWomen    []HomePageRankingRider       `json:"top_women"`
	TopJuniors  []HomePageRankingRider       `json:"top_juniors"`
	News        []NewsArticle `json:"news"`
}

type HomePageRace struct {
	RaceID      string    `json:"race_id"`
	RaceName    string `json:"race_name"`
	FirstPlaceId       string    `json:"first_place_id"`
	FirstPlaceName   string `json:"first_place_name"`
	SecondPlaceId      string    `json:"second_place_id"`
	SecondPlaceName  string `json:"second_place_name"`
	ThirdPlaceId       string    `json:"third_place_id"`
	ThirdPlaceName   string `json:"third_place_name"`
	Date        string `json:"date"`
}


type NewsArticle struct {
	Header  string `json:"header"`
	Content string `json:"content"`
	Date    string `json:"date"`
}


type UpcomingRace struct {
	RaceID      string  `json:"race_id"`
	RaceName    string `json:"race_name"`
	Location    string `json:"location"`
	Category    string `json:"category"`
	Temperature int    `json:"temperature"`
}


type HomePageRankingRider struct {
	RiderID string    `json:"rider_id"`
	Name   string `json:"name"`
	Team   string `json:"team"`
	ProfilePhotoURL string `json:"profile_photo_url"`
	Points float64    `json:"points"`
}

type RiderRaceResult struct {
	RaceName string `json:"raceName"`
	Position string `json:"position"`
}

type RidersTop100 struct {
	RiderID    string             `json:"riderId"`
	FirstName  string             `json:"firstName"`
	LastName   string             `json:"lastName"`
	Team       string             `json:"team"`
	Points     float64            `json:"points"`
	LastRaces  []RiderRaceResult  `json:"lastRaces"`
}
