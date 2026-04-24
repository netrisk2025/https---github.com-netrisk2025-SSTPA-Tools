import type { ToolModule } from "@sstpa/addon-sdk"

const navigatorTool: ToolModule = {
  manifest: {
    id: "navigator",
    name: "Navigator Tool",
    runtime: "popup",
    summary: "Hierarchy search, SoI selection, and clone-target selection.",
    inputContracts: ["current-soi", "search-query", "selection-mode"],
    outputContracts: ["selected-soi", "selected-node", "clone-target"],
  },
  load: async () => ({ component: "navigator" }),
}

export default navigatorTool
