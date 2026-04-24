// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
package apihttp

import "net/http"

func (api api) openapiHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/yaml")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(openAPISpec))
}

const openAPISpec = `openapi: 3.1.0
info:
  title: SSTPA Tool Backend API
  version: 0.1.0
  summary: REST contracts for SSTPA graph, mutation, messaging, and reference data slices.
servers:
  - url: /api/v1
paths:
  /health:
    get:
      responses:
        "200":
          description: Backend health response.
  /nodes:
    get:
      parameters:
        - name: type
          in: query
          schema: { type: string }
        - name: page
          in: query
          schema: { type: integer, minimum: 1 }
        - name: limit
          in: query
          schema: { type: integer, minimum: 1, maximum: 200 }
      responses:
        "200": { description: Paginated SSTPA node summaries. }
  /nodes/{hid}:
    get:
      parameters:
        - name: hid
          in: path
          required: true
          schema: { type: string }
      responses:
        "200": { description: SSTPA node with all properties. }
  /nodes/uuid/{uuid}:
    get:
      parameters:
        - name: uuid
          in: path
          required: true
          schema: { type: string }
      responses:
        "200": { description: SSTPA node with all properties. }
  /nodes/{hid}/context:
    get:
      parameters:
        - name: hid
          in: path
          required: true
          schema: { type: string }
      responses:
        "200": { description: Node context, parent relationships, and hierarchy path. }
  /hierarchy:
    get:
      responses:
        "200": { description: Compact Capability/System hierarchy graph. }
  /search:
    get:
      parameters:
        - name: q
          in: query
          schema: { type: string }
        - name: type
          in: query
          schema: { type: string }
        - name: page
          in: query
          schema: { type: integer, minimum: 1 }
        - name: limit
          in: query
          schema: { type: integer, minimum: 1, maximum: 200 }
      responses:
        "200": { description: Paginated SSTPA node search results. }
  /validate/relationship:
    post:
      responses:
        "200": { description: Validity and reason for a proposed relationship. }
  /mutations:
    post:
      responses:
        "201": { description: Transactional mutation CommitReport. }
  /messages:
    get:
      responses:
        "200": { description: Paginated mailbox messages. }
    post:
      responses:
        "201": { description: Created direct message. }
  /messages/unread-count:
    get:
      responses:
        "200": { description: Unread message count for the current user. }
  /messages/{messageId}:
    get:
      responses:
        "200": { description: Message detail. }
    delete:
      responses:
        "200": { description: Soft-deleted message detail. }
  /messages/{messageId}/reply:
    post:
      responses:
        "201": { description: Reply message. }
  /messages/{messageId}/read:
    post:
      responses:
        "200": { description: Message marked read. }
  /reference/frameworks:
    get:
      responses:
        "200": { description: Available reference frameworks. }
  /reference/items:
    get:
      responses:
        "200": { description: Paginated reference items by framework and type. }
  /reference/items/{externalID}:
    get:
      responses:
        "200": { description: Reference item by ExternalID. }
  /reference/items/uuid/{uuid}:
    get:
      responses:
        "200": { description: Reference item by uuid. }
  /reference/items/{uuid}/related:
    get:
      responses:
        "200": { description: Related reference items. }
  /reference/search:
    get:
      responses:
        "200": { description: Paginated reference item search. }
  /reference/validate-assignment:
    post:
      responses:
        "200": { description: Validity and reason for a proposed REFERENCES assignment. }
  /references/assignments/{sourceHID}:
    get:
      responses:
        "200": { description: Reference assignments for a selected SSTPA node. }
  /references/assignments:
    post:
      responses:
        "201": { description: Transactional reference assignment mutation. }
    delete:
      responses:
        "200": { description: Transactional reference assignment removal. }
  /users:
    get:
      parameters:
        - name: page
          in: query
          schema: { type: integer, minimum: 1 }
        - name: limit
          in: query
          schema: { type: integer, minimum: 1, maximum: 200 }
      responses:
        "200": { description: Paginated registered users. }
    post:
      responses:
        "201": { description: Newly registered user. }
        "409": { description: User email already registered. }
  /users/{uuid}:
    get:
      parameters:
        - name: uuid
          in: path
          required: true
          schema: { type: string }
      responses:
        "200": { description: Registered user. }
        "404": { description: User not found. }
  /admins:
    get:
      parameters:
        - name: page
          in: query
          schema: { type: integer, minimum: 1 }
        - name: limit
          in: query
          schema: { type: integer, minimum: 1, maximum: 200 }
      responses:
        "200": { description: Paginated registered admins. }
    post:
      responses:
        "201": { description: Newly registered admin. }
        "409": { description: Admin email already registered. }
  /admins/{uuid}:
    get:
      parameters:
        - name: uuid
          in: path
          required: true
          schema: { type: string }
      responses:
        "200": { description: Registered admin. }
        "404": { description: Admin not found. }
components:
  schemas:
    CommitReport:
      type: object
      required: [commitId, nodesChanged, relationshipsChanged, messagesGenerated, recipientsNotified]
      properties:
        commitId: { type: string }
        nodesChanged: { type: array, items: { type: string } }
        relationshipsChanged: { type: array, items: { type: string } }
        messagesGenerated: { type: integer }
        recipientsNotified: { type: array, items: { type: string } }
`
