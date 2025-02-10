package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/malbanese/adguardhomestats/internal/proxy"
	"github.com/malbanese/adguardhomestats/pkg/client"
)

const (
	DefaultPort = 3001
	DefaultUrl  = "http://127.0.0.1"
	DefaultHost = "127.0.0.1"
)

func main() {
	// Setup command line flags
	username := flag.String("u", "", "Username for authentication (required)")
	password := flag.String("p", "", "Password for authentication (required)")
	url := flag.String("url", DefaultUrl, "Base URL to request")
	host := flag.String("host", DefaultHost, "Host the proxy will bind with")
	port := flag.String("port", strconv.Itoa(DefaultPort), "Port the proxy will bind with")

	// Parse command line flags
	flag.Parse()

	if *url == "" || *username == "" || *password == "" {
		// Invalid command line args
		flag.Usage()
	} else {
		httpClient := http.Client{}
		aghClient := client.NewClient(&httpClient, *url, *username, *password)
		server := proxy.NewProxyServer(aghClient, *host+":"+*port)
		log.Fatal(server.Start())
	}
}
