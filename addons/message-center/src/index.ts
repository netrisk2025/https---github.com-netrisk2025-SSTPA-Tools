// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
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
