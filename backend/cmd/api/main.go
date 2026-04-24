package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"sstpa-tool/backend/internal/config"
	apihttp "sstpa-tool/backend/internal/http"
	"sstpa-tool/backend/internal/version"
)

func main() {
	cfg := config.Load()
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           apihttp.NewRouter(version.Dev),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}

	slog.Info("starting sstpa api", "addr", server.Addr)

	err := server.ListenAndServe()
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return
	}

	slog.Error("sstpa api exited", "error", err)
	os.Exit(1)
}
