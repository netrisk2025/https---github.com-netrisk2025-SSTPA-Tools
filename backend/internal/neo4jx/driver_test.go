// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package neo4jx

import "testing"

func TestConfigEnabled(t *testing.T) {
	if (Config{}).Enabled() {
		t.Fatal("empty config should be disabled")
	}

	if !(Config{URI: "bolt://localhost:7687"}).Enabled() {
		t.Fatal("URI should enable Neo4j")
	}
}
