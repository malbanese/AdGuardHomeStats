package proxy

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	statClient "github.com/malbanese/adguardhomestats/pkg/client"
)

const (
	statsEndpointPath = "/control/stats"
)

func SetupHttpRoutes(client *http.Client, baseProxyUrl string, username string, password string) {
	statsEndpoint := statClient.AghStatsEndpointInfo{
		Username: username,
		Password: password,
		Url:      formatStatsUrl(baseProxyUrl),
	}

	http.HandleFunc(statsEndpointPath, func(w http.ResponseWriter, r *http.Request) {
		statsHttpHandler(client, &statsEndpoint, w, r)
	})
}

func statsHttpHandler(
	client *http.Client,
	endpoint *statClient.AghStatsEndpointInfo,
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf("Accepting request from %s", r.RemoteAddr)

	var response statClient.AghStatsResponse

	err := statClient.FetchStats(client, &response, endpoint)
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

func formatStatsUrl(url string) string {
	if !strings.HasSuffix(url, statsEndpointPath) {
		url = strings.TrimSuffix(url, "/")
		url += statsEndpointPath
	}

	return url
}
