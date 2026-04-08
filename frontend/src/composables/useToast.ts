/**
 * 通用 Toast 通知系统 composable
 *
 * 基于 provide/inject 的全局单例模式，支持 success/error/warning/info 四种类型。
 * 最大同时显示 5 条 toast，超出自动移除最早的。
 */

import { ref, readonly } from 'vue'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface ToastAction {
  label: string
  handler: () => void
}

export interface ToastItem {
  id: string
  type: ToastType
  title: string
  message?: string
  duration: number   // 0 = 不自动消失
  action?: ToastAction
  createdAt: number
  progress: number   // 0~1，用于进度条
}

const MAX_TOASTS = 5

// 全局 toast 队列（单例）
const toasts = ref<ToastItem[]>([])
let idCounter = 0

function generateId(): string {
  return `toast-${++idCounter}-${Date.now()}`
}

function addToast(item: Omit<ToastItem, 'id' | 'createdAt' | 'progress'>): string {
  const id = generateId()
  const newToast: ToastItem = {
    ...item,
    id,
    createdAt: Date.now(),
    progress: 1,
  }

  toasts.value.push(newToast)

  // 超过最大数量，移除最早的
  while (toasts.value.length > MAX_TOASTS) {
    toasts.value.shift()
  }

  // 启动自动消失计时器
  if (item.duration > 0) {
    startDismissTimer(id, item.duration)
  }

  return id
}

const dismissTimers = new Map<string, ReturnType<typeof setTimeout>>()

function startDismissTimer(id: string, duration: number) {
  const existing = dismissTimers.get(id)
  if (existing) clearTimeout(existing)

  const timer = setTimeout(() => {
    removeToast(id)
    dismissTimers.delete(id)
  }, duration)

  dismissTimers.set(id, timer)
}

function removeToast(id: string) {
  const idx = toasts.value.findIndex(t => t.id === id)
  if (idx !== -1) {
    toasts.value.splice(idx, 1)
  }
  // 清除计时器
  const timer = dismissTimers.get(id)
  if (timer) {
    clearTimeout(timer)
    dismissTimers.delete(id)
  }
}

function clearAllToasts() {
  for (const [id, timer] of dismissTimers) {
    clearTimeout(timer)
  }
  dismissTimers.clear()
  toasts.value = []
}

/**
 * useToast composable — 应用中任何组件均可调用
 *
 * 用法:
 * const toast = useToast()
 * toast.success('保存成功')
 * toast.error('导出失败', 'PDF导出过程中发生错误')
 * toast.warning('自动备份未启用', undefined, { action: { label: '去设置', handler: openSettings } })
 */
export function useToast() {
  /**
   * 显示一条 toast
   */
  function show(options: {
    type: ToastType
    title: string
    message?: string
    duration?: number
    action?: ToastAction
  }): string {
    return addToast({
      type: options.type,
      title: options.title,
      message: options.message,
      duration: options.duration ?? getDefaultDuration(options.type),
      action: options.action,
    })
  }

  /**
   * 成功提示
   */
  function success(title: string, message?: string, action?: ToastAction): string {
    return show({ type: 'success', title, message, action })
  }

  /**
   * 错误提示
   */
  function error(title: string, message?: string, action?: ToastAction): string {
    return show({ type: 'error', title, message, action, duration: getDefaultDuration('error') })
  }

  /**
   * 警告提示
   */
  function warning(title: string, message?: string, action?: ToastAction): string {
    return show({ type: 'warning', title, message, action })
  }

  /**
   * 信息提示
   */
  function info(title: string, message?: string, action?: ToastAction): string {
    return show({ type: 'info', title, message, action })
  }

  /**
   * 手动移除指定 toast
   */
  function remove(id: string) {
    removeToast(id)
  }

  /**
   * 清除所有 toast
   */
  function clear() {
    clearAllToasts()
  }

  /**
   * 更新 toast 进度（用于进度条动画）
   */
  function updateProgress(id: string, progress: number) {
    const toast = toasts.value.find(t => t.id === id)
    if (toast) {
      toast.progress = Math.max(0, Math.min(1, progress))
    }
  }

  return {
    toasts: readonly(toasts),
    show,
    success,
    error,
    warning,
    info,
    remove,
    clear,
    updateProgress,
  }
}

/**
 * 根据类型返回默认显示时长
 */
function getDefaultDuration(type: ToastType): number {
  switch (type) {
    case 'success': return 3000
    case 'error': return 0    // 错误不自动消失，需手动关闭
    case 'warning': return 5000
    case 'info': return 4000
    default: return 4000
  }
}
