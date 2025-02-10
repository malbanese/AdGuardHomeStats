package proxy

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/malbanese/adguardhomestats/pkg/client"
)

func SetupHttpRoutes(aghClient client.AghClient) {
	http.HandleFunc(client.EndpointStatsPath, func(w http.ResponseWriter, r *http.Request) {
		statsHttpHandler(aghClient, w, r)
	})
}

func statsHttpHandler(
	aghClient client.AghClient,
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf("Accepting request from %s", r.RemoteAddr)

	var response client.AghStatsResponse
	err := aghClient.FetchStats(&response)

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
