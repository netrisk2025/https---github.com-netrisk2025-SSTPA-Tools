package nist

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sstpa-tool/reference-pipeline/internal/provenance"
)

const frameworkName = "nist-sp800-53"

type StageOptions struct {
	CatalogPath  string
	LicensePath  string
	OutDir       string
	ManifestPath string
}

type StageResult struct {
	Version            string
	ManifestPath       string
	CollectionCount    int
	ItemCount          int
	EdgeCandidateCount int
	CitationCount      int
}

type catalogDocument struct {
	Catalog catalog `json:"catalog"`
}

type catalog struct {
	Metadata metadata `json:"metadata"`
	Groups   []group  `json:"groups"`
}

type metadata struct {
	Title        string `json:"title"`
	Version      string `json:"version"`
	LastModified string `json:"last-modified"`
}

type group struct {
	ID       string    `json:"id"`
	Class    string    `json:"class"`
	Title    string    `json:"title"`
	Controls []control `json:"controls"`
}

type control struct {
	ID       string    `json:"id"`
	Class    string    `json:"class"`
	Title    string    `json:"title"`
	Props    []prop    `json:"props"`
	Links    []link    `json:"links"`
	Controls []control `json:"controls"`
}

type prop struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Value string `json:"value"`
}

type link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Text string `json:"text"`
}

type Collection struct {
	CollectionID   string `json:"collectionId"`
	CollectionType string `json:"collectionType"`
	Name           string `json:"name"`
	SourcePath     string `json:"sourcePath"`
}

type Item struct {
	Framework  string     `json:"framework"`
	Version    string     `json:"version"`
	ItemID     string     `json:"itemId"`
	ItemType   string     `json:"itemType"`
	Class      string     `json:"class,omitempty"`
	Title      string     `json:"title"`
	ParentID   string     `json:"parentId,omitempty"`
	FamilyID   string     `json:"familyId,omitempty"`
	SourcePath string     `json:"sourcePath"`
	Props      []Property `json:"props,omitempty"`
}

type Property struct {
	Name  string `json:"name"`
	Class string `json:"class,omitempty"`
	Value string `json:"value"`
}

type EdgeCandidate struct {
	FromID           string `json:"fromId"`
	RelationshipType string `json:"relationshipType"`
	ToRef            string `json:"toRef"`
	ToKind           string `json:"toKind"`
	ToItemID         string `json:"toItemId,omitempty"`
	SourcePath       string `json:"sourcePath"`
}

type Citation struct {
	FromID     string `json:"fromId"`
	Relation   string `json:"relation"`
	Href       string `json:"href"`
	Text       string `json:"text,omitempty"`
	SourcePath string `json:"sourcePath"`
}

type Metadata struct {
	Framework          string         `json:"framework"`
	Version            string         `json:"version"`
	SourceTitle        string         `json:"sourceTitle"`
	SourceLastModified string         `json:"sourceLastModified"`
	GeneratedAt        string         `json:"generatedAt"`
	Counts             map[string]int `json:"counts"`
}

type pendingEdge struct {
	FromID           string
	RelationshipType string
	ToRef            string
	SourcePath       string
}

