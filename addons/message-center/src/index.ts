import type { ToolModule } from "@sstpa/addon-sdk"

const messageCenterTool: ToolModule = {
  manifest: {
    id: "message-center",
    name: "Message Center",
    runtime: "popup",
    summary: "Mailbox access for direct messages and owner-notification traffic.",
    inputContracts: ["current-user", "current-soi"],
    outputContracts: ["read-state-change", "reply-draft"],
  },
  load: async () => ({ component: "message-center" }),
}

export default messageCenterTool
