package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/malbanese/adguardhomestats/pkg/client"
)

const (
	StatsEndpoint = "/control/stats"
)

func StatsHttpHandler(
	endpoint *client.AdGuardHomeEndpointInfo,
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf("Accepting request from %s", r.RemoteAddr)

	var response client.AdGuardHomeStatsResponse

	err := client.FetchStats(&response, endpoint)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Println("Error while writing the stats response: %w", err)
	}
}

func FormatStatsUrl(url string) string {
	if !strings.HasSuffix(url, StatsEndpoint) {
		url = strings.TrimSuffix(url, "/")
		url += StatsEndpoint
	}

	return url
}
