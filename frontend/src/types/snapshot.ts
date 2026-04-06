export interface Snapshot {
  id: string
  resumeId: string
  label: string
  note: string
  triggerType: 'manual' | 'auto_pdf_export' | 'rollback'
  jsonData: string
  markdownContent: string
  templateId: string
  customCss: string
  createdAt: string
}

export interface SnapshotListItem {
  id: string
  resumeId: string
  label: string
  note: string
  triggerType: string
  createdAt: string
}

export interface DiffResult {
  snapshot1Id: string
  snapshot2Id: string
  snapshot1Label: string
  snapshot2Label: string
  contentDiff: string
  stats: {
    addedLines: number
    removedLines: number
    changedLines: number
  }
}
