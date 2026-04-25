// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import type { ToolModule } from "@sstpa/addon-sdk"

const stateTool: ToolModule = {
  manifest: {
    id: "state",
    name: "State Tool",
    runtime: "popup",
    summary: "State transition observation, analysis, editing, and export for the current SoI.",
    inputContracts: ["current-soi", "data-drawer", "selected-state"],
    outputContracts: ["state-transition-commit", "state-diagram-export"],
    requiredContexts: ["current-soi", "current-user", "data-drawer"],
    graphScopes: [
      {
        nodeTypes: ["State", "Hazard", "Countermeasure", "Requirement"],
        relationshipTypes: ["TRANSITIONS_TO", "HAS_HAZARD", "APPLIES_TO_STATE", "HAS_REQUIREMENT"],
        allowsCrossSoI: false,
      },
    ],
    requirementIds: ["1.3.4.3-001", "1.3.6-004", "3.4.5.1-001"],
  },
  load: async () => ({ component: "state" }),
}

export default stateTool
