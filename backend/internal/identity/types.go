// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package identity

type NodeType string

const (
	NodeTypeCapability        NodeType = "Capability"
	NodeTypeSandbox           NodeType = "Sandbox"
	NodeTypeSystem            NodeType = "System"
	NodeTypeEnvironment       NodeType = "Environment"
	NodeTypeConnection        NodeType = "Connection"
	NodeTypeInterface         NodeType = "Interface"
	NodeTypeFunction          NodeType = "Function"
	NodeTypeElement           NodeType = "Element"
	NodeTypePurpose           NodeType = "Purpose"
	NodeTypeState             NodeType = "State"
	NodeTypeControlStructure  NodeType = "ControlStructure"
	NodeTypeAsset             NodeType = "Asset"
	NodeTypeSecurity          NodeType = "Security"
	NodeTypeFunctionalFlow    NodeType = "FunctionalFlow"
	NodeTypeRegime            NodeType = "Regime"
	NodeTypeConstraint        NodeType = "Constraint"
	NodeTypeRequirement       NodeType = "Requirement"
	NodeTypeValidation        NodeType = "Validation"
	NodeTypeControl           NodeType = "Control"
	NodeTypeCountermeasure    NodeType = "Countermeasure"
	NodeTypeVerification      NodeType = "Verification"
	NodeTypeControlAlgorithm  NodeType = "ControlAlgorithm"
	NodeTypeProcessModel      NodeType = "ProcessModel"
	NodeTypeControlAction     NodeType = "ControlAction"
	NodeTypeFeedback          NodeType = "Feedback"
	NodeTypeControlledProcess NodeType = "ControlledProcess"
	NodeTypeHazard            NodeType = "Hazard"
	NodeTypeLoss              NodeType = "Loss"
	NodeTypeAttack            NodeType = "Attack"
	NodeTypeGoal              NodeType = "Goal"
	NodeTypeStrategy          NodeType = "Strategy"
	NodeTypeContext           NodeType = "Context"
	NodeTypeJustification     NodeType = "Justification"
	NodeTypeAssumption        NodeType = "Assumption"
	NodeTypeSolution          NodeType = "Solution"
	NodeTypeSSTPATool         NodeType = "SSTPA_Tool"
	NodeTypeUserRegistry      NodeType = "UserRegistry"
	NodeTypeAdminRegistry     NodeType = "AdminRegistry"
	NodeTypeUser              NodeType = "User"
	NodeTypeAdmin             NodeType = "Admin"
	NodeTypeMasterRegime      NodeType = "MasterRegime"
)

var typeIDs = map[NodeType]string{
	NodeTypeCapability:        "CAP",
	NodeTypeSandbox:           "SB",
	NodeTypeSystem:            "SYS",
	NodeTypeEnvironment:       "ENV",
	NodeTypeConnection:        "CNN",
	NodeTypeInterface:         "INT",
	NodeTypeFunction:          "FUN",
	NodeTypeElement:           "EL",
	NodeTypePurpose:           "PUR",
	NodeTypeState:             "ST",
	NodeTypeControlStructure:  "CS",
	NodeTypeAsset:             "AST",
	NodeTypeSecurity:          "SEC",
	NodeTypeFunctionalFlow:    "FF",
	NodeTypeRegime:            "REG",
	NodeTypeConstraint:        "CONSTR",
	NodeTypeRequirement:       "REQ",
	NodeTypeValidation:        "VAL",
	NodeTypeControl:           "CTRL",
	NodeTypeCountermeasure:    "CM",
	NodeTypeVerification:      "VER",
	NodeTypeControlAlgorithm:  "CAL",
	NodeTypeProcessModel:      "PM",
	NodeTypeControlAction:     "ACT",
	NodeTypeFeedback:          "FB",
	NodeTypeControlledProcess: "CP",
	NodeTypeHazard:            "HAZ",
	NodeTypeLoss:              "LOS",
	NodeTypeAttack:            "ATK",
	NodeTypeGoal:              "G",
	NodeTypeStrategy:          "SGY",
	NodeTypeContext:           "CX",
	NodeTypeAssumption:        "ASS",
	NodeTypeJustification:     "JUS",
	NodeTypeSolution:          "SOL",
	NodeTypeSSTPATool:         "SST",
	NodeTypeUserRegistry:      "URG",
	NodeTypeAdminRegistry:     "ARG",
	NodeTypeUser:              "USR",
	NodeTypeAdmin:             "ADM",
	NodeTypeMasterRegime:      "MRG",
}

var orderedTypes = []NodeType{
	NodeTypeCapability,
	NodeTypeSandbox,
	NodeTypeSystem,
	NodeTypeEnvironment,
	NodeTypeConnection,
	NodeTypeInterface,
	NodeTypeFunction,
	NodeTypeElement,
	NodeTypePurpose,
	NodeTypeState,
	NodeTypeControlStructure,
	NodeTypeAsset,
	NodeTypeSecurity,
	NodeTypeFunctionalFlow,
	NodeTypeRegime,
	NodeTypeConstraint,
	NodeTypeRequirement,
	NodeTypeValidation,
	NodeTypeControl,
	NodeTypeCountermeasure,
	NodeTypeVerification,
	NodeTypeControlAlgorithm,
	NodeTypeProcessModel,
	NodeTypeControlAction,
	NodeTypeFeedback,
	NodeTypeControlledProcess,
	NodeTypeHazard,
	NodeTypeLoss,
	NodeTypeAttack,
	NodeTypeGoal,
	NodeTypeStrategy,
	NodeTypeContext,
	NodeTypeJustification,
	NodeTypeAssumption,
	NodeTypeSolution,
	NodeTypeSSTPATool,
	NodeTypeUserRegistry,
	NodeTypeAdminRegistry,
	NodeTypeUser,
	NodeTypeAdmin,
	NodeTypeMasterRegime,
}

func TypeID(nodeType NodeType) (string, bool) {
	id, ok := typeIDs[nodeType]
	return id, ok
}

func AllTypes() []NodeType {
	return append([]NodeType(nil), orderedTypes...)
}

func IsValidTypeID(id string) bool {
	for _, candidate := range typeIDs {
		if candidate == id {
			return true
		}
	}

	return false
}
