import { ref } from 'vue'
import {
  ShowSaveBackupDialog,
  ShowOpenBackupDialog,
  ExportBackupToPath,
  RestoreFromBackup,
  ListBackups,
  GetAutoBackupSettings,
  SetAutoBackupSettings,
} from '../wailsjs/wailsjs/go/main/App'

export interface BackupInfo {
  filename: string
  path: string
  createdAt: string
  encrypted: boolean
  dbSize: number
  tables: Record<string, number>
  checksum: string
  version: string
  app: string
  compression: string
}

export interface BackupSettings {
  enabled: string
  interval: string
}

export function useBackup() {
  const isExporting = ref(false)
  const isImporting = ref(false)
  const isLoading = ref(false)
  const backupList = ref<BackupInfo[]>([])
  const error = ref<string>('')
  const success = ref<string>('')

  function clearMessages() {
    error.value = ''
    success.value = ''
  }

  // 导出备份
  async function exportBackup(password: string = ''): Promise<string | null> {
    isExporting.value = true
    clearMessages()
    try {
      // 1. 选择保存路径
      const savePath = await ShowSaveBackupDialog()
      if (!savePath) return null

      // 2. 创建并保存备份
      const result = await ExportBackupToPath(savePath, password)
      success.value = `备份已导出到: ${result}`
      await loadBackupList()
      return result
    } catch (e: any) {
      error.value = e?.message || '导出失败'
      throw e
    } finally {
      isExporting.value = false
    }
  }

  // 导入恢复
  async function importBackup(filePath?: string, password: string = ''): Promise<boolean> {
    isImporting.value = true
    clearMessages()
    try {
      let selectedPath = filePath

      // 1. 选择备份文件（如果未提供路径）
      if (!selectedPath) {
        selectedPath = await ShowOpenBackupDialog()
        if (!selectedPath) return false
      }

      // 2. 恢复备份
      await RestoreFromBackup(selectedPath, password)

      success.value = '恢复成功，页面将重新加载...'

      // 3. 恢复成功后刷新页面
      setTimeout(() => {
        window.location.reload()
      }, 1500)

      return true
    } catch (e: any) {
      error.value = e?.message || '恢复失败'
      throw e
    } finally {
      isImporting.value = false
    }
  }

  // 加载备份列表
  async function loadBackupList(): Promise<void> {
    isLoading.value = true
    clearMessages()
    try {
      const json = await ListBackups()
      if (json) {
        backupList.value = JSON.parse(json)
      } else {
        backupList.value = []
      }
    } catch (e: any) {
      error.value = e?.message || '加载备份列表失败'
    } finally {
      isLoading.value = false
    }
  }

  // 获取自动备份设置
  async function getAutoBackupSettings(): Promise<BackupSettings> {
    try {
      const json = await GetAutoBackupSettings()
      if (json) {
        return JSON.parse(json)
      }
    } catch (e: any) {
      // 忽略错误，返回默认值
    }
    return { enabled: 'false', interval: 'daily' }
  }

  // 设置自动备份
  async function setAutoBackupSettings(enabled: boolean, interval: string): Promise<void> {
    try {
      await SetAutoBackupSettings(enabled, interval)
      success.value = '自动备份设置已保存'
    } catch (e: any) {
      error.value = e?.message || '保存设置失败'
      throw e
    }
  }

  // 格式化文件大小
  function formatSize(bytes: number): string {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
  }

  // 格式化时间
  function formatTime(isoString: string): string {
    try {
      const date = new Date(isoString)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      })
    } catch {
      return isoString
    }
  }

  return {
    isExporting,
    isImporting,
    isLoading,
    backupList,
    error,
    success,
    exportBackup,
    importBackup,
    loadBackupList,
    getAutoBackupSettings,
    setAutoBackupSettings,
    formatSize,
    formatTime,
    clearMessages,
  }
}
