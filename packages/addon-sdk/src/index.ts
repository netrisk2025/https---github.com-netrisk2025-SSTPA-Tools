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
