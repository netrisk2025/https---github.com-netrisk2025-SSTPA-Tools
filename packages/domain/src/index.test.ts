import { describe, expect, it } from "vitest"

import { hasVerificationGap } from "./index"

describe("hasVerificationGap", () => {
  it("flags approved requirements without verification", () => {
    expect(
      hasVerificationGap({
        id: "REQ-1",
        sections: ["1.3.2"],
        status: "approved",
        verificationIds: [],
      }),
    ).toBe(true)
  })

  it("allows candidate requirements to remain unmapped", () => {
    expect(
      hasVerificationGap({
        id: "REQ-2",
        sections: ["2.2.10.8.1"],
        status: "candidate",
        verificationIds: [],
      }),
    ).toBe(false)
  })
})
