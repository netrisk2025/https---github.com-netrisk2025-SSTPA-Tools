// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package testhelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Neo4jFixture struct {
	URI       string
	User      string
	Password  string
	Driver    neo4j.DriverWithContext
	terminate func(context.Context) error
}

func StartNeo4j(t testing.TB) Neo4jFixture {
	t.Helper()

	ctx := context.Background()
	const password = "test-password"

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "neo4j:5.26-community",
			ExposedPorts: []string{"7687/tcp"},
			Env: map[string]string{
				"NEO4J_AUTH": "neo4j/" + password,
			},
			WaitingFor: wait.ForListeningPort("7687/tcp").WithStartupTimeout(2 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		t.Skip("docker unavailable: " + err.Error())
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("get neo4j container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "7687/tcp")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("get neo4j mapped port: %v", err)
	}

	uri := fmt.Sprintf("bolt://%s:%s", host, port.Port())
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth("neo4j", password, ""))
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("create neo4j driver: %v", err)
	}

	verifyCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := driver.VerifyConnectivity(verifyCtx); err != nil {
		_ = driver.Close(ctx)
		_ = container.Terminate(ctx)
		t.Skip("docker unavailable: " + err.Error())
	}

	t.Cleanup(func() {
		_ = driver.Close(context.Background())
		_ = container.Terminate(context.Background())
	})

	return Neo4jFixture{
		URI:      uri,
		User:     "neo4j",
		Password: password,
		Driver:   driver,
		terminate: func(ctx context.Context) error {
			return container.Terminate(ctx)
		},
	}
}
