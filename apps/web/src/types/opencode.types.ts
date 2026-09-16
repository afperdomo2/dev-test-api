export interface OpenCodeUsageWindow {
  status: string
  percent: number
  resetsAt: string
}

export interface OpenCodeUsage {
  rolling: OpenCodeUsageWindow
  weekly: OpenCodeUsageWindow
  monthly: OpenCodeUsageWindow
}
