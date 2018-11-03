package mock

import (
	"context"

	"github.com/paulbes/go-pedal/pkg/api"
)

// NewStation creates a mocked station
func NewStation() api.Station {
	return api.Station{
		ID:            1,
		InService:     true,
		Title:         "Antarctica",
		Subtitle:      "Close to the penguins",
		NumberOfLocks: 10,
		Center: api.Coord{
			Latitude:  59.00001,
			Longitude: 59.00002,
		},
		Bounds: []api.Coord{
			{
				Latitude:  59.10,
				Longitude: 59.09,
			},
		},
		Availability: api.Availability{
			Bikes: 5,
			Locks: 5,
		},
		Closed: false,
	}
}

type stationStore struct {
	GetFn  func(id int) (api.Station, error)
	ListFn func() ([]api.Station, error)
}

// Get returns the values of the mocked function
func (s *stationStore) Get(id int) (api.Station, error) {
	return s.GetFn(id)
}

// List returns the values of the mocked function
func (s *stationStore) List() ([]api.Station, error) {
	return s.ListFn()
}

// NewStationStore creates a mocked station store using the provided
// input values
func NewStationStore(station api.Station, err error) api.StationStore {
	return &stationStore{
		GetFn: func(int) (api.Station, error) {
			return station, err
		},
		ListFn: func() ([]api.Station, error) {
			return []api.Station{station}, err
		},
	}
}

type stationService struct {
	GetFn  func(ctx context.Context, id int) (api.Station, error)
	ListFn func(ctx context.Context) ([]api.Station, error)
}

// Get returns the value of the mocked function
func (s *stationService) Get(ctx context.Context, id int) (api.Station, error) {
	return s.GetFn(ctx, id)
}

// List returns the value of the mocked function
func (s *stationService) List(ctx context.Context) ([]api.Station, error) {
	return s.ListFn(ctx)
}

// NewStationService creates a mocked station service using the provided
// inputs values
func NewStationService(station api.Station, err error) api.StationService {
	return &stationService{
		GetFn: func(context.Context, int) (api.Station, error) {
			return station, err
		},
		ListFn: func(context.Context) ([]api.Station, error) {
			return []api.Station{station}, err
		},
	}
}
