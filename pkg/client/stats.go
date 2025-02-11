package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type StatsResponse struct {
	TimeUnits           string  `json:"time_units"`
	NumDNSQueries       int     `json:"num_dns_queries"`
	NumBlockedFiltering int     `json:"num_blocked_filtering"`
	AvgProcessingTime   float64 `json:"avg_processing_time"`
}

type StatsClient struct {
	Client *http.Client
	Auth   AuthType
	URL    string
}

func (c *StatsClient) FetchStats(response *StatsResponse) error {
	// Setup initial request
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Add authentication to our request
	c.Auth.Authenticate(req)

	// Execute the request
	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	// Check the response was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Error: %s", resp.Status)
	}

	// Read the response data
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}

	// Decode the JSON, to strip undesired fields
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return fmt.Errorf("error parsing body json: %w", err)
	}

	return nil
}
