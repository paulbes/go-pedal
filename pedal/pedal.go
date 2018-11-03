package pedal

import (
	"log"
	"time"

	"github.com/paulbes/go-pedal/pedal/client"
	"github.com/paulbes/go-pedal/pedal/model"
)

// Pedlar defines the available methods for
// interacting with the Oslo City Bike API
type Pedlar interface {
	Stations() (map[int]*model.Station, error)
}

// pedlar contains some basic data that
// is required to load oslo city bike data
type pedlar struct {
	client      client.Client
	stations    map[int]*model.Station
	refreshRate time.Duration
	lastUpdate  time.Time
}

// New creates a new client for interacting
// with the Oslo City Bike API
func New(client client.Client) Pedlar {
	return &pedlar{
		client:      client,
		stations:    map[int]*model.Station{},
		refreshRate: 0 * time.Second,
		lastUpdate:  time.Now().Add(-1 * time.Hour),
	}
}

// Stations returns a list of all stations
// and their availability and status
func (p *pedlar) Stations() (map[int]*model.Station, error) {
	// Populate the stations, if we have new stations,
	// lets force an update
	newStations, stations, err := p.doPopulateStations(p.stations)
	if err != nil {
		return nil, err
	}

	updated, stations, err := p.doUpdateAvailability(newStations, stations)
	if err != nil {
		return nil, err
	}

	if updated {
		stations, err = p.doUpdateStatus(stations)
		if err != nil {
			return nil, err
		}
	}

	p.stations = stations
	return p.stations, nil
}

func (p *pedlar) doPopulateStations(s map[int]*model.Station) (bool, map[int]*model.Station, error) {
	stations, err := p.client.Stations()
	if err != nil {
		return false, nil, err
	}

	newStations := false
	for _, station := range stations.Stations {
		if _, hasKey := s[station.ID]; !hasKey {
			s[station.ID] = station
			newStations = true
		}
	}

	return newStations, s, nil
}

func (p *pedlar) doUpdateAvailability(force bool, s map[int]*model.Station) (bool, map[int]*model.Station, error) {
	// Determine if we should update the availability of bikes and locks
	if p.lastUpdate.Add(p.refreshRate).Before(time.Now()) || force {
		availability, err := p.client.Availability()
		if err != nil {
			return false, nil, err
		}
		p.refreshRate = time.Duration(availability.RefreshRate) * time.Second
		p.lastUpdate = availability.UpdatedAt
		for _, station := range availability.Stations {
			if s, hasKey := s[station.ID]; hasKey {
				s.Availability = station.Availability
			} else {
				log.Printf("availability: could not find station with ID: %d, skipping", station.ID)
			}
		}
	}

	return true, s, nil
}

func (p *pedlar) doUpdateStatus(s map[int]*model.Station) (map[int]*model.Station, error) {
	status, err := p.client.Status()
	if err != nil {
		return s, err
	}
	if status.AllStationsClosed {
		for _, station := range s {
			station.Closed = true
		}
	} else {
		for _, closedID := range status.StationsClosed {
			if station, hasKey := s[closedID]; hasKey {
				station.Closed = true
			} else {
				log.Printf("close station: could not find station with id: %d", closedID)
			}
		}
	}

	return s, nil
}
