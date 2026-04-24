// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import { describe, expect, it } from "vitest"

import { APIError, SSTPAClient } from "./index"

describe("SSTPAClient", () => {
  it("posts mutations with actor headers and body fallback", async () => {
    const requests: Request[] = []
    const client = new SSTPAClient({
      baseUrl: "http://localhost:8080/api/v1/",
      actor: { name: "Alice", email: "alice@example.test" },
      fetchImpl: async (input, init) => {
        requests.push(new Request(input, init))
        return Response.json({
          commitId: "commit-1",
          nodesChanged: ["CAP__0"],
          relationshipsChanged: [],
          messagesGenerated: 0,
          recipientsNotified: [],
        })
      },
    })

    const report = await client.commitMutation({
      operations: [{ kind: "update_node", hid: "CAP__0", properties: { Name: "Changed" } }],
    })

    expect(report.commitId).toBe("commit-1")
    expect(requests[0].url).toBe("http://localhost:8080/api/v1/mutations")
    expect(requests[0].headers.get("X-SSTPA-User-Email")).toBe("alice@example.test")
    await expect(requests[0].json()).resolves.toMatchObject({
      actor: { name: "Alice", email: "alice@example.test" },
    })
  })

  it("throws APIError with server message", async () => {
    const client = new SSTPAClient({
      fetchImpl: async () => Response.json({ error: "invalid relationship" }, { status: 400 }),
    })

    await expect(
      client.validateRelationship({ relationshipName: "BAD", fromType: "System", toType: "System" }),
    ).rejects.toEqual(new APIError(400, "invalid relationship"))
  })
})
