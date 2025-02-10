package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AghStatsResponse struct {
	TimeUnits           string  `json:"time_units"`
	NumDnsQueries       int     `json:"num_dns_queries"`
	NumBlockedFiltering int     `json:"num_blocked_filtering"`
	AvgProcessingTime   float64 `json:"avg_processing_time"`
}

func FetchStats(response *AghStatsResponse, client *http.Client, url string, username string, password string) error {
	// Setup initial request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// AdGuard Home uses basic authentication
	req.SetBasicAuth(username, password)

	// Execute the request
	resp, err := client.Do(req)
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
