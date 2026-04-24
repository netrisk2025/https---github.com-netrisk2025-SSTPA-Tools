// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"sstpa-tool/backend/internal/identity"
)

type nodeResponse struct {
	HID           string         `json:"hid"`
	UUID          string         `json:"uuid"`
	TypeName      string         `json:"typeName"`
	ContainingSoI string         `json:"containingSoI"`
	Properties    map[string]any `json:"properties"`
}

type relationshipResponse struct {
	FromHID string `json:"fromHid"`
	ToHID   string `json:"toHid"`
	Type    string `json:"type"`
}

type hierarchyResponse struct {
	Nodes         []nodeSummary          `json:"nodes"`
	Relationships []relationshipResponse `json:"relationships"`
}

type nodeSummary struct {
	HID              string `json:"hid"`
	UUID             string `json:"uuid"`
	TypeName         string `json:"typeName"`
	Name             string `json:"name"`
	ShortDescription string `json:"shortDescription"`
	ContainingSoI    string `json:"containingSoI"`
}

type nodeContextResponse struct {
	Node                nodeResponse           `json:"node"`
	ContainingSoI       string                 `json:"containingSoI"`
	ParentRelationships []relationshipResponse `json:"parentRelationships"`
	Path                []string               `json:"path"`
}

func (api api) getNodeByHIDHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	node, err := api.findNode(request.Context(), "n.HID = $value", chi.URLParam(request, "hid"))
	if err != nil {
		handleNeo4jReadError(writer, err, "node not found")
		return
	}

	writeJSON(writer, http.StatusOK, node)
}

func (api api) getNodeByUUIDHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	node, err := api.findNode(request.Context(), "n.uuid = $value", chi.URLParam(request, "uuid"))
	if err != nil {
		handleNeo4jReadError(writer, err, "node not found")
		return
	}

	writeJSON(writer, http.StatusOK, node)
}

