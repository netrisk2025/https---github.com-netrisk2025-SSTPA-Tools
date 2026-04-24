// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
export type RequirementStatus = "candidate" | "approved" | "deferred" | "superseded" | "rejected"

export type ReferencePipelineStage = "raw" | "staged" | "normalized" | "graphImported"

export interface RequirementRecord {
  id: string
  sections: string[]
  status: RequirementStatus
  verificationIds: string[]
}

export interface ReferenceArtifact {
  frameworkName: string
  frameworkVersion: string
  stage: ReferencePipelineStage
  filePath: string
}

export function hasVerificationGap(requirement: RequirementRecord) {
  return requirement.status === "approved" && requirement.verificationIds.length === 0
}
