// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import messageCenterTool from "@sstpa/message-center-tool"
import navigatorTool from "@sstpa/navigator-tool"
import requirementsTool from "@sstpa/requirements-tool"
import assetManagerTool from "@sstpa/asset-manager-tool"
import flowTool from "@sstpa/flow-tool"
import goalKeeperTool from "@sstpa/goal-keeper-tool"
import lossTool from "@sstpa/loss-tool"
import referenceCatalogTool from "@sstpa/reference-catalog-tool"
import stateTool from "@sstpa/state-tool"

export const toolRegistry = [
  navigatorTool,
  requirementsTool,
  stateTool,
  flowTool,
  assetManagerTool,
  lossTool,
  goalKeeperTool,
  referenceCatalogTool,
  messageCenterTool,
]
