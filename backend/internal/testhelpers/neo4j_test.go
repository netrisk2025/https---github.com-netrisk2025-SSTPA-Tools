// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package testhelpers

import (
	"context"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestStartNeo4jSmoke(t *testing.T) {
	fixture := StartNeo4j(t)

	ctx := context.Background()
	session := fixture.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, "RETURN 1 AS n", nil)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return record.Values[0], nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if value != int64(1) {
		t.Fatalf("RETURN 1 result = %#v, want int64(1)", value)
	}
}
