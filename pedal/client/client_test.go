package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/paulbes/go-pedal/pedal/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestHttpClient_Stations(t *testing.T) {
	cli, err := NewHTTPClient("myID", 1)
	assert.Nil(t, err)

	station := &model.Station{
		ID:            157,
		InService:     true,
		Title:         "Nylandsveien",
		Subtitle:      "mellom Norbygata og Urtegata",
		NumberOfLocks: 30,
		Center: model.Coord{
			Latitude:  59.91562,
			Longitude: 10.762248,
		},
		Bounds: []model.Coord{
			{
				Latitude:  59.915418602160436,
				Longitude: 10.762068629264832,
			},
		},
	}

	testCases := []struct {
		Name      string
		Mock      func()
		Expect    interface{}
		ExpectErr bool
	}{
		{
			Name: "Stations Ok",
			Mock: func() {
				gock.New(BaseURL).Get("stations").Reply(http.StatusOK).File("fixtures/get.stations.json")
			},
			Expect: &model.Stations{
				Stations: []*model.Station{station},
			},
		},
		{
			Name: "Stations malformed",
			Mock: func() {
				gock.New(BaseURL).Get("stations").Reply(http.StatusOK).File("fixtures/get.stations.malformed.json")
			},
			Expect:    "invalid character 'h' looking for beginning of value",
			ExpectErr: true,
		},
		{
			Name: "Stations bad request",
			Mock: func() {
				gock.New(BaseURL).Get("stations").Reply(http.StatusBadRequest).File("fixtures/get.stations.bad.json")
			},
			Expect:    "failed to invoke API, got error code: 400, reason: {\n  \"error\": \"something\"\n}",
			ExpectErr: true,
		},
	}

	for _, tc := range testCases {
		gock.Clean()
		tc.Mock()

		stations, err := cli.Stations()
		if tc.ExpectErr {
			assert.Equal(t, tc.Expect, err.Error(), tc.Name)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, stations, tc.Expect)
		}
	}
}

func TestHttpClient_Availability(t *testing.T) {
	cli, err := NewHTTPClient("myID", 1)
	assert.Nil(t, err)

	t1, err := time.Parse(time.RFC3339, "2018-11-03T15:17:00+00:00")
	assert.Nil(t, err)
	availability := &model.StationAvailability{
		Stations: []struct {
			ID           int                `json:"ID"`
			Availability model.Availability `json:"availability"`
		}{
			{
				ID: 177,
				Availability: model.Availability{
					Bikes: 0,
					Locks: 28,
				},
			},
		},
		UpdatedAt:   t1,
		RefreshRate: 10.0,
	}

	testCases := []struct {
		Name      string
		Mock      func()
		Expect    interface{}
		ExpectErr bool
	}{
		{
			Name: "Availability Ok",
			Mock: func() {
				gock.New(BaseURL).Get("stations/availability").Reply(http.StatusOK).File("fixtures/get.availability.json")
			},
			Expect: availability,
		},
		{
			Name: "Availability malformed",
			Mock: func() {
				gock.New(BaseURL).Get("stations/availability").Reply(http.StatusOK).File("fixtures/get.availability.malformed.json")
			},
			Expect:    "unexpected end of JSON input",
			ExpectErr: true,
		},
		{
			Name: "Availability bad request",
			Mock: func() {
				gock.New(BaseURL).Get("stations/availability").Reply(http.StatusBadRequest).File("fixtures/get.availability.bad.json")
			},
			Expect:    "failed to invoke API, got error code: 400, reason: {\n  \"error\": \"something\"\n}",
			ExpectErr: true,
		},
	}

	for _, tc := range testCases {
		gock.Clean()
		tc.Mock()
		avail, err := cli.Availability()
		if tc.ExpectErr {
			assert.Equal(t, tc.Expect, err.Error(), tc.Name)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, avail, tc.Expect)
		}
	}
}

func TestHttpClient_Status(t *testing.T) {
	cli, err := NewHTTPClient("myID", 1)
	assert.Nil(t, err)

	status := &model.Status{
		AllStationsClosed: false,
		StationsClosed: []int{
			100,
		},
	}

	testCases := []struct {
		Name      string
		Mock      func()
		Expect    interface{}
		ExpectErr bool
	}{
		{
			Name: "Status Ok",
			Mock: func() {
				gock.New(BaseURL).Get("status").Reply(http.StatusOK).File("fixtures/get.status.json")
			},
			Expect: status,
		},
		{
			Name: "Status malformed",
			Mock: func() {
				gock.New(BaseURL).Get("status").Reply(http.StatusOK).File("fixtures/get.status.malformed.json")
			},
			Expect:    "invalid character '}' after object key",
			ExpectErr: true,
		},
		{
			Name: "Status bad request",
			Mock: func() {
				gock.New(BaseURL).Get("status").Reply(http.StatusBadRequest).File("fixtures/get.status.bad.json")
			},
			Expect:    "failed to invoke API, got error code: 400, reason: {\n  \"error\": \"something\"\n}",
			ExpectErr: true,
		},
	}

	for _, tc := range testCases {
		gock.Clean()
		tc.Mock()
		avail, err := cli.Status()
		if tc.ExpectErr {
			assert.Equal(t, tc.Expect, err.Error(), tc.Name)
		} else {
			assert.Nil(t, err, tc.Name)
			assert.Equal(t, tc.Expect, avail, tc.Name)
		}
	}
}
