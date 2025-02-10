package proxy

import (
	"net/http"

	"github.com/malbanese/adguardhomestats/internal/handler"
	"github.com/malbanese/adguardhomestats/pkg/client"
)

func SetupHttpRoutes(baseProxyUrl string, username string, password string) {
	statsEndpoint := client.AdGuardHomeEndpointInfo{
		Username: username,
		Password: password,
		Url:      handler.FormatStatsUrl(baseProxyUrl),
	}

	http.HandleFunc(handler.StatsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		handler.StatsHttpHandler(&statsEndpoint, w, r)
	})
}
