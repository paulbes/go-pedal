package pedal

import (
	"fmt"
	"testing"
	"time"

	modmock "github.com/paulbes/go-pedal/pedal/model/mock"

	"github.com/paulbes/go-pedal/pedal/client/mock"
	"github.com/paulbes/go-pedal/pedal/model"
	"github.com/stretchr/testify/assert"
)

func TestPedlar_doPopulateStations(t *testing.T) {
	testCases := []struct {
		Name      string
		Stations  *model.Stations
		Initial   map[int]*model.Station
		Err       error
		ExpectErr bool
		Expect    interface{}
		ExpectNew bool
	}{
		{
			Name:      "Populating works",
			Stations:  &model.Stations{Stations: []*model.Station{{ID: 1}}},
			Expect:    map[int]*model.Station{1: {ID: 1}},
			Initial:   map[int]*model.Station{},
			ExpectNew: true,
		},
		{
			Name:      "With err",
			Err:       fmt.Errorf("nope"),
			ExpectErr: true,
			Expect:    "nope",
		},
		{
			Name:      "No new",
			Stations:  &model.Stations{Stations: []*model.Station{{ID: 1}}},
			Expect:    map[int]*model.Station{1: {ID: 1}},
			Initial:   map[int]*model.Station{1: {ID: 1}},
			ExpectNew: false,
		},
	}

	for _, tc := range testCases {
		client := mock.NewClient(tc.Stations, nil, nil, tc.Err)
		p := pedlar{
			client: client,
		}
		gotNew, got, err := p.doPopulateStations(tc.Initial)
		if tc.ExpectErr {
			assert.Equal(t, tc.Expect, err.Error(), tc.Name)
		} else {
			assert.Nil(t, err, tc.Name)
			assert.Equal(t, tc.ExpectNew, gotNew)
			assert.Equal(t, tc.Expect, got)
		}
	}
}

func TestPedlar_doUpdateAvailability(t *testing.T) {
	testCases := []struct {
		Name         string
		Initial      map[int]*model.Station
		Force        bool
		LastUpdate   time.Time
		Availability *model.StationAvailability
		Err          error
		ExpectErr    bool
		Expect       interface{}
		ExpectUpdate bool
	}{
		{
			Name:         "Updating works",
			Availability: modmock.NewStationAvailability(),
			LastUpdate:   time.Now().Add(-10 * time.Second),
			Expect:       map[int]*model.Station{1: modmock.NewStation()},
			Initial: map[int]*model.Station{
				1: func() *model.Station {
					s := modmock.NewStation()
					s.Availability = model.Availability{}
					return s
				}(),
			},
			ExpectUpdate: true,
		},
		{
			Name:      "With err",
			Err:       fmt.Errorf("nope"),
			ExpectErr: true,
			Expect:    "nope",
		},
		{
			Name:         "Forced",
			Availability: modmock.NewStationAvailability(),
			LastUpdate:   time.Now().Add(20 * time.Second),
			Force:        true,
			Expect:       map[int]*model.Station{1: modmock.NewStation()},
			Initial: map[int]*model.Station{
				1: func() *model.Station {
					s := modmock.NewStation()
					s.Availability = model.Availability{}
					return s
				}(),
			},
			ExpectUpdate: true,
		},
		{
			Name:         "No update",
			Availability: modmock.NewStationAvailability(),
			LastUpdate:   time.Now().Add(20 * time.Second),
			Expect: map[int]*model.Station{
				1: func() *model.Station {
					s := modmock.NewStation()
					s.Availability = model.Availability{}
					return s
				}(),
			},
			Initial: map[int]*model.Station{
				1: func() *model.Station {
					s := modmock.NewStation()
					s.Availability = model.Availability{}
					return s
				}(),
			},
			ExpectUpdate: true,
		},
	}

	for _, tc := range testCases {
		client := mock.NewClient(nil, tc.Availability, nil, tc.Err)
		p := pedlar{
			client:     client,
			lastUpdate: tc.LastUpdate,
		}
		updated, got, err := p.doUpdateAvailability(tc.Force, tc.Initial)
		if tc.ExpectErr {
			assert.Equal(t, tc.Expect, err.Error(), tc.Name)
		} else {
			assert.Nil(t, err, tc.Name)
			assert.Equal(t, tc.ExpectUpdate, updated)
			assert.Equal(t, tc.Expect, got)
		}
	}
}

func TestPedlar_doUpdateStatus(t *testing.T) {
	testCases := []struct {
		Name      string
		Status    *model.Status
		Initial   map[int]*model.Station
		Err       error
		ExpectErr bool
		Expect    interface{}
	}{
		{
			Name:    "Closing all works",
			Status:  &model.Status{AllStationsClosed: true},
			Expect:  map[int]*model.Station{1: {ID: 1, Closed: true}},
			Initial: map[int]*model.Station{1: {ID: 1}},
		},
		{
			Name:      "With err",
			Err:       fmt.Errorf("nope"),
			ExpectErr: true,
			Expect:    "nope",
		},
		{
			Name:    "Closing one works",
			Status:  &model.Status{AllStationsClosed: false, StationsClosed: []int{1}},
			Expect:  map[int]*model.Station{1: {ID: 1, Closed: true}},
			Initial: map[int]*model.Station{1: {ID: 1}},
		},
		{
			Name:    "Closing none works",
			Status:  &model.Status{AllStationsClosed: false, StationsClosed: []int{}},
			Expect:  map[int]*model.Station{1: {ID: 1}},
			Initial: map[int]*model.Station{1: {ID: 1}},
		},
	}

	for _, tc := range testCases {
		client := mock.NewClient(nil, nil, tc.Status, tc.Err)
		p := pedlar{
			client: client,
		}
		got, err := p.doUpdateStatus(tc.Initial)
		if tc.ExpectErr {
			assert.Equal(t, tc.Expect, err.Error(), tc.Name)
		} else {
			assert.Nil(t, err, tc.Name)
			assert.Equal(t, tc.Expect, got, tc.Name)
		}
	}
}
