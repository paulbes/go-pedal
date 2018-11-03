package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/paulbes/go-pedal/pedal/model"
)

// Client defines the interface that
// an API client must implement
type Client interface {
	Stations() (*model.Stations, error)
	Availability() (*model.StationAvailability, error)
	Status() (*model.Status, error)
}

// BaseURL provides the base for performing queries
var BaseURL = "https://oslobysykkel.no/api/v1/"

type httpClient struct {
	baseURL          string
	clientIdentifier string
	client           *http.Client
}

// NewHTTPClient creates an http client that can communicate with the
// oslo city bike API
func NewHTTPClient(clientID string, timeoutInSec int) (Client, error) {
	if len(clientID) == 0 {
		return nil, fmt.Errorf("client identifier is required")
	}
	return &httpClient{
		baseURL:          BaseURL,
		clientIdentifier: clientID,
		client: &http.Client{
			Timeout: time.Duration(timeoutInSec) * time.Second,
		},
	}, nil
}

// Status loads the status of the stations
func (c *httpClient) Status() (*model.Status, error) {
	var status = struct {
		Status model.Status `json:"status"`
	}{}

	err := c.do("status", &status)
	if err != nil {
		return nil, err
	}

	return &status.Status, err
}

// Stations loads all known stations from the API
func (c *httpClient) Stations() (*model.Stations, error) {
	var stations model.Stations

	err := c.do("stations", &stations)
	if err != nil {
		return nil, err
	}

	return &stations, err
}

// Availability fetches the availability of bikes and locks at all
// locations.
func (c *httpClient) Availability() (*model.StationAvailability, error) {
	var stationAvailability model.StationAvailability

	err := c.do("stations/availability", &stationAvailability)
	if err != nil {
		return nil, err
	}

	return &stationAvailability, nil
}

// do executes a request towards the oslo city bike API
func (c *httpClient) do(endpoint string, to interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.baseURL, endpoint), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Client-Identifier", c.clientIdentifier)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		var reason string
		if err != nil {
			reason = "unknown"
		} else {
			reason = string(data)
		}

		return fmt.Errorf("failed to invoke API, got error code: %d, reason: %s", resp.StatusCode, reason)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, to)
	if err != nil {
		return err
	}

	return nil
}
