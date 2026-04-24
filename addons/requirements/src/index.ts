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
