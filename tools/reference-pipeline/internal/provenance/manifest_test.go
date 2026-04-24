// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package provenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManifestValidate(t *testing.T) {
	manifest := Manifest{
		SchemaVersion: SchemaVersion,
		Framework:     "nist-sp800-53",
		Version:       "5.2.0",
		Stage:         "staged",
		GeneratedAt:   "2026-04-22T12:00:00Z",
		RawArtifacts: []Artifact{
			{
				Path:   "reference-data/raw/nist-sp800-53/v5.2.0/NIST_SP-800-53_rev5_catalog.json",
				Role:   "catalog",
				SHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Bytes:  123,
			},
		},
		StagedArtifacts: []Artifact{
			{
				Path:        "reference-data/staged/nist-sp800-53/v5.2.0/items.ndjson",
				Role:        "items",
				SHA256:      "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				Bytes:       456,
				RecordCount: 3,
			},
		},
	}

	if err := manifest.Validate(); err != nil {
		t.Fatalf("expected manifest to validate: %v", err)
	}
}

func TestManifestValidateRejectsDuplicateArtifactPaths(t *testing.T) {
	manifest := Manifest{
		SchemaVersion: SchemaVersion,
		Framework:     "nist-sp800-53",
		Version:       "5.2.0",
		Stage:         "staged",
		GeneratedAt:   "2026-04-22T12:00:00Z",
		RawArtifacts: []Artifact{
			{
				Path:   "same.yaml",
				Role:   "catalog",
				SHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Bytes:  123,
			},
		},
		StagedArtifacts: []Artifact{
			{
				Path:   "same.yaml",
				Role:   "items",
				SHA256: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				Bytes:  456,
			},
		},
	}

	if err := manifest.Validate(); err == nil {
		t.Fatal("expected duplicate path validation error")
	}
}

func TestBuildArtifact(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write sample file: %v", err)
	}

	artifact, err := BuildArtifact(path, "sample", 0)
	if err != nil {
		t.Fatalf("build artifact: %v", err)
	}

	if artifact.Role != "sample" {
		t.Fatalf("unexpected role %q", artifact.Role)
	}

	if artifact.Bytes != 5 {
		t.Fatalf("unexpected byte count %d", artifact.Bytes)
	}

	if artifact.SHA256 != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Fatalf("unexpected sha256 %q", artifact.SHA256)
	}
}
