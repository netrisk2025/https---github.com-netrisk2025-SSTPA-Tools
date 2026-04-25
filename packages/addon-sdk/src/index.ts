// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
export type ToolRuntime = "panel" | "popup"

export type ToolContextKind =
  | "current-soi"
  | "current-user"
  | "data-drawer"
  | "selected-node"
  | "selected-relationship"
  | "graph-selection"

export type DiagramStorageProperty =
  | "ControlStructureJSON"
  | "FunctionalFlowJSON"
  | "AttackTreeJSON"
  | "GoalStructure"

export interface NodeRef {
  hid: string
  uuid: string
  typeName: string
  name?: string
}

export interface CurrentSoIContext {
  system: NodeRef
  hidIndex: string
}

export interface DataDrawerContext {
  node?: NodeRef
  relationshipType?: string
  mode: "view" | "edit" | "associate"
}

export interface DiagramViewState {
  schemaVersion: string
  rootHid?: string
  rootUuid?: string
  viewport: { x: number; y: number }
  zoom: number
  nodePositions: Record<string, { x: number; y: number }>
  collapsed: string[]
  filters: Record<string, string | number | boolean>
  displayOptions: Record<string, string | number | boolean>
  staleReferences?: string[]
}

export interface DiagramPersistenceContract {
  storageProperty: DiagramStorageProperty
  authoritativeSource: "neo4j-graph"
  reconcilesStaleReferences: boolean
}

export interface GraphScopeContract {
  nodeTypes: string[]
  relationshipTypes: string[]
  allowsCrossSoI: boolean
}

export interface ToolManifest {
  id: string
  name: string
  runtime: ToolRuntime
  summary: string
  inputContracts: string[]
  outputContracts: string[]
  requiredContexts?: ToolContextKind[]
  graphScopes?: GraphScopeContract[]
  diagramPersistence?: DiagramPersistenceContract
  requirementIds?: string[]
}

export interface ToolModule {
  manifest: ToolManifest
  load: () => Promise<unknown>
}
