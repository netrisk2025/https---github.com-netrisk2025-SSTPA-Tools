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
  },
  load: async () => ({ component: "requirements" }),
}

export default requirementsTool
