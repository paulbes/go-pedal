package pedal_test

import (
	"fmt"
	"testing"

	"github.com/paulbes/go-pedal/pedal"
	"github.com/paulbes/go-pedal/pedal/client/mock"
	"github.com/paulbes/go-pedal/pedal/model"
	modmock "github.com/paulbes/go-pedal/pedal/model/mock"
	"github.com/stretchr/testify/assert"
)

func TestPedlar_Stations(t *testing.T) {
	testCases := []struct {
		Name         string
		Stations     *model.Stations
		Availability *model.StationAvailability
		Status       *model.Status
		Err          error
		ExpectErr    bool
		Expect       interface{}
	}{
		{
			Name:         "Getting stations works",
			Stations:     modmock.NewStations(),
			Availability: modmock.NewStationAvailability(),
			Status:       modmock.NewStatus(),
			Expect: map[int]*model.Station{
				1: modmock.NewStation(),
			},
		},
		{
			Name:         "With error fails",
			Stations:     modmock.NewStations(),
			Availability: modmock.NewStationAvailability(),
			Status:       modmock.NewStatus(),
			Err:          fmt.Errorf("nope"),
			ExpectErr:    true,
			Expect:       "nope",
		},
	}

	for _, tc := range testCases {
		client := mock.NewClient(tc.Stations, tc.Availability, tc.Status, tc.Err)
		p := pedal.New(client)
		got, err := p.Stations()
		if tc.ExpectErr {
			assert.Equal(t, tc.Expect, err.Error(), tc.Name)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.Expect, got, tc.Name)
		}
	}
}
