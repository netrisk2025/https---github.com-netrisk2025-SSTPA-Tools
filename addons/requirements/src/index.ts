// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import type { ToolModule } from "@sstpa/addon-sdk"

const requirementsTool: ToolModule = {
  manifest: {
    id: "requirements",
    name: "Requirements Tool",
    runtime: "popup",
    summary: "Requirement hierarchy visualization, editing, and traceability work.",
    inputContracts: ["selected-requirement", "current-soi"],
    outputContracts: ["updated-requirement", "traceability-request"],
    requiredContexts: ["current-soi", "current-user", "data-drawer"],
    graphScopes: [
      {
        nodeTypes: ["Requirement", "Capability", "Purpose", "Connection", "Interface", "Function", "Element", "Constraint", "Countermeasure"],
        relationshipTypes: ["HAS_REQUIREMENT", "PARENTS", "VERIFIED_BY"],
        allowsCrossSoI: true,
      },
    ],
    requirementIds: ["1.3.4.7-001", "1.3.4.8-001", "3.4.2.2-001", "3.4.2.16-001"],
  },
  load: async () => ({ component: "requirements" }),
}

export default requirementsTool
