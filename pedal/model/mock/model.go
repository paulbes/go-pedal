package mock

import (
	"time"

	"github.com/paulbes/go-pedal/pedal/model"
)

// NewStation creates a mocked station
func NewStation() *model.Station {
	return &model.Station{
		ID:            1,
		InService:     true,
		Title:         "Antarctica",
		Subtitle:      "Close to the penguins",
		NumberOfLocks: 10,
		Center: model.Coord{
			Latitude:  59.00001,
			Longitude: 59.00002,
		},
		Bounds: []model.Coord{
			{
				Latitude:  59.10,
				Longitude: 59.09,
			},
		},
		Availability: model.Availability{
			Bikes: 5,
			Locks: 5,
		},
		Closed: false,
	}
}

// NewStations returns mocked stations
func NewStations() *model.Stations {
	return &model.Stations{
		Stations: []*model.Station{
			NewStation(),
		},
	}
}

// NewStationAvailability returns mocked station availability
func NewStationAvailability() *model.StationAvailability {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	return &model.StationAvailability{
		Stations: []struct {
			ID           int                `json:"ID"`
			Availability model.Availability `json:"availability"`
		}{
			{
				ID: 1,
				Availability: model.Availability{
					Bikes: 5,
					Locks: 5,
				},
			},
		},
		UpdatedAt:   t1,
		RefreshRate: 10.0,
	}
}

// NewStatus creates a mocked status object
func NewStatus() *model.Status {
	return &model.Status{
		AllStationsClosed: false,
		StationsClosed:    []int{},
	}
}
