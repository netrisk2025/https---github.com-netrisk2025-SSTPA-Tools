// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"sstpa-tool/backend/internal/config"
	apihttp "sstpa-tool/backend/internal/http"
	"sstpa-tool/backend/internal/neo4jx"
	"sstpa-tool/backend/internal/schema"
	"sstpa-tool/backend/internal/version"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	driver, err := neo4jx.Open(ctx, neo4jx.Config{
		URI:      cfg.Neo4jURI,
		User:     cfg.Neo4jUser,
		Password: cfg.Neo4jPassword,
		Database: cfg.Neo4jDatabase,
		Timeout:  cfg.Neo4jTimeout,
	})
	if err != nil {
		slog.Error("neo4j connection failed", "error", err)
		os.Exit(1)
	}
	if driver != nil {
		defer driver.Close(ctx)
		if err := schema.Bootstrap(ctx, driver, cfg.Neo4jDatabase); err != nil {
			slog.Error("neo4j schema bootstrap failed", "error", err)
			os.Exit(1)
		}
		slog.Info("neo4j schema ready", "database", cfg.Neo4jDatabase)
	} else {
		slog.Info("neo4j disabled; set SSTPA_NEO4J_URI to enable graph persistence")
	}

	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           apihttp.NewRouterWithOptions(apihttp.RouterOptions{Version: version.Dev, Driver: driver, DatabaseName: cfg.Neo4jDatabase}),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}

	slog.Info("starting sstpa api", "addr", server.Addr)

	err = server.ListenAndServe()
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return
	}

	slog.Error("sstpa api exited", "error", err)
	os.Exit(1)
}
