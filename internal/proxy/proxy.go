package proxy

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/malbanese/adguardhomestats/pkg/client"
)

type ProxyServer struct {
	server *http.Server
}

func (p *ProxyServer) Start() error {
	return p.server.ListenAndServe()
}

func NewProxyServer(
	aghClient client.AghClient,
	addr string,
) *ProxyServer {
	mux := http.NewServeMux()
	bindRoutes(mux, aghClient)

	return &ProxyServer{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func bindRoutes(mux *http.ServeMux, aghClient client.AghClient) {
	mux.HandleFunc(client.EndpointStatsPath, func(w http.ResponseWriter, r *http.Request) {
		var response client.StatsResponse
		proxy(w, r, func() (any, error) {
			return &response, aghClient.FetchStats(&response)
		})
	})
}

func proxy(
	w http.ResponseWriter,
	r *http.Request,
	method func() (any, error),
) {
	log.Printf("Accepting request [%s] from [%s]", r.URL, r.RemoteAddr)

	response, err := method()
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
