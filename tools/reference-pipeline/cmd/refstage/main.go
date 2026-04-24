// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package main

import (
	"flag"
	"fmt"
	"os"

	"sstpa-tool/reference-pipeline/internal/nist"
	"sstpa-tool/reference-pipeline/internal/provenance"
)

func main() {
	if hasArg("--print-layout") {
		printLayout()
		return
	}

	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "manifest":
		runManifest(os.Args[2:])
	case "stage":
		runStage(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func runManifest(args []string) {
	if len(args) < 2 || args[0] != "validate" {
		usage()
		os.Exit(2)
	}

	manifest, err := provenance.ReadFile(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read manifest: %v\n", err)
		os.Exit(1)
	}

	if err := manifest.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid manifest: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(
		"manifest valid: framework=%s version=%s raw=%d staged=%d\n",
		manifest.Framework,
		manifest.Version,
		len(manifest.RawArtifacts),
		len(manifest.StagedArtifacts),
	)
}

func runStage(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}

	switch args[0] {
	case "nist":
		runStageNIST(args[1:])
	default:
		usage()
		os.Exit(2)
	}
}

func runStageNIST(args []string) {
	fs := flag.NewFlagSet("stage nist", flag.ExitOnError)
	catalogPath := fs.String("catalog", "", "path to the raw NIST OSCAL catalog JSON")
	licensePath := fs.String("license", "", "path to the raw NIST license file")
	outDir := fs.String("out-dir", "", "directory for staged output files")
	manifestPath := fs.String("manifest", "", "path to write the provenance manifest")
	fs.Parse(args)

	if *catalogPath == "" || *licensePath == "" || *outDir == "" || *manifestPath == "" {
		fs.Usage()
		os.Exit(2)
	}

	result, err := nist.Stage(nist.StageOptions{
		CatalogPath:  *catalogPath,
		LicensePath:  *licensePath,
		OutDir:       *outDir,
		ManifestPath: *manifestPath,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "stage nist: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(
		"staged nist-sp800-53 %s: collections=%d items=%d edge-candidates=%d citations=%d\nmanifest: %s\n",
		result.Version,
		result.CollectionCount,
		result.ItemCount,
		result.EdgeCandidateCount,
		result.CitationCount,
		result.ManifestPath,
	)
}

func hasArg(want string) bool {
	for _, arg := range os.Args[1:] {
		if arg == want {
			return true
		}
	}

	return false
}

func printLayout() {
	fmt.Println("reference-data/raw")
	fmt.Println("reference-data/staged")
	fmt.Println("reference-data/normalized")
	fmt.Println("reference-data/manifests")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  refstage --print-layout\n")
	fmt.Fprintf(os.Stderr, "  refstage manifest validate <path>\n")
	fmt.Fprintf(os.Stderr, "  refstage stage nist --catalog <path> --license <path> --out-dir <dir> --manifest <path>\n")
}
