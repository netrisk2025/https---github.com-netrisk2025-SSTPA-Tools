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

  it("creates a user and lists registered users", async () => {
    const requests: Request[] = []
    const client = new SSTPAClient({
      baseUrl: "http://localhost:8080/api/v1/",
      fetchImpl: async (input, init) => {
        requests.push(new Request(input, init))
        if (init?.method === "POST") {
          return Response.json(
            {
              hid: "USR__1",
              uuid: "u-1",
              typeName: "User",
              userName: "Alice",
              userEmail: "alice@example.test",
              created: "2026-04-24T12:00:00Z",
              lastTouch: "2026-04-24T12:00:00Z",
            },
            { status: 201 },
          )
        }
        return Response.json({
          items: [
            {
              hid: "USR__1",
              uuid: "u-1",
              typeName: "User",
              userName: "Alice",
              userEmail: "alice@example.test",
              created: "2026-04-24T12:00:00Z",
              lastTouch: "2026-04-24T12:00:00Z",
            },
          ],
          page: 1,
          limit: 50,
          total: 1,
        })
      },
    })

    const created = await client.createUser({ userName: "Alice", userEmail: "alice@example.test" })
    expect(created.hid).toBe("USR__1")
    expect(requests[0].method).toBe("POST")
    expect(requests[0].url).toBe("http://localhost:8080/api/v1/users")
    await expect(requests[0].json()).resolves.toEqual({
      userName: "Alice",
      userEmail: "alice@example.test",
    })

    const list = await client.listUsers()
    expect(list.total).toBe(1)
    expect(list.items[0].userEmail).toBe("alice@example.test")
  })

  it("bootstraps the installer admin and user", async () => {
    const requests: Request[] = []
    const client = new SSTPAClient({
      baseUrl: "http://localhost:8080/api/v1/",
      fetchImpl: async (input, init) => {
        requests.push(new Request(input, init))
        return Response.json(
          {
            admin: {
              hid: "ADM__1",
              uuid: "a-1",
              typeName: "Admin",
              userName: "Installer",
              userEmail: "installer@example.test",
              created: "2026-04-24T12:00:00Z",
              lastTouch: "2026-04-24T12:00:00Z",
            },
            user: {
              hid: "USR__1",
              uuid: "u-1",
              typeName: "User",
              userName: "Installer",
              userEmail: "installer@example.test",
              created: "2026-04-24T12:00:00Z",
              lastTouch: "2026-04-24T12:00:00Z",
            },
          },
          { status: 201 },
        )
      },
    })

    const created = await client.bootstrapInstaller({
      installerName: "Installer",
      installerEmail: "installer@example.test",
    })
    expect(created.admin.hid).toBe("ADM__1")
    expect(created.user.hid).toBe("USR__1")
    expect(requests[0].url).toBe("http://localhost:8080/api/v1/onboarding/bootstrap")
    await expect(requests[0].json()).resolves.toEqual({
      installerName: "Installer",
      installerEmail: "installer@example.test",
    })
  })

  it("creates an admin with an admin actor", async () => {
    const requests: Request[] = []
    const client = new SSTPAClient({
      baseUrl: "http://localhost:8080/api/v1/",
      actor: { name: "Installer", email: "installer@example.test", admin: true },
      fetchImpl: async (input, init) => {
        requests.push(new Request(input, init))
        return Response.json(
          {
            hid: "ADM__2",
            uuid: "a-2",
            typeName: "Admin",
            userName: "Root",
            userEmail: "root@example.test",
            created: "2026-04-24T12:00:00Z",
            lastTouch: "2026-04-24T12:00:00Z",
          },
          { status: 201 },
        )
      },
    })

    const created = await client.createAdmin({ userName: "Root", userEmail: "root@example.test" })
    expect(created.hid).toBe("ADM__2")
    expect(requests[0].url).toBe("http://localhost:8080/api/v1/admins")
    expect(requests[0].headers.get("X-SSTPA-User")).toBe("Installer")
    expect(requests[0].headers.get("X-SSTPA-User-Email")).toBe("installer@example.test")
    expect(requests[0].headers.get("X-SSTPA-Admin")).toBe("true")
  })
})
