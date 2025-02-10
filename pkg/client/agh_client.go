package client

import (
	"net/http"
	"strings"
)

const (
	EndpointStatsPath = "/control/stats"
)

type AghClient interface {
	FetchStats(response *StatsResponse) error
}

type aghClientImpl struct {
	client   *http.Client
	username string
	password string
	statsUrl string
}

func (c *aghClientImpl) FetchStats(response *StatsResponse) error {
	return fetchStats(response, c.client, c.statsUrl, c.username, c.password)
}

func NewClient(client *http.Client, baseProxyUrl string, username string, password string) AghClient {
	baseProxyUrl = strings.TrimSuffix(baseProxyUrl, "/")
	return &aghClientImpl{
		client:   client,
		statsUrl: baseProxyUrl + EndpointStatsPath,
		username: username,
		password: password,
	}
}
