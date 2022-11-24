package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jsawo/colin-server/ws"
)

var (
	server *http.Server
	done   = make(chan interface{})
)

const (
	serverPort = "9111"
)

func StartServer() {
	fmt.Println("Starting HTTP server")

	server = makeHTTPServer(serverPort)

	go ws.WriteToWS()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ERR: Failed to start http server: %v\n", err)
		}
	}()

	fmt.Printf("Listening for HTTP requests at port %v\n", serverPort)

	<-done
}

func makeHTTPServer(serverPort string) *http.Server {
	mux := &http.ServeMux{}

	registerRoutes(mux)

	srv := &http.Server{
		ReadTimeout:       120 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 120 * time.Second,
		Handler:           mux,
		Addr:              ":" + serverPort,
	}

	return srv
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws", ws.ServeWS)
	mux.HandleFunc("/toc", serveTOC)
}

func serveTOC(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]map[string]string)
	for _, collector := range CollectorConfigs {
		if collector.Enabled {
			resp[collector.Channel] = map[string]string{
				"title":       collector.Title,
				"description": collector.Description,
				"type":        collector.Type.ToString(),
				"frequency":   collector.Frequency.String(),
			}
		}
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
