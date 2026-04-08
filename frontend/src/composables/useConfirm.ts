/**
 * 通用确认对话框 composable
 *
 * 提供编程式确认对话框，替代 window.confirm()。
 * 使用全局单例模式，应用中任何组件均可调用。
 *
 * 用法:
 * const { confirm } = useConfirm()
 * const ok = await confirm({
 *   title: '删除简历',
 *   message: '该简历将移入回收站，30天内可恢复',
 *   confirmLabel: '删除',
 *   type: 'danger'
 * })
 */

import { ref } from 'vue'

export interface ConfirmOptions {
  title: string
  message: string
  confirmLabel?: string
  cancelLabel?: string
  type?: 'default' | 'danger'
}

// 全局确认对话框状态（单例）
const dialogState = ref<{
  visible: boolean
  title: string
  message: string
  confirmLabel: string
  cancelLabel: string
  type: 'default' | 'danger'
}>({
  visible: false,
  title: '',
  message: '',
  confirmLabel: '确认',
  cancelLabel: '取消',
  type: 'default',
})

let resolvePromise: ((value: boolean) => void) | null = null

function show(options: ConfirmOptions): Promise<boolean> {
  return new Promise((resolve) => {
    dialogState.value = {
      visible: true,
      title: options.title,
      message: options.message,
      confirmLabel: options.confirmLabel ?? '确认',
      cancelLabel: options.cancelLabel ?? '取消',
      type: options.type ?? 'default',
    }
    resolvePromise = resolve
  })
}

function handleConfirm() {
  if (resolvePromise) {
    resolvePromise(true)
    resolvePromise = null
  }
  dialogState.value.visible = false
}

function handleCancel() {
  if (resolvePromise) {
    resolvePromise(false)
    resolvePromise = null
  }
  dialogState.value.visible = false
}

export function useConfirm() {
  return {
    /** 对话框状态（用于模板绑定） */
    dialogState,
    /** 显示确认对话框，返回 Promise<boolean> */
    confirm: show,
    /** 内部：确认按钮回调 */
    onConfirm: handleConfirm,
    /** 内部：取消按钮回调 */
    onCancel: handleCancel,
  }
}
