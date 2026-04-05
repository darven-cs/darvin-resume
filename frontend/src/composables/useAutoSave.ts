import { ref } from 'vue'
import type { Ref } from 'vue'
import { UpdateResume } from '../wailsjs/wailsjs/go/main/App'

/**
 * 自动保存 Composable
 * 管理编辑器内容的自动保存状态和逻辑
 *
 * 保存状态:
 * - saved: 内容已保存
 * - saving: 保存中
 * - unsaved: 有未保存的变更
 * - error: 保存失败
 *
 * 触发条件 per D-26:
 * - 内容变更后 30 秒间隔（定时器自动触发）
 * - 外部手动调用 triggerSave()（AI 操作完成、页面切换）
 */
export function useAutoSave(options: {
  resumeId: Ref<string>
  getData: () => { markdownContent: string; jobTarget: string }
}) {
  const saveStatus: Ref<'saved' | 'saving' | 'unsaved' | 'error'> = ref<'saved' | 'saving' | 'unsaved' | 'error'>('saved')
  const errorMessage: Ref<string> = ref('')

  let isDirty = false
  let autoSaveTimer: ReturnType<typeof setInterval> | null = null

  /**
   * 标记内容已修改
   * 调用后 isDirty = true，等待自动保存定时器触发
   */
  function markDirty(): void {
    isDirty = true
    if (saveStatus.value === 'saved') {
      saveStatus.value = 'unsaved'
    }
  }

  /**
   * 立即执行保存操作
   * 调用 UpdateResume 将数据写入后端
   */
  async function triggerSave(): Promise<void> {
    const id = options.resumeId.value
    if (!id) return

    saveStatus.value = 'saving'
    errorMessage.value = ''

    try {
      const data = options.getData()
      await UpdateResume(id, JSON.stringify(data))
      saveStatus.value = 'saved'
      isDirty = false
      errorMessage.value = ''
    } catch (err) {
      saveStatus.value = 'error'
      errorMessage.value = String(err)
      console.error('自动保存失败:', err)
    }
  }

  /**
   * 启动 30 秒自动保存定时器
   * 如果有未保存的变更，定时器会触发保存
   */
  function startAutoSave(): void {
    stopAutoSave()
    autoSaveTimer = setInterval(() => {
      if (isDirty) {
        triggerSave()
      }
    }, 30_000)
  }

  /**
   * 停止自动保存定时器
   * 离开编辑器页面时调用
   */
  function stopAutoSave(): void {
    if (autoSaveTimer !== null) {
      clearInterval(autoSaveTimer)
      autoSaveTimer = null
    }
  }

  return {
    saveStatus,
    errorMessage,
    markDirty,
    triggerSave,
    startAutoSave,
    stopAutoSave
  }
}
