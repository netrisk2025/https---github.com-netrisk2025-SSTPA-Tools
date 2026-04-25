// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import type { ToolModule } from "@sstpa/addon-sdk"

const referenceCatalogTool: ToolModule = {
  manifest: {
    id: "reference-catalog",
    name: "Reference Catalog Tool",
    runtime: "popup",
    summary: "Read-only reference framework research and REFERENCES assignment workflows.",
    inputContracts: ["current-soi", "data-drawer", "reference-search-query"],
    outputContracts: ["reference-assignment", "reference-removal", "reference-selection"],
    requiredContexts: ["current-soi", "current-user", "data-drawer"],
    graphScopes: [
      {
        nodeTypes: ["ReferenceItem", "Element", "Hazard", "Attack", "Control", "Countermeasure"],
        relationshipTypes: ["REFERENCES"],
        allowsCrossSoI: true,
      },
    ],
    requirementIds: ["1.3.10-005", "1.5.6-003", "3.4.4.1-001", "3.4.4.3-001"],
  },
  load: async () => ({ component: "reference-catalog" }),
}

export default referenceCatalogTool
