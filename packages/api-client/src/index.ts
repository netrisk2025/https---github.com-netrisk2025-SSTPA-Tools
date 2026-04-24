// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
export interface SSTPAClientOptions {
  baseUrl?: string
  fetchImpl?: typeof fetch
  actor?: Actor
}

export interface Actor {
  name: string
  email: string
  admin?: boolean
}

export interface NodeSummary {
  hid: string
  uuid: string
  typeName: string
  name: string
  shortDescription: string
  containingSoI: string
}

export interface NodeDetail extends NodeSummary {
  properties: Record<string, unknown>
}

export interface ListResponse<TItem> {
  items: TItem[]
  page: number
  limit: number
  total: number
}

export interface RelationshipValidationRequest {
  relationshipName: string
  fromType: string
  toType: string
  fromHid?: string
  toHid?: string
}

export interface ValidationResponse {
  valid: boolean
  reason: string
}

export type MutationOperationKind = "create_node" | "update_node" | "create_relationship"

export interface MutationOperation {
  kind: MutationOperationKind
  nodeType?: string
  hid?: string
  uuid?: string
  properties?: Record<string, unknown>
  relationshipName?: string
  fromHid?: string
  fromType?: string
  toHid?: string
  toType?: string
  relationshipProperties?: Record<string, unknown>
}

export interface MutationRequest {
  actor?: Actor
  commitId?: string
  versionId?: string
  operations: MutationOperation[]
}

export interface CommitReport {
  commitId: string
  nodesChanged: string[]
  relationshipsChanged: string[]
  messagesGenerated: number
  recipientsNotified: string[]
}

export interface MessageSummary {
  messageId: string
  subject: string
  sentAt: string
  sender: string
  senderEmail: string
  messageType: string
  isRead: boolean
  relatedNodeHids: string[]
  relatedRelationshipTypes: string[]
}

export interface ReferenceItem {
  uuid: string
  frameworkName: string
  frameworkVersion: string
  externalId: string
  externalType: string
  name: string
  shortDescription: string
  longDescription?: string
  sourceUri: string
  properties: Record<string, unknown>
}

export interface ReferenceAssignmentRequest {
  actor?: Actor
  sourceHid: string
  referenceUuid?: string
  externalId?: string
  frameworkName?: string
  frameworkVersion?: string
  commitId?: string
}

export interface ReferenceAssignment {
  sourceHid: string
  referenceItem: ReferenceItem
}

export interface ReferenceAssignmentMutationResponse {
  assignment: ReferenceAssignment
  commitReport: CommitReport
}

export interface OnboardingRecord {
  hid: string
  uuid: string
  typeName: string
  userName: string
  userEmail: string
  created: string
  lastTouch: string
}

export interface CreateOnboardingRequest {
  userName: string
  userEmail: string
}

export class APIError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = "APIError"
    this.status = status
  }
}

export class SSTPAClient {
  private readonly baseUrl: string
  private readonly fetchImpl: typeof fetch
  private readonly actor?: Actor

  constructor(options: SSTPAClientOptions = {}) {
    this.baseUrl = trimTrailingSlash(options.baseUrl ?? "/api/v1")
    this.fetchImpl = options.fetchImpl ?? fetch
    this.actor = options.actor
  }

  getNodeByHID(hid: string) {
    return this.request<NodeDetail>(`/nodes/${encodeURIComponent(hid)}`)
  }

  searchNodes(params: { q?: string; type?: string; page?: number; limit?: number } = {}) {
    return this.request<ListResponse<NodeSummary>>(`/search${queryString(params)}`)
  }

  validateRelationship(payload: RelationshipValidationRequest) {
    return this.request<ValidationResponse>("/validate/relationship", {
      method: "POST",
      body: payload,
    })
  }

  commitMutation(payload: MutationRequest) {
    return this.request<CommitReport>("/mutations", {
      method: "POST",
      body: withActor(payload, this.actor),
    })
  }

  listMessages(params: { page?: number; limit?: number; sort?: string; direction?: "asc" | "desc" } = {}) {
    return this.request<ListResponse<MessageSummary>>(`/messages${queryString(params)}`)
  }

  unreadMessageCount() {
    return this.request<{ unreadCount: number }>("/messages/unread-count")
  }

  searchReferenceItems(
    params: {
      q?: string
      frameworkName?: string
      frameworkVersion?: string
      externalType?: string
      page?: number
      limit?: number
    } = {},
  ) {
    return this.request<ListResponse<ReferenceItem>>(`/reference/search${queryString(params)}`)
  }

  validateReferenceAssignment(payload: ReferenceAssignmentRequest) {
    return this.request<ValidationResponse>("/reference/validate-assignment", {
      method: "POST",
      body: payload,
    })
  }

  createReferenceAssignment(payload: ReferenceAssignmentRequest) {
    return this.request<ReferenceAssignmentMutationResponse>("/references/assignments", {
      method: "POST",
      body: withActor(payload, this.actor),
    })
  }

  deleteReferenceAssignment(payload: ReferenceAssignmentRequest) {
    return this.request<ReferenceAssignmentMutationResponse>("/references/assignments", {
      method: "DELETE",
      body: withActor(payload, this.actor),
    })
  }

  listUsers(params: { page?: number; limit?: number } = {}) {
    return this.request<ListResponse<OnboardingRecord>>(`/users${queryString(params)}`)
  }

  getUser(uuid: string) {
    return this.request<OnboardingRecord>(`/users/${encodeURIComponent(uuid)}`)
  }

  createUser(payload: CreateOnboardingRequest) {
    return this.request<OnboardingRecord>("/users", { method: "POST", body: payload })
  }

  listAdmins(params: { page?: number; limit?: number } = {}) {
    return this.request<ListResponse<OnboardingRecord>>(`/admins${queryString(params)}`)
  }

  getAdmin(uuid: string) {
    return this.request<OnboardingRecord>(`/admins/${encodeURIComponent(uuid)}`)
  }

  createAdmin(payload: CreateOnboardingRequest) {
    return this.request<OnboardingRecord>("/admins", { method: "POST", body: payload })
  }

  private async request<TResponse>(path: string, options: { method?: string; body?: unknown } = {}) {
    const headers = new Headers()
    if (this.actor) {
      headers.set("X-SSTPA-User", this.actor.name)
      headers.set("X-SSTPA-User-Email", this.actor.email)
      if (this.actor.admin) {
        headers.set("X-SSTPA-Admin", "true")
      }
    }
    if (options.body !== undefined) {
      headers.set("Content-Type", "application/json")
    }

    const response = await this.fetchImpl(`${this.baseUrl}${path}`, {
      method: options.method ?? "GET",
      headers,
      body: options.body === undefined ? undefined : JSON.stringify(options.body),
    })
    if (!response.ok) {
      throw new APIError(response.status, await errorMessage(response))
    }

    return (await response.json()) as TResponse
  }
}

function withActor<TPayload extends { actor?: Actor }>(payload: TPayload, actor: Actor | undefined) {
  if (payload.actor || !actor) {
    return payload
  }

  return { ...payload, actor }
}

function queryString(params: Record<string, string | number | undefined>) {
  const search = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== "") {
      search.set(key, String(value))
    }
  }

  const value = search.toString()
  return value === "" ? "" : `?${value}`
}

async function errorMessage(response: Response) {
  try {
    const payload = (await response.json()) as { error?: string }
    return payload.error ?? response.statusText
  } catch {
    return response.statusText
  }
}

function trimTrailingSlash(value: string) {
  return value.endsWith("/") ? value.slice(0, -1) : value
}
