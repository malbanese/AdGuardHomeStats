package client

import (
	"net/http"
	"strings"
)

const (
	EndpointStatsPath = "/control/stats"
)

type AghClient interface {
	FetchStats(response *AghStatsResponse) error
}

type aghClientImpl struct {
	client   *http.Client
	username string
	password string
	statsUrl string
}

func (c *aghClientImpl) FetchStats(response *AghStatsResponse) error {
	return FetchStats(response, c.client, c.statsUrl, c.username, c.password)
}

func NewClient(client *http.Client, baseProxyUrl string, username string, password string) AghClient {
	return &aghClientImpl{
		client:   client,
		statsUrl: formatUrlForPath(baseProxyUrl, EndpointStatsPath),
		username: username,
		password: password,
	}
}

func formatUrlForPath(url string, path string) string {
	if !strings.HasSuffix(url, path) {
		url = strings.TrimSuffix(url, "/")
		url += path
	}

	return url
}
