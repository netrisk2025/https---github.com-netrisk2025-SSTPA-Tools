// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
export type ToolRuntime = "panel" | "popup"

export interface ToolManifest {
  id: string
  name: string
  runtime: ToolRuntime
  summary: string
  inputContracts: string[]
  outputContracts: string[]
}

export interface ToolModule {
  manifest: ToolManifest
  load: () => Promise<unknown>
}
