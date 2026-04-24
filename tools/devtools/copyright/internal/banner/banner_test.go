// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package banner

import "testing"

func TestHasBannerDetectsExistingBanner(t *testing.T) {
	src := Prepend("package sample\n")
	if !HasBanner(src) {
		t.Fatal("expected HasBanner to be true after Prepend")
	}
}

func TestHasBannerDetectsMissingBanner(t *testing.T) {
	if HasBanner("package sample\n") {
		t.Fatal("expected HasBanner to be false on bare source")
	}
}

func TestPrependIsIdempotent(t *testing.T) {
	once := Prepend("package sample\n")
	twice := Prepend(once)
	if once != twice {
		t.Fatalf("expected Prepend to be idempotent\nonce=%q\ntwice=%q", once, twice)
	}
}

func TestPrependPreservesShebang(t *testing.T) {
	src := "#!/usr/bin/env node\nconsole.log('hi')\n"
	got := Prepend(src)
	if !stringsHasPrefix(got, "#!/usr/bin/env node") {
		t.Fatal("shebang must remain first line")
	}
	if !HasBanner(got) {
		t.Fatal("banner must be present after shebang")
	}
}

func stringsHasPrefix(value string, prefix string) bool {
	if len(value) < len(prefix) {
		return false
	}

	return value[:len(prefix)] == prefix
}
