export type ReportStatus = 'pending' | 'reviewed'
export type ReportTargetType = 'question' | 'tip'

export interface Report {
  id: string
  reporterUserId: string
  targetType: ReportTargetType
  targetId: string
  reason: string
  status: ReportStatus
  createdAt: string
  updatedAt: string
}

export type ReportResponse = Report

export interface CreateReportRequest {
  targetType: ReportTargetType
  targetId: string
  reason: string
}

export interface ReviewReportRequest {
  unpublish?: boolean
}