package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jyork/yafib/v2/internal/engine"
	"github.com/jyork/yafib/v2/internal/server"
)

func main() {
	handler := server.NewServer(engine.NewFibonacciEngine())
	mux := http.NewServeMux()
	mux.HandleFunc("/api/yafib", handler.YafibHandler)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
