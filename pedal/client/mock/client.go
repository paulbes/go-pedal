package mock

import (
	cli "github.com/paulbes/go-pedal/pedal/client"
	"github.com/paulbes/go-pedal/pedal/model"
)

type client struct {
	StationsFn     func() (*model.Stations, error)
	AvailabilityFn func() (*model.StationAvailability, error)
	StatusFn       func() (*model.Status, error)
}

// Stations returns the output of the mocked stations function
func (c *client) Stations() (*model.Stations, error) {
	return c.StationsFn()
}

// Availability returns the output of the mocked availability function
func (c *client) Availability() (*model.StationAvailability, error) {
	return c.AvailabilityFn()
}

// Status returns the output of the mocked status function
func (c *client) Status() (*model.Status, error) {
	return c.StatusFn()
}

// NewClient creates a new mock that returns the provided arguments
func NewClient(stations *model.Stations, availability *model.StationAvailability, status *model.Status, err error) cli.Client {
	return &client{
		StationsFn: func() (*model.Stations, error) {
			return stations, err
		},
		AvailabilityFn: func() (*model.StationAvailability, error) {
			return availability, err
		},
		StatusFn: func() (*model.Status, error) {
			return status, err
		},
	}
}
