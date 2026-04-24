// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package nist

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"sstpa-tool/reference-pipeline/internal/provenance"
)

func TestStage(t *testing.T) {
	tempDir := t.TempDir()
	outDir := filepath.Join(tempDir, "staged", "nist-sp800-53", "5.2.0")
	manifestPath := filepath.Join(tempDir, "manifests", "nist-sp800-53", "5.2.0.yaml")

	result, err := Stage(StageOptions{
		CatalogPath:  filepath.Join("testdata", "sample_catalog.json"),
		LicensePath:  filepath.Join("testdata", "LICENSE.md"),
		OutDir:       outDir,
		ManifestPath: manifestPath,
	})
	if err != nil {
		t.Fatalf("stage: %v", err)
	}

	if result.Version != "5.2.0" {
		t.Fatalf("unexpected version %q", result.Version)
	}

	if result.CollectionCount != 1 {
		t.Fatalf("unexpected collection count %d", result.CollectionCount)
	}

	if result.ItemCount != 4 {
		t.Fatalf("unexpected item count %d", result.ItemCount)
	}

	if result.EdgeCandidateCount != 7 {
		t.Fatalf("unexpected edge candidate count %d", result.EdgeCandidateCount)
	}

	if result.CitationCount != 1 {
		t.Fatalf("unexpected citation count %d", result.CitationCount)
	}

	var items []Item
	readNDJSON(t, filepath.Join(outDir, "items.ndjson"), &items)
	if items[0].ItemType != "family" {
		t.Fatalf("expected first item to be family, got %q", items[0].ItemType)
	}

	var edges []EdgeCandidate
	readNDJSON(t, filepath.Join(outDir, "edge-candidates.ndjson"), &edges)

	foundFragment := false
	foundRequired := false
	for _, edge := range edges {
		if edge.RelationshipType == "incorporated-into" && edge.ToKind == "fragment" {
			foundFragment = true
		}
		if edge.RelationshipType == "required" && edge.ToItemID == "AC-2" {
			foundRequired = true
		}
	}

	if !foundFragment {
		t.Fatal("expected fragment edge candidate")
	}

	if !foundRequired {
		t.Fatal("expected required edge candidate to resolve to AC-2")
	}

	manifest, err := provenance.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}

	if err := manifest.Validate(); err != nil {
		t.Fatalf("validate manifest: %v", err)
	}
}

func readNDJSON[T any](t *testing.T, path string, out *[]T) {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record T
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			t.Fatalf("decode ndjson %s: %v", path, err)
		}
		*out = append(*out, record)
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("scan ndjson %s: %v", path, err)
	}
}
