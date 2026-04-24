// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"sstpa-tool/devtools/copyright/internal/banner"
)

var roots = []string{"backend", "apps", "addons", "packages", "tools"}

var skipDirs = map[string]struct{}{
	".git":         {},
	".vite":        {},
	"coverage":     {},
	"dist":         {},
	"node_modules": {},
	"target":       {},
}

var sourceExts = map[string]struct{}{
	".cjs": {},
	".go":  {},
	".js":  {},
	".mjs": {},
	".rs":  {},
	".ts":  {},
	".tsx": {},
}

func main() {
	check := flag.Bool("check", false, "exit nonzero if any source file lacks the banner")
	apply := flag.Bool("apply", false, "prepend the banner to source files that lack it")
	flag.Parse()

	if *check == *apply {
		fmt.Fprintln(os.Stderr, "exactly one of --check or --apply is required")
		os.Exit(2)
	}

	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "find repo root: %v\n", err)
		os.Exit(1)
	}

	var missing []string
	for _, root := range roots {
		rootPath := filepath.Join(repoRoot, root)
		err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.IsDir() {
				if _, skip := skipDirs[d.Name()]; skip {
					return filepath.SkipDir
				}
				return nil
			}

			if _, ok := sourceExts[filepath.Ext(path)]; !ok {
				return nil
			}

			relPath, err := filepath.Rel(repoRoot, path)
			if err != nil {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			if banner.HasBanner(string(content)) {
				return nil
			}

			if *check {
				missing = append(missing, relPath)
				return nil
			}

			return os.WriteFile(path, []byte(banner.Prepend(string(content))), 0o644)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "walk %s: %v\n", root, err)
			os.Exit(1)
		}
	}

	if len(missing) == 0 {
		return
	}

	fmt.Fprintln(os.Stderr, "files missing copyright banner:")
	for _, path := range missing {
		fmt.Fprintf(os.Stderr, "  %s\n", path)
	}
	os.Exit(1)
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if hasFile(filepath.Join(dir, "go.work")) && hasFile(filepath.Join(dir, "Makefile")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.work and Makefile not found from %s", dir)
		}
		dir = parent
	}
}

func hasFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
