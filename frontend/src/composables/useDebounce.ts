import { ref } from 'vue'

/**
 * 防抖 Hook
 * @param fn 要防抖的函数
 * @param delay 延迟时间（毫秒），默认 150ms per D-10
 */
export function useDebounce<T extends (...args: any[]) => any>(
  fn: T,
  delay: number = 150
) {
  const timer = ref<ReturnType<typeof setTimeout> | null>(null)
  const pending = ref(false)

  const debouncedFn = (...args: Parameters<T>) => {
    if (timer.value) {
      clearTimeout(timer.value)
    }
    pending.value = true
    timer.value = setTimeout(() => {
      fn(...args)
      pending.value = false
    }, delay)
  }

  const cancel = () => {
    if (timer.value) {
      clearTimeout(timer.value)
      timer.value = null
      pending.value = false
    }
  }

  const flush = (...args: Parameters<T>) => {
    cancel()
    fn(...args)
  }

  return { debouncedFn, cancel, flush, pending }
}

/**
 * 独立的防抖工具函数（用于 Ref 值监听）
 */
export function createDebounce<T>(fn: (val: T) => void, delay: number = 150) {
  let timer: ReturnType<typeof setTimeout> | null = null

  return (value: T) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      fn(value)
      timer = null
    }, delay)
  }
}
