package pedal

import (
	"log"
	"time"
)

// Pedlar defines the available methods for
// interacting with the Oslo City Bike API
type Pedlar interface {
	Stations() (map[int]*Station, error)
}

// Client defines the interface that
// an API client must implement
type Client interface {
	Stations() (*Stations, error)
	Availability() (*StationAvailability, error)
	Status() (*Status, error)
}

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

// Stations represents a single bike station
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

// pedlar contains some basic data that
// is required to load oslo city bike data
type pedlar struct {
	client      Client
	stations    map[int]*Station
	refreshRate time.Duration
	lastUpdate  time.Time
}

// New creates a new client for interacting
// with the Oslo City Bike API
func New(client Client) Pedlar {
	return &pedlar{
		client:      client,
		stations:    map[int]*Station{},
		refreshRate: 0 * time.Second,
		lastUpdate:  time.Now().Add(-1 * time.Hour),
	}
}

// Stations returns a list of all stations
// and their availability and status
func (p *pedlar) Stations() (map[int]*Station, error) {
	err := p.populateStations()
	if err != nil {
		return nil, err
	}
	return p.stations, nil
}

func (p *pedlar) populateStations() error {
	stations, err := p.client.Stations()
	if err != nil {
		return err
	}
	updateAvailability := false
	for _, station := range stations.Stations {
		if _, hasKey := p.stations[station.ID]; !hasKey {
			p.stations[station.ID] = station
			updateAvailability = true
		}
	}

	// Determine if we should update the availability of bikes and locks
	if p.lastUpdate.Add(p.refreshRate).After(time.Now()) || updateAvailability {
		availability, err := p.client.Availability()
		if err != nil {
			return err
		}
		p.refreshRate = time.Duration(availability.RefreshRate) * time.Second
		p.lastUpdate = availability.UpdatedAt
		for _, station := range availability.Stations {
			if s, hasKey := p.stations[station.ID]; hasKey {
				s.Availability = station.Availability
			} else {
				log.Printf("availability: could not find station with ID: %d, skipping", station.ID)
			}
		}

		// While we are first updating, might as well get status
		status, err := p.client.Status()
		if err != nil {
			return err
		}
		if status.AllStationsClosed {
			for _, station := range p.stations {
				station.Closed = false
			}
		} else {
			for _, closedID := range status.StationsClosed {
				if station, hasKey := p.stations[closedID]; hasKey {
					station.Closed = true
				} else {
					log.Printf("close station: could not find station with id: %d", closedID)
				}
			}
		}
	}

	return nil
}
