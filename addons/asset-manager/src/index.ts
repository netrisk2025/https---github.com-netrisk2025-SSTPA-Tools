// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import type { ToolModule } from "@sstpa/addon-sdk"

const assetManagerTool: ToolModule = {
  manifest: {
    id: "asset-manager",
    name: "Asset Manager Tool",
    runtime: "popup",
    summary: "Table-first Asset, Regime, Loss, and Root Goal management for the current SoI.",
    inputContracts: ["current-soi", "data-drawer", "selected-node"],
    outputContracts: ["asset-commit", "regime-clone", "loss-goal-generation-request"],
    requiredContexts: ["current-soi", "current-user", "data-drawer"],
    graphScopes: [
      {
        nodeTypes: ["Asset", "Regime", "MasterRegime", "Loss", "Goal", "Environment"],
        relationshipTypes: ["HAS_REGIME", "DERIVED_FROM", "HAS_LOSS", "HAS_GOAL", "HAS_ENVIRONMENT"],
        allowsCrossSoI: false,
      },
    ],
    requirementIds: [
      "3.4.7.1-001",
      "3.4.7.1-003",
      "3.4.7.6-001",
      "3.4.7.7-001",
      "3.4.7.8-005",
      "3.4.7.10-001",
    ],
  },
  load: async () => ({ component: "asset-manager" }),
}

export default assetManagerTool
