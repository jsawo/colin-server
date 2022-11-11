package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jsawo/colin-server/ws"
)

var (
	server *http.Server
	done   chan (interface{})
)

const (
	serverPort = "9111"
)

func StartServer() {
	fmt.Println("Starting HTTP server")
	done = make(chan interface{})

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
}
