// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import type { ToolModule } from "@sstpa/addon-sdk"

const goalKeeperTool: ToolModule = {
  manifest: {
    id: "goal-keeper",
    name: "Goal Keeper Tool",
    runtime: "popup",
    summary: "GSN assurance-case DAG editing, validation, evidence association, and layout persistence.",
    inputContracts: ["current-soi", "data-drawer", "selected-goal-structure"],
    outputContracts: ["gsn-commit", "diagram-layout-commit", "evidence-selection"],
    requiredContexts: ["current-soi", "current-user", "data-drawer"],
    graphScopes: [
      {
        nodeTypes: ["Goal", "Strategy", "Context", "Justification", "Assumption", "Solution"],
        relationshipTypes: ["SUPPORTED_BY", "IN_CONTEXT_OF", "HAS_VALIDATION", "HAS_VERIFICATION", "HAS_LOSS"],
        allowsCrossSoI: false,
      },
    ],
    diagramPersistence: {
      storageProperty: "GoalStructure",
      authoritativeSource: "neo4j-graph",
      reconcilesStaleReferences: true,
    },
    requirementIds: [
      "3.4.11.1-001",
      "3.4.11.1-003",
      "3.4.11.1-005",
      "3.4.11.5-002",
      "3.4.11.6-001",
      "3.4.11.9-001",
      "3.4.11.12-001",
    ],
  },
  load: async () => ({ component: "goal-keeper" }),
}

export default goalKeeperTool
