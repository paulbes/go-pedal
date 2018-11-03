package server

import (
	"context"

	"github.com/paulbes/go-pedal/pkg/api"
)

// Note: the service layer implements the business logic of the application
// we currently don't have much of that though.

type stationService struct {
	store api.StationStore
}

func (s *stationService) Get(ctx context.Context, id int) (api.Station, error) {
	return s.store.Get(id)
}

func (s *stationService) List(ctx context.Context) ([]api.Station, error) {
	return s.store.List()
}

// NewStationService returns an initialised station service
func NewStationService(store api.StationStore) api.StationService {
	return &stationService{
		store: store,
	}
}
