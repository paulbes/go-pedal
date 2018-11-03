package api

import "context"

// Note: this represents our API domain model
// I have decoupled this from the pedlar model
// so I can change them independenly, if required

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

// StationService defines what methods a station
// service implementation must implement
type StationService interface {
	Get(ctx context.Context, id int) (Station, error)
	List(ctx context.Context) ([]Station, error)
}

// StationStore defines what methods a station
// storage implementation must implement
type StationStore interface {
	Get(id int) (Station, error)
	List() ([]Station, error)
}
