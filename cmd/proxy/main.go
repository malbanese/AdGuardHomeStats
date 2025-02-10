package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/malbanese/adguardhomestats/internal/proxy"
)

const (
	DefaultPort = 3001
)

func main() {
	// Setup command line flags
	url := flag.String("url", "", "Base URL to request (required)")
	username := flag.String("u", "", "Username for authentication (required)")
	password := flag.String("p", "", "Password for authentication (required)")
	port := flag.String("port", strconv.Itoa(DefaultPort), "Port the proxy will be started on")

	// Parse command line flags
	flag.Parse()

	// Validate command line flags
	if *url == "" || *username == "" || *password == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Setup the routes and start listening
	proxy.SetupHttpRoutes(*url, *username, *password)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