func Stage(opts StageOptions) (StageResult, error) {
	doc, err := readCatalog(opts.CatalogPath)
	if err != nil {
		return StageResult{}, err
	}

	if doc.Catalog.Metadata.Version == "" {
		return StageResult{}, fmt.Errorf("catalog metadata.version is required")
	}

	if err := os.MkdirAll(opts.OutDir, 0o755); err != nil {
		return StageResult{}, err
	}

	items, itemIDs, pendingEdges, citations := extract(doc)
	collections := []Collection{
		{
			CollectionID:   "catalog",
			CollectionType: "catalog",
			Name:           doc.Catalog.Metadata.Title,
			SourcePath:     "catalog",
		},
	}
	edges := resolveEdges(itemIDs, pendingEdges)

	metadata := Metadata{
		Framework:          frameworkName,
		Version:            doc.Catalog.Metadata.Version,
		SourceTitle:        doc.Catalog.Metadata.Title,
		SourceLastModified: doc.Catalog.Metadata.LastModified,
		GeneratedAt:        time.Now().UTC().Format(time.RFC3339),
		Counts: map[string]int{
			"collections":    len(collections),
			"items":          len(items),
			"edgeCandidates": len(edges),
			"citations":      len(citations),
		},
	}

	metadataPath := filepath.Join(opts.OutDir, "metadata.json")
	collectionsPath := filepath.Join(opts.OutDir, "collections.ndjson")
	itemsPath := filepath.Join(opts.OutDir, "items.ndjson")
	edgesPath := filepath.Join(opts.OutDir, "edge-candidates.ndjson")
	citationsPath := filepath.Join(opts.OutDir, "citations.ndjson")

	if err := writeJSON(metadataPath, metadata); err != nil {
		return StageResult{}, err
	}

	if err := writeNDJSON(collectionsPath, collections); err != nil {
		return StageResult{}, err
	}

	if err := writeNDJSON(itemsPath, items); err != nil {
		return StageResult{}, err
	}

	if err := writeNDJSON(edgesPath, edges); err != nil {
		return StageResult{}, err
	}

	if err := writeNDJSON(citationsPath, citations); err != nil {
		return StageResult{}, err
	}

	manifest := provenance.Manifest{
		SchemaVersion:      provenance.SchemaVersion,
		Framework:          frameworkName,
		Version:            doc.Catalog.Metadata.Version,
		Stage:              "staged",
		GeneratedAt:        time.Now().UTC().Format(time.RFC3339),
		SourceTitle:        doc.Catalog.Metadata.Title,
		SourceLastModified: doc.Catalog.Metadata.LastModified,
	}

	rawArtifacts, err := buildArtifacts([]artifactSpec{
		{Path: opts.CatalogPath, Role: "catalog"},
		{Path: opts.LicensePath, Role: "license"},
	})
	if err != nil {
		return StageResult{}, err
	}
	manifest.RawArtifacts = rawArtifacts

	stagedArtifacts, err := buildArtifacts([]artifactSpec{
		{Path: metadataPath, Role: "metadata", RecordCount: 1},
		{Path: collectionsPath, Role: "collections", RecordCount: len(collections)},
		{Path: itemsPath, Role: "items", RecordCount: len(items)},
		{Path: edgesPath, Role: "edge-candidates", RecordCount: len(edges)},
		{Path: citationsPath, Role: "citations", RecordCount: len(citations)},
	})
	if err != nil {
		return StageResult{}, err
	}
	manifest.StagedArtifacts = stagedArtifacts

	if err := provenance.WriteFile(opts.ManifestPath, manifest); err != nil {
		return StageResult{}, err
	}

	return StageResult{
		Version:            doc.Catalog.Metadata.Version,
		ManifestPath:       filepath.ToSlash(filepath.Clean(opts.ManifestPath)),
		CollectionCount:    len(collections),
		ItemCount:          len(items),
		EdgeCandidateCount: len(edges),
		CitationCount:      len(citations),
	}, nil
}

func readCatalog(path string) (catalogDocument, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return catalogDocument{}, err
	}

	var doc catalogDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return catalogDocument{}, err
	}

	return doc, nil
}

func extract(doc catalogDocument) ([]Item, map[string]struct{}, []pendingEdge, []Citation) {
	var items []Item
	itemIDs := make(map[string]struct{})
	var pendingEdges []pendingEdge
	var citations []Citation

	addItem := func(item Item) {
		items = append(items, item)
		itemIDs[item.ItemID] = struct{}{}
	}

	for groupIndex, family := range doc.Catalog.Groups {
		familyPath := fmt.Sprintf("catalog.groups[%d]", groupIndex)
		addItem(Item{
			Framework:  frameworkName,
			Version:    doc.Catalog.Metadata.Version,
			ItemID:     strings.ToUpper(family.ID),
			ItemType:   "family",
			Class:      family.Class,
			Title:      family.Title,
			SourcePath: familyPath,
		})

		for controlIndex, topControl := range family.Controls {
			controlPath := fmt.Sprintf("%s.controls[%d]", familyPath, controlIndex)
			addControl(
				doc.Catalog.Metadata.Version,
				strings.ToUpper(family.ID),
				strings.ToUpper(family.ID),
				topControl,
				controlPath,
				&items,
				itemIDs,
				&pendingEdges,
				&citations,
			)

			pendingEdges = append(pendingEdges, pendingEdge{
				FromID:           strings.ToUpper(family.ID),
				RelationshipType: "has-child",
				ToRef:            "#" + strings.ToUpper(topControl.ID),
				SourcePath:       controlPath,
			})
		}
	}

	return items, itemIDs, pendingEdges, citations
}

