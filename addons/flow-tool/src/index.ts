// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import type { ToolModule } from "@sstpa/addon-sdk"

const flowTool: ToolModule = {
  manifest: {
    id: "flow",
    name: "Flow Tool",
    runtime: "popup",
    summary: "Functional flow and STPA control-flow diagram work for the current SoI.",
    inputContracts: ["current-soi", "data-drawer", "selected-flow"],
    outputContracts: ["functional-flow-commit", "control-flow-commit", "diagram-layout-commit"],
    requiredContexts: ["current-soi", "current-user", "data-drawer"],
    graphScopes: [
      {
        nodeTypes: ["FunctionalFlow", "Function", "Interface", "Connection", "Element", "Asset"],
        relationshipTypes: ["FLOWS_TO_FUNCTION", "FLOWS_TO_INTERFACE", "CONTAINS", "PARTICIPATES_IN", "CONNECTS"],
        allowsCrossSoI: false,
      },
      {
        nodeTypes: ["ControlStructure", "ControlAlgorithm", "ProcessModel", "ControlledProcess", "ControlAction", "Feedback"],
        relationshipTypes: ["IMPLEMENTS", "GENERATES", "COMMANDS", "PRODUCES", "INFORMS", "TUNES"],
        allowsCrossSoI: false,
      },
    ],
    diagramPersistence: {
      storageProperty: "FunctionalFlowJSON",
      authoritativeSource: "neo4j-graph",
      reconcilesStaleReferences: true,
    },
    requirementIds: ["1.3.4.5-001", "1.3.4.12-001", "3.4.6.1-001", "3.4.6.2-002"],
  },
  load: async () => ({ component: "flow" }),
}

export default flowTool
