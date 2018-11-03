package http

import (
	"fmt"
	"testing"

	mock3 "github.com/paulbes/go-pedal/pedal/model/mock"

	"github.com/paulbes/go-pedal/pedal"
	"github.com/paulbes/go-pedal/pedal/client/mock"
	"github.com/paulbes/go-pedal/pkg/api"
	mock2 "github.com/paulbes/go-pedal/pkg/api/mock"
	"github.com/stretchr/testify/assert"
)

func TestStationStore_Get(t *testing.T) {
	testCases := []struct {
		Name      string
		ID        int
		Err       error
		ExpectErr bool
		Expect    interface{}
	}{
		{
			Name:      "Get station",
			ID:        1,
			ExpectErr: false,
			Expect:    mock2.NewStation(),
		},
		{
			Name:      "Get station, bad id",
			ID:        1000,
			ExpectErr: true,
			Expect:    "notfound: could not find station: no such id: 1000",
		},
		{
			Name:      "Get station, storage error",
			ID:        1,
			Err:       fmt.Errorf("could not connect to API"),
			ExpectErr: true,
			Expect:    "io: failed to read station: could not connect to API",
		},
	}

	for _, tc := range testCases {
		client := mock.NewClient(mock3.NewStations(), mock3.NewStationAvailability(), mock3.NewStatus(), tc.Err)
		store := NewStationStore(pedal.New(client))
		got, err := store.Get(tc.ID)
		if tc.ExpectErr {
			assert.Equal(t, err.Error(), tc.Expect)
		} else {
			assert.Nil(t, err, tc.Name)
			assert.Equal(t, tc.Expect, got)
		}
	}
}

func TestStationStore_List(t *testing.T) {
	testCases := []struct {
		Name      string
		Err       error
		ExpectErr bool
		Expect    interface{}
	}{
		{
			Name:      "List stations",
			ExpectErr: false,
			Expect:    []api.Station{mock2.NewStation()},
		},
		{
			Name:      "List stations, storage error",
			Err:       fmt.Errorf("could not connect to API"),
			ExpectErr: true,
			Expect:    "io: failed to read stations: could not connect to API",
		},
	}

	for _, tc := range testCases {
		client := mock.NewClient(mock3.NewStations(), mock3.NewStationAvailability(), mock3.NewStatus(), tc.Err)
		store := NewStationStore(pedal.New(client))
		got, err := store.List()
		if tc.ExpectErr {
			assert.Equal(t, err.Error(), tc.Expect)
		} else {
			assert.Nil(t, err, tc.Name)
			assert.Equal(t, tc.Expect, got)
		}
	}
}