func (api api) listNodesHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	page, err := parsePagination(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	nodeType, err := parseOptionalNodeType(request.URL.Query().Get("type"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	result, err := api.searchNodes(request.Context(), searchNodeInput{
		Query:    "",
		NodeType: nodeType,
		Page:     page,
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, result)
}

func (api api) searchHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	page, err := parsePagination(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	nodeType, err := parseOptionalNodeType(request.URL.Query().Get("type"))
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	result, err := api.searchNodes(request.Context(), searchNodeInput{
		Query:    request.URL.Query().Get("q"),
		NodeType: nodeType,
		Page:     page,
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, result)
}

func (api api) hierarchyHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	session := api.driver.NewSession(request.Context(), neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(request.Context())

	value, err := session.ExecuteRead(request.Context(), func(tx neo4j.ManagedTransaction) (any, error) {
		nodeRecords, err := tx.Run(request.Context(), `
MATCH (n:SSTPANode)
WHERE n:Capability OR n:Sandbox OR n:System OR n:Element
RETURN properties(n) AS properties
ORDER BY n.HID
`, nil)
		if err != nil {
			return nil, err
		}
		records, err := nodeRecords.Collect(request.Context())
		if err != nil {
			return nil, err
		}

		nodes := make([]nodeSummary, 0, len(records))
		for _, record := range records {
			props := propertiesFromRecord(record, "properties")
			nodes = append(nodes, nodeSummaryFromProperties(props))
		}

		relationshipRecords, err := tx.Run(request.Context(), `
MATCH (parent)-[r:HAS_SYSTEM|PARENTS]->(child:System)
WHERE parent:Capability OR parent:Sandbox OR parent:Element
RETURN parent.HID AS fromHID, child.HID AS toHID, type(r) AS type
ORDER BY fromHID, toHID
`, nil)
		if err != nil {
			return nil, err
		}
		relRecords, err := relationshipRecords.Collect(request.Context())
		if err != nil {
			return nil, err
		}

		relationships := make([]relationshipResponse, 0, len(relRecords))
		for _, record := range relRecords {
			relationships = append(relationships, relationshipFromRecord(record))
		}

		return hierarchyResponse{Nodes: nodes, Relationships: relationships}, nil
	})
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, value)
}

func (api api) nodeContextHandler(writer http.ResponseWriter, request *http.Request) {
	if !api.requireDriver(writer) {
		return
	}

	hid := chi.URLParam(request, "hid")
	session := api.driver.NewSession(request.Context(), neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(request.Context())

	value, err := session.ExecuteRead(request.Context(), func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(request.Context(), `
MATCH (n:SSTPANode {HID: $hid})
OPTIONAL MATCH (parent)-[r:HAS_SYSTEM|PARENTS|HAS_ELEMENT|HAS_FUNCTION|HAS_INTERFACE|HAS_REQUIREMENT|HAS_HAZARD|HAS_CONTROL|HAS_COUNTERMEASURE]->(n)
OPTIONAL MATCH path = (root:Capability)-[*1..10]->(n)
RETURN properties(n) AS properties,
       collect(DISTINCT {fromHID: parent.HID, toHID: n.HID, type: type(r)}) AS parents,
       [node IN nodes(path) | node.HID] AS path
LIMIT 1
`, map[string]any{"hid": hid})
		if err != nil {
			return nil, err
		}
		record, err := result.Single(request.Context())
		if err != nil {
			return nil, err
		}

		props := propertiesFromRecord(record, "properties")
		node := nodeFromProperties(props)
		parentsValue, _ := record.Get("parents")
		pathValue, _ := record.Get("path")
		return nodeContextResponse{
			Node:                node,
			ContainingSoI:       node.ContainingSoI,
			ParentRelationships: relationshipMapsToResponses(parentsValue),
			Path:                anySliceToStrings(pathValue),
		}, nil
	})
	if err != nil {
		handleNeo4jReadError(writer, err, "node context not found")
		return
	}

	writeJSON(writer, http.StatusOK, value)
}

type searchNodeInput struct {
	Query    string
	NodeType identity.NodeType
	Page     pageRequest
}

func (api api) findNode(ctx context.Context, whereClause string, value string) (nodeResponse, error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
MATCH (n:SSTPANode)
WHERE ` + whereClause + `
RETURN properties(n) AS properties
LIMIT 1
`
		records, err := tx.Run(ctx, query, map[string]any{"value": value})
		if err != nil {
			return nodeResponse{}, err
		}

		record, err := records.Single(ctx)
		if err != nil {
			return nodeResponse{}, err
		}

		return nodeFromProperties(propertiesFromRecord(record, "properties")), nil
	})
	if err != nil {
		return nodeResponse{}, err
	}

	node, ok := result.(nodeResponse)
	if !ok {
		return nodeResponse{}, errors.New("unexpected node response")
	}

	return node, nil
}

func (api api) searchNodes(ctx context.Context, input searchNodeInput) (listResponse[nodeSummary], error) {
	session := api.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: api.databaseName})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"query":    input.Query,
			"lower":    strings.ToLower(input.Query),
			"typeName": string(input.NodeType),
			"skip":     input.Page.Offset,
			"limit":    input.Page.Limit,
		}
		where := `
($typeName = "" OR n.TypeName = $typeName)
AND (
  $query = ""
  OR n.HID = $query
  OR n.uuid = $query
  OR toLower(coalesce(n.Name, "")) CONTAINS $lower
  OR toLower(coalesce(n.ShortDescription, "")) CONTAINS $lower
)
`
		records, err := tx.Run(ctx, `
MATCH (n:SSTPANode)
WHERE `+where+`
RETURN properties(n) AS properties
ORDER BY n.HID
SKIP $skip
LIMIT $limit
`, params)
		if err != nil {
			return nil, err
		}
		collected, err := records.Collect(ctx)
		if err != nil {
			return nil, err
		}

		items := make([]nodeSummary, 0, len(collected))
		for _, record := range collected {
			items = append(items, nodeSummaryFromProperties(propertiesFromRecord(record, "properties")))
		}

		countRecords, err := tx.Run(ctx, `
MATCH (n:SSTPANode)
WHERE `+where+`
RETURN count(n) AS total
`, params)
		if err != nil {
			return nil, err
		}
		countRecord, err := countRecords.Single(ctx)
		if err != nil {
			return nil, err
		}
		totalValue, _ := countRecord.Get("total")
		total, _ := totalValue.(int64)

		return listResponse[nodeSummary]{Items: items, Page: input.Page.Page, Limit: input.Page.Limit, Total: total}, nil
	})
	if err != nil {
		return listResponse[nodeSummary]{}, err
	}

	response, ok := result.(listResponse[nodeSummary])
	if !ok {
		return listResponse[nodeSummary]{}, errors.New("unexpected search response")
	}

	return response, nil
}

func parseOptionalNodeType(value string) (identity.NodeType, error) {
	if value == "" {
		return "", nil
	}
	nodeType := identity.NodeType(value)
	if _, ok := identity.TypeID(nodeType); !ok {
		return "", errors.New("unknown node type")
	}

	return nodeType, nil
}

func handleNeo4jReadError(writer http.ResponseWriter, err error, notFoundMessage string) {
	if isNoRecordsError(err) {
		writeError(writer, http.StatusNotFound, notFoundMessage)
		return
	}

	writeError(writer, http.StatusInternalServerError, err.Error())
}

func isNoRecordsError(err error) bool {
	return neo4j.IsUsageError(err) && strings.Contains(err.Error(), "Result contains no more records")
}

func propertiesFromRecord(record *neo4j.Record, key string) map[string]any {
	value, _ := record.Get(key)
	properties, _ := value.(map[string]any)
	if properties == nil {
		return map[string]any{}
	}

	return properties
}

func nodeFromProperties(properties map[string]any) nodeResponse {
	node := nodeResponse{
		HID:           stringProperty(properties, "HID"),
		UUID:          stringProperty(properties, "uuid"),
		TypeName:      stringProperty(properties, "TypeName"),
		ContainingSoI: stringProperty(properties, "SoI"),
		Properties:    properties,
	}
	if node.ContainingSoI == "" && node.TypeName == string(identity.NodeTypeSystem) {
		node.ContainingSoI = node.HID
	}

	return node
}

func nodeSummaryFromProperties(properties map[string]any) nodeSummary {
	node := nodeFromProperties(properties)
	return nodeSummary{
		HID:              node.HID,
		UUID:             node.UUID,
		TypeName:         node.TypeName,
		Name:             stringProperty(properties, "Name"),
		ShortDescription: stringProperty(properties, "ShortDescription"),
		ContainingSoI:    node.ContainingSoI,
	}
}

func relationshipFromRecord(record *neo4j.Record) relationshipResponse {
	fromHID, _ := record.Get("fromHID")
	toHID, _ := record.Get("toHID")
	relationshipType, _ := record.Get("type")
	return relationshipResponse{
		FromHID: stringValue(fromHID),
		ToHID:   stringValue(toHID),
		Type:    stringValue(relationshipType),
	}
}

func relationshipMapsToResponses(value any) []relationshipResponse {
	items, _ := value.([]any)
	responses := make([]relationshipResponse, 0, len(items))
	for _, item := range items {
		relationshipMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		relationship := relationshipResponse{
			FromHID: stringValue(relationshipMap["fromHID"]),
			ToHID:   stringValue(relationshipMap["toHID"]),
			Type:    stringValue(relationshipMap["type"]),
		}
		if relationship.FromHID == "" || relationship.ToHID == "" || relationship.Type == "" {
			continue
		}
		responses = append(responses, relationship)
	}

	return responses
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}
