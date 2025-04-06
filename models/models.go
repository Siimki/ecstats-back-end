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