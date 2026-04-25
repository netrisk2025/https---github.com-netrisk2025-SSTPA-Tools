// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import type { ToolModule } from "@sstpa/addon-sdk"

const lossTool: ToolModule = {
  manifest: {
    id: "loss",
    name: "Loss Tool",
    runtime: "popup",
    summary: "AttackTreeJSON viewing, editing, validation, and Loss analysis for the selected SoI.",
    inputContracts: ["current-soi", "data-drawer", "selected-loss"],
    outputContracts: ["attack-tree-commit", "derived-asset-request", "residual-risk-update"],
    requiredContexts: ["current-soi", "current-user", "data-drawer"],
    graphScopes: [
      {
        nodeTypes: ["Loss", "Attack", "Asset", "Countermeasure", "Environment", "State", "Element"],
        relationshipTypes: ["HAS_ATTACK", "HAS_COUNTERMEASURE", "HAS_ENVIRONMENT", "HAS_STATE", "HAS_ELEMENT"],
        allowsCrossSoI: false,
      },
    ],
    diagramPersistence: {
      storageProperty: "AttackTreeJSON",
      authoritativeSource: "neo4j-graph",
      reconcilesStaleReferences: true,
    },
    requirementIds: ["1.3.1.5-002", "1.3.4.11-001", "3.4.7.11-003", "3.4.10.1-001", "3.4.10.1-002"],
  },
  load: async () => ({ component: "loss" }),
}

export default lossTool
