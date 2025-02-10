package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AdGuardHomeStatsResponse struct {
	TimeUnits           string  `json:"time_units"`
	NumDnsQueries       int     `json:"num_dns_queries"`
	NumBlockedFiltering int     `json:"num_blocked_filtering"`
	AvgProcessingTime   float64 `json:"avg_processing_time"`
}

type AdGuardHomeEndpointInfo struct {
	Username string
	Password string
	Url      string
}

func FetchStats(response *AdGuardHomeStatsResponse, endpoint *AdGuardHomeEndpointInfo) error {
	// Setup initial request
	req, err := http.NewRequest("GET", endpoint.Url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// AdGuard Home uses basic plain-text authentication
	req.SetBasicAuth(endpoint.Username, endpoint.Password)

	// AdGuard Home uses basic http
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	// Check the response was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Error: %d, %s", resp.StatusCode, resp.Status)
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