func addControl(
	version string,
	familyID string,
	parentID string,
	src control,
	sourcePath string,
	items *[]Item,
	itemIDs map[string]struct{},
	pendingEdges *[]pendingEdge,
	citations *[]Citation,
) {
	itemID := strings.ToUpper(src.ID)
	item := Item{
		Framework:  frameworkName,
		Version:    version,
		ItemID:     itemID,
		ItemType:   itemTypeFor(src.Class, src.ID),
		Class:      src.Class,
		Title:      src.Title,
		ParentID:   parentID,
		FamilyID:   familyID,
		SourcePath: sourcePath,
		Props:      mapProps(src.Props),
	}

	*items = append(*items, item)
	itemIDs[item.ItemID] = struct{}{}

	for linkIndex, link := range src.Links {
		linkPath := fmt.Sprintf("%s.links[%d]", sourcePath, linkIndex)
		if link.Rel == "reference" {
			*citations = append(*citations, Citation{
				FromID:     itemID,
				Relation:   link.Rel,
				Href:       link.Href,
				Text:       link.Text,
				SourcePath: linkPath,
			})
			continue
		}

		*pendingEdges = append(*pendingEdges, pendingEdge{
			FromID:           itemID,
			RelationshipType: link.Rel,
			ToRef:            link.Href,
			SourcePath:       linkPath,
		})
	}

	for childIndex, child := range src.Controls {
		childPath := fmt.Sprintf("%s.controls[%d]", sourcePath, childIndex)
		addControl(version, familyID, itemID, child, childPath, items, itemIDs, pendingEdges, citations)
		*pendingEdges = append(*pendingEdges, pendingEdge{
			FromID:           itemID,
			RelationshipType: "has-child",
			ToRef:            "#" + strings.ToUpper(child.ID),
			SourcePath:       childPath,
		})
	}
}

func mapProps(props []prop) []Property {
	if len(props) == 0 {
		return nil
	}

	mapped := make([]Property, 0, len(props))
	for _, prop := range props {
		mapped = append(mapped, Property{
			Name:  prop.Name,
			Class: prop.Class,
			Value: prop.Value,
		})
	}

	return mapped
}

func itemTypeFor(class, id string) string {
	switch {
	case strings.Contains(strings.ToLower(class), "enhancement"):
		return "control-enhancement"
	case strings.Contains(id, "."):
		return "control-enhancement"
	default:
		return "control"
	}
}

func resolveEdges(itemIDs map[string]struct{}, pending []pendingEdge) []EdgeCandidate {
	edges := make([]EdgeCandidate, 0, len(pending))
	for _, edge := range pending {
		toKind := "external"
		toItemID := ""
		if strings.HasPrefix(edge.ToRef, "#") {
			ref := strings.ToUpper(strings.TrimPrefix(edge.ToRef, "#"))
			if _, ok := itemIDs[ref]; ok {
				toKind = "item"
				toItemID = ref
			} else {
				toKind = "fragment"
			}
		}

		edges = append(edges, EdgeCandidate{
			FromID:           edge.FromID,
			RelationshipType: edge.RelationshipType,
			ToRef:            edge.ToRef,
			ToKind:           toKind,
			ToItemID:         toItemID,
			SourcePath:       edge.SourcePath,
		})
	}

	return edges
}

func writeJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

func writeNDJSON[T any](path string, records []T) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	for _, record := range records {
		if err := encoder.Encode(record); err != nil {
			return err
		}
	}

	return writer.Flush()
}

type artifactSpec struct {
	Path        string
	Role        string
	RecordCount int
}

func buildArtifacts(specs []artifactSpec) ([]provenance.Artifact, error) {
	artifacts := make([]provenance.Artifact, 0, len(specs))
	for _, spec := range specs {
		artifact, err := provenance.BuildArtifact(spec.Path, spec.Role, spec.RecordCount)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}

	return artifacts, nil
}
