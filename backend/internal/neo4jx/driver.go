// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package neo4jx

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Config struct {
	URI      string
	User     string
	Password string
	Database string
	Timeout  time.Duration
}

func (cfg Config) Enabled() bool {
	return cfg.URI != ""
}

func Open(ctx context.Context, cfg Config) (neo4j.DriverWithContext, error) {
	if !cfg.Enabled() {
		return nil, nil
	}

	if cfg.User == "" {
		return nil, fmt.Errorf("neo4j user is required when URI is configured")
	}

	driver, err := neo4j.NewDriverWithContext(cfg.URI, neo4j.BasicAuth(cfg.User, cfg.Password, ""))
	if err != nil {
		return nil, err
	}

	pingCtx := ctx
	cancel := func() {}
	if cfg.Timeout > 0 {
		pingCtx, cancel = context.WithTimeout(ctx, cfg.Timeout)
	}
	defer cancel()

	if err := driver.VerifyConnectivity(pingCtx); err != nil {
		_ = driver.Close(ctx)
		return nil, err
	}

	return driver, nil
}
