// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package identity

import (
	"regexp"
	"testing"
)

func TestNewUUIDFormatAndUniqueness(t *testing.T) {
	pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	first := NewUUID()
	second := NewUUID()

	if first == second {
		t.Fatal("expected two generated UUIDs to differ")
	}

	if !pattern.MatchString(first) {
		t.Fatalf("first UUID has unexpected format: %q", first)
	}
}
