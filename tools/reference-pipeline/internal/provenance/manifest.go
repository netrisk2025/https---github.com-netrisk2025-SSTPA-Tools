package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const SchemaVersion = "1.0.0"

type Artifact struct {
	Path        string `yaml:"path" json:"path"`
	Role        string `yaml:"role" json:"role"`
	SHA256      string `yaml:"sha256" json:"sha256"`
	Bytes       int64  `yaml:"bytes" json:"bytes"`
	RecordCount int    `yaml:"recordCount,omitempty" json:"recordCount,omitempty"`
}

type Manifest struct {
	SchemaVersion      string     `yaml:"schemaVersion" json:"schemaVersion"`
	Framework          string     `yaml:"framework" json:"framework"`
	Version            string     `yaml:"version" json:"version"`
	Stage              string     `yaml:"stage" json:"stage"`
	GeneratedAt        string     `yaml:"generatedAt" json:"generatedAt"`
	SourceTitle        string     `yaml:"sourceTitle,omitempty" json:"sourceTitle,omitempty"`
	SourceLastModified string     `yaml:"sourceLastModified,omitempty" json:"sourceLastModified,omitempty"`
	RawArtifacts       []Artifact `yaml:"rawArtifacts" json:"rawArtifacts"`
	StagedArtifacts    []Artifact `yaml:"stagedArtifacts" json:"stagedArtifacts"`
}

func (manifest Manifest) Validate() error {
	if manifest.SchemaVersion == "" {
		return fmt.Errorf("schema version is required")
	}

	if manifest.Framework == "" {
		return fmt.Errorf("framework is required")
	}

	if manifest.Version == "" {
		return fmt.Errorf("version is required")
	}

	if manifest.Stage == "" {
		return fmt.Errorf("stage is required")
	}

	if manifest.GeneratedAt == "" {
		return fmt.Errorf("generated at is required")
	}

	if len(manifest.RawArtifacts) == 0 {
		return fmt.Errorf("at least one raw artifact is required")
	}

	if len(manifest.StagedArtifacts) == 0 {
		return fmt.Errorf("at least one staged artifact is required")
	}

	seen := make(map[string]struct{})
	for _, artifact := range append(append([]Artifact{}, manifest.RawArtifacts...), manifest.StagedArtifacts...) {
		if err := artifact.Validate(); err != nil {
			return fmt.Errorf("artifact %q: %w", artifact.Path, err)
		}

		if _, ok := seen[artifact.Path]; ok {
			return fmt.Errorf("duplicate artifact path %q", artifact.Path)
		}
		seen[artifact.Path] = struct{}{}
	}

	return nil
}

func (artifact Artifact) Validate() error {
	if artifact.Path == "" {
		return fmt.Errorf("path is required")
	}

	if artifact.Role == "" {
		return fmt.Errorf("role is required")
	}

	if len(artifact.SHA256) != 64 {
		return fmt.Errorf("sha256 must be 64 hex characters")
	}

	for _, ch := range artifact.SHA256 {
		if !strings.ContainsRune("0123456789abcdef", ch) {
			return fmt.Errorf("sha256 must be lowercase hex")
		}
	}

	if artifact.Bytes <= 0 {
		return fmt.Errorf("bytes must be greater than zero")
	}

	if artifact.RecordCount < 0 {
		return fmt.Errorf("record count must not be negative")
	}

	return nil
}

func ReadFile(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}

	return manifest, nil
}

func WriteFile(path string, manifest Manifest) error {
	if err := manifest.Validate(); err != nil {
		return err
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func BuildArtifact(path, role string, recordCount int) (Artifact, error) {
	sum, bytes, err := hashFile(path)
	if err != nil {
		return Artifact{}, err
	}

	return Artifact{
		Path:        filepath.ToSlash(filepath.Clean(path)),
		Role:        role,
		SHA256:      sum,
		Bytes:       bytes,
		RecordCount: recordCount,
	}, nil
}

func hashFile(path string) (string, int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", 0, err
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), int64(len(data)), nil
}
