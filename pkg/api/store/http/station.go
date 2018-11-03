package http

import (
	"fmt"

	"github.com/paulbes/go-pedal/pedal/model"

	"github.com/paulbes/go-pedal/pedal"
	"github.com/paulbes/go-pedal/pkg/api"
	"github.com/paulbes/go-pedal/pkg/errors"
)

// Note: here we implement the store layer, in this case that is an
// http layer towards pedlar.

type stationStore struct {
	pedlar pedal.Pedlar
}

// Get reads the stations from the pedlar client and returns
// the station that was requested
func (s *stationStore) Get(id int) (api.Station, error) {
	stations, err := s.pedlar.Stations()
	if err != nil {
		return api.Station{}, errors.New(err, "failed to read station", errors.IO)
	}

	station, hasKey := stations[id]
	if !hasKey {
		return api.Station{}, errors.New(fmt.Errorf("no such id: %d", id), "could not find station", errors.NotFound)
	}

	return convertStation(station), nil
}

// List reads the stations from the pedlar client and returns
// all the stations
func (s *stationStore) List() ([]api.Station, error) {
	stations, err := s.pedlar.Stations()
	if err != nil {
		return nil, errors.New(err, "failed to read stations", errors.IO)
	}
	var res []api.Station
	for _, station := range stations {
		res = append(res, convertStation(station))
	}
	return res, nil
}

// convertStation maps stations between the two domain
// models
func convertStation(station *model.Station) api.Station {
	res := api.Station{
		ID:            station.ID,
		InService:     station.InService,
		Title:         station.Title,
		Subtitle:      station.Subtitle,
		NumberOfLocks: station.NumberOfLocks,
		Center: api.Coord{
			Latitude:  station.Center.Latitude,
			Longitude: station.Center.Longitude,
		},
		Availability: api.Availability{
			Bikes: station.Availability.Bikes,
			Locks: station.Availability.Locks,
		},
		Closed: station.Closed,
	}
	for _, coord := range station.Bounds {
		res.Bounds = append(res.Bounds, api.Coord{
			Latitude:  coord.Latitude,
			Longitude: coord.Longitude,
		})
	}
	return res
}

// NewStationStore creates a new station store
func NewStationStore(pedlar pedal.Pedlar) api.StationStore {
	return &stationStore{
		pedlar: pedlar,
	}
}
