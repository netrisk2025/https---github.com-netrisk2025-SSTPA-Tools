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
