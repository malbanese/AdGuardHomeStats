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
)

func main() {
	// Setup command line flags
	username := flag.String("u", "", "Username for authentication (required)")
	password := flag.String("p", "", "Password for authentication (required)")
	url := flag.String("url", DefaultUrl, "Base URL to request")
	port := flag.String("port", strconv.Itoa(DefaultPort), "Port the proxy will be started on")

	// Parse command line flags
	flag.Parse()

	if *url == "" || *username == "" || *password == "" {
		// Invalid command line args
		flag.Usage()
	} else {
		// Setup the routes and start listening
		aghClient := client.NewClient(&http.Client{}, *url, *username, *password)
		proxy.SetupHttpRoutes(aghClient)
		log.Fatal(http.ListenAndServe(":"+*port, nil))
	}
}
