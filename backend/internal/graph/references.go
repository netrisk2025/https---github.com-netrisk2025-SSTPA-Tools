// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package graph

import (
	"strings"

	"sstpa-tool/backend/internal/identity"
)

type ReferenceAssignment struct {
	SourceType    identity.NodeType
	FrameworkName string
	ExternalType  string
}

func ReferenceAssignmentAllowed(sourceType identity.NodeType, frameworkName string, externalType string) bool {
	normalizedFramework := normalizeReferenceText(frameworkName)
	normalizedType := normalizeReferenceText(externalType)
	for _, assignment := range referenceAssignmentCatalog {
		if assignment.SourceType != sourceType {
			continue
		}
		if normalizeReferenceText(assignment.FrameworkName) != normalizedFramework {
			continue
		}
		if normalizeReferenceText(assignment.ExternalType) == normalizedType {
			return true
		}
	}

	return false
}

func ReferenceAssignmentCatalog() []ReferenceAssignment {
	return append([]ReferenceAssignment(nil), referenceAssignmentCatalog...)
}

func normalizeReferenceText(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(value)), " ")
}

var referenceAssignmentCatalog = []ReferenceAssignment{
	{SourceType: identity.NodeTypeControl, FrameworkName: "NIST SP 800-53", ExternalType: "Control"},
	{SourceType: identity.NodeTypeControl, FrameworkName: "NIST SP800-53", ExternalType: "Control"},
	{SourceType: identity.NodeTypeControl, FrameworkName: "MITRE ATT&CK", ExternalType: "Mitigation"},
	{SourceType: identity.NodeTypeControl, FrameworkName: "MITRE EMB3D", ExternalType: "Mitigation"},

	{SourceType: identity.NodeTypeElement, FrameworkName: "MITRE EMB3D", ExternalType: "Property"},
	{SourceType: identity.NodeTypeElement, FrameworkName: "MITRE EMB3D", ExternalType: "Threat"},

	{SourceType: identity.NodeTypeSystem, FrameworkName: "MITRE EMB3D", ExternalType: "Property"},
	{SourceType: identity.NodeTypeSystem, FrameworkName: "NIST SP 800-53", ExternalType: "Control"},
	{SourceType: identity.NodeTypeSystem, FrameworkName: "NIST SP800-53", ExternalType: "Control"},

	{SourceType: identity.NodeTypeHazard, FrameworkName: "MITRE EMB3D", ExternalType: "Threat"},
	{SourceType: identity.NodeTypeHazard, FrameworkName: "MITRE ATT&CK", ExternalType: "Technique"},

	{SourceType: identity.NodeTypeAttack, FrameworkName: "MITRE ATT&CK", ExternalType: "Tactic"},
	{SourceType: identity.NodeTypeAttack, FrameworkName: "MITRE ATT&CK", ExternalType: "Technique"},
	{SourceType: identity.NodeTypeAttack, FrameworkName: "MITRE EMB3D", ExternalType: "Threat"},

	{SourceType: identity.NodeTypeCountermeasure, FrameworkName: "MITRE EMB3D", ExternalType: "Mitigation"},
	{SourceType: identity.NodeTypeCountermeasure, FrameworkName: "MITRE ATT&CK", ExternalType: "Mitigation"},
	{SourceType: identity.NodeTypeCountermeasure, FrameworkName: "NIST SP 800-53", ExternalType: "Control"},
	{SourceType: identity.NodeTypeCountermeasure, FrameworkName: "NIST SP800-53", ExternalType: "Control"},
}
