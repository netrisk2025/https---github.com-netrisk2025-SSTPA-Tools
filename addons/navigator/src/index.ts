// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import type { ToolModule } from "@sstpa/addon-sdk"

const navigatorTool: ToolModule = {
  manifest: {
    id: "navigator",
    name: "Navigator Tool",
    runtime: "popup",
    summary: "Hierarchy search, SoI selection, and clone-target selection.",
    inputContracts: ["current-soi", "search-query", "selection-mode"],
    outputContracts: ["selected-soi", "selected-node", "clone-target"],
    requiredContexts: ["current-soi", "current-user", "graph-selection"],
    graphScopes: [
      {
        nodeTypes: ["Capability", "Sandbox", "System", "Element"],
        relationshipTypes: ["HAS_SYSTEM", "PARENTS"],
        allowsCrossSoI: true,
      },
    ],
    requirementIds: ["3.4.1-001", "3.4.1-002", "3.4.1.1-002", "3.4.1.1-006"],
  },
  load: async () => ({ component: "navigator" }),
}

export default navigatorTool
