package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/malbanese/adguardhomestats/pkg/client"
)

// HTTP routes the proxy server will bind to
type Routes struct {
	Stats string
}

// Base server definition
type Server struct {
	Stats StatsClient
}

// Client to be used for fetching AdGuard Home statistics
type StatsClient interface {
	FetchStats(response *client.StatsResponse) error
}

// Returns an HTTP server which will operate on the given host and port.
// All relevant proxy routes have been mounted.
func (p *Server) NewHTTPServer(routes Routes, host string, port uint) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(routes.Stats, p.onStatsRequest)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}

	return server
}

func (p *Server) onStatsRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Accepting stats request [%s] from [%s]", r.URL, r.RemoteAddr)
	var response client.StatsResponse
	err := p.Stats.FetchStats(&response)
	writeProxyResponse(w, response, err)
}

func writeProxyResponse(
	w http.ResponseWriter,
	response any,
	err error,
) {
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Println("Error while writing the stats response: %w", err)
		return
	}
}
