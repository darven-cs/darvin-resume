import { ref } from 'vue'
import type { Snapshot, SnapshotListItem, DiffResult } from '@/types/snapshot'
import {
  ListSnapshots,
  CreateSnapshot,
  DiffSnapshots,
  RollbackToSnapshot,
  DeleteSnapshot,
  GetSnapshot,
  GetSnapshotMarkdown,
} from '@/wailsjs/go/main/App'

export function useSnapshot() {
  const snapshots = ref<SnapshotListItem[]>([])
  const isLoading = ref(false)
  const diffResult = ref<DiffResult | null>(null)

  // 加载快照列表
  async function loadSnapshots(resumeId: string) {
    isLoading.value = true
    try {
      const list = await ListSnapshots(resumeId)
      snapshots.value = list || []
    } finally {
      isLoading.value = false
    }
  }

  // 创建快照
  async function createSnapshot(
    resumeId: string,
    label: string,
    note: string = '',
    triggerType: 'manual' | 'auto_pdf_export' | 'rollback' = 'manual'
  ): Promise<Snapshot> {
    const snap = await CreateSnapshot(resumeId, label, note, triggerType)
    // 刷新列表
    await loadSnapshots(resumeId)
    return snap as Snapshot
  }

  // 自动创建快照（PDF导出后）
  async function autoCreateSnapshot(
    resumeId: string,
    label: string = 'PDF 导出快照',
    note: string = ''
  ): Promise<Snapshot> {
    return createSnapshot(resumeId, label, note, 'auto_pdf_export')
  }

  // 对比两个快照
  async function diffSnapshots(id1: string, id2: string): Promise<DiffResult> {
    const result = await DiffSnapshots(id1, id2) as DiffResult
    diffResult.value = result
    return result
  }

  // 回滚到指定快照
  async function rollback(resumeId: string, snapshotId: string): Promise<Snapshot | null> {
    const snap = await RollbackToSnapshot(resumeId, snapshotId)
    // 刷新列表
    await loadSnapshots(resumeId)
    return snap as Snapshot | null
  }

  // 删除快照
  async function deleteSnapshot(snapshotId: string, resumeId: string) {
    await DeleteSnapshot(snapshotId)
    // 刷新列表
    await loadSnapshots(resumeId)
  }

  // 获取快照详情（用于回滚时预览）
  async function getSnapshot(snapshotId: string): Promise<Snapshot> {
    return await GetSnapshot(snapshotId) as Snapshot
  }

  // 获取快照 Markdown 内容（用于编辑器加载）
  async function getSnapshotMarkdown(snapshotId: string): Promise<{
    markdown: string
    templateId: string
    customCss: string
    jsonData: string
  }> {
    const result = await GetSnapshotMarkdown(snapshotId)
    return {
      markdown: result[0],
      templateId: result[1],
      customCss: result[2],
      jsonData: result[3],
    }
  }

  return {
    snapshots,
    isLoading,
    diffResult,
    loadSnapshots,
    createSnapshot,
    autoCreateSnapshot,
    diffSnapshots,
    rollback,
    deleteSnapshot,
    getSnapshot,
    getSnapshotMarkdown,
  }
}
