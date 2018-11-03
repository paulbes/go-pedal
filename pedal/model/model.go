package model

import "time"

// Status describes whether a station is open
// or not
type Status struct {
	AllStationsClosed bool  `json:"all_stations_closed"`
	StationsClosed    []int `json:"stations_closed"`
}

// Availability describes how many locks
// or bikes are available at a given location
type Availability struct {
	Bikes int `json:"bikes"`
	Locks int `json:"locks"`
}

// Coord represents a lat and long coordinate
type Coord struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Stations represents all bike stations
type Stations struct {
	Stations []*Station `json:"stations"`
}

// Station represents a single bike station
type Station struct {
	ID            int          `json:"id"`
	InService     bool         `json:"in_service"`
	Title         string       `json:"title"`
	Subtitle      string       `json:"subtitle"`
	NumberOfLocks int          `json:"number_of_locks"`
	Center        Coord        `json:"center"`
	Bounds        []Coord      `json:"bounds"`
	Availability  Availability `json:"-"`
	Closed        bool         `json:"-"`
}

// StationAvailability represents the availability of
// bikes and locks, with refresh rate, etc.
type StationAvailability struct {
	Stations []struct {
		ID           int          `json:"ID"`
		Availability Availability `json:"availability"`
	} `json:"stations"`
	UpdatedAt   time.Time `json:"updated_at"`
	RefreshRate float32   `json:"refresh_rate"`
}
