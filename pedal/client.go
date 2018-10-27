package pedal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// BaseURL provides the base for performing queries
var BaseURL = "https://oslobysykkel.no/api/v1/"

type httpClient struct {
	baseURL          string
	clientIdentifier string
	client           *http.Client
}

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
func (c *httpClient) Status() (*Status, error) {
	var status Status

	err := c.do("status", &status)
	if err != nil {
		return nil, err
	}

	return &status, err
}

// Stations loads all known stations from the API
func (c *httpClient) Stations() (*Stations, error) {
	var stations Stations

	err := c.do("stations", &stations)
	if err != nil {
		return nil, err
	}

	return &stations, err
}

// Availability fetches the availability of bikes and locks at all
// locations.
func (c *httpClient) Availability() (*StationAvailability, error) {
	var stationAvailability StationAvailability

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
		return fmt.Errorf("expected success, got: %d", resp.StatusCode)
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
