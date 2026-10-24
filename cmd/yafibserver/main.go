package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jyork/yafib/v2/internal/engine"
	"github.com/jyork/yafib/v2/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	handler := server.NewServer(engine.NewFibonacciEngine(), logger)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/yafib", handler.YafibHandler)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed", "addr", srv.Addr, "error", err)
		os.Exit(1)
	}
}
