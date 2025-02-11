package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/malbanese/adguardhomestats/internal/proxy"
	"github.com/malbanese/adguardhomestats/pkg/client"
)

// Default values to be used in arguments
const (
	DefaultPort uint   = 3001
	DefaultUrl  string = "http://127.0.0.1"
	DefaultHost string = "127.0.0.1"
)

// Proxy routes
var routes = proxy.ProxyRoutes{
	Stats: "/control/stats",
}

// Argument convenience bundling
type ProxyArgs struct {
	Username *string
	Password *string
	Url      *string
	Host     *string
	Port     *uint
}

func (a *ProxyArgs) IsValid() bool {
	return *a.Username != "" && *a.Password != "" && *a.Url != "" && *a.Host != ""
}

// Main entry point
func main() {
	args := ProxyArgs{
		Username: flag.String("u", "", "Username for authentication (required)"),
		Password: flag.String("p", "", "Password for authentication (required)"),
		Url:      flag.String("url", DefaultUrl, "Base URL to request"),
		Host:     flag.String("host", DefaultHost, "Host the proxy will bind with"),
		Port:     flag.Uint("port", DefaultPort, "Port the proxy will bind with"),
	}

	flag.Parse()

	if !args.IsValid() {
		flag.Usage()
	} else {
		startServer(&args)
	}
}

func startServer(args *ProxyArgs) {
	httpClient := http.Client{}

	authType := client.BasicAuthType{
		Username: *args.Username,
		Password: *args.Password,
	}

	url := strings.TrimSuffix(*args.Url, "/")

	server := proxy.ProxyServer{
		Stats: &client.StatsClient{
			Url:    url + routes.Stats,
			Auth:   &authType,
			Client: &httpClient,
		},
	}

	httpServer := server.NewHttpServer(routes, *args.Host, *args.Port)
	log.Fatal(httpServer.ListenAndServe())
}
