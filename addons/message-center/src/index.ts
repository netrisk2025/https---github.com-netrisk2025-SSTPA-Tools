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
    requiredContexts: ["current-user"],
    graphScopes: [
      {
        nodeTypes: ["User", "Mailbox", "Message"],
        relationshipTypes: ["HAS_MAILBOX", "HAS_MESSAGE", "REPLIES_TO"],
        allowsCrossSoI: true,
      },
    ],
    requirementIds: ["3.3.1-001", "3.3.1-002", "3.3.1-003", "3.3.1-004", "3.3.1-005"],
  },
  load: async () => ({ component: "message-center" }),
}

export default messageCenterTool
