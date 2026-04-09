import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import * as monaco from 'monaco-editor'
import { GetTheme, SetTheme } from '../wailsjs/wailsjs/go/main/App'
import { EventsOn } from '../wailsjs/wailsjs/runtime/runtime'

export type ThemeMode = 'light' | 'dark' | 'system'

const THEME_STORAGE_KEY = 'darvin-theme'

// 全局共享状态（单例模式）
const mode = ref<ThemeMode>('system')
const isDark = ref(false)
let mediaQuery: MediaQueryList | null = null
let mediaHandler: ((e: MediaQueryListEvent) => void) | null = null

function applyTheme(dark: boolean) {
  isDark.value = dark
  document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light')

  // 同步 Monaco 编辑器主题
  try {
    monaco.editor.setTheme(dark ? 'vs-dark' : 'vs')
  } catch {
    // Monaco 可能还没初始化
  }
}

function resolveIsDark(m: ThemeMode): boolean {
  if (m === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  return m === 'dark'
}

function startMediaListener() {
  if (mediaQuery) return
  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaHandler = (e: MediaQueryListEvent) => {
    if (mode.value === 'system') {
      applyTheme(e.matches)
    }
  }
  mediaQuery.addEventListener('change', mediaHandler)
}

function stopMediaListener() {
  if (mediaQuery && mediaHandler) {
    mediaQuery.removeEventListener('change', mediaHandler)
    mediaQuery = null
    mediaHandler = null
  }
}

export function useTheme() {
  const currentMode = computed(() => mode.value)

  async function setTheme(newMode: ThemeMode) {
    mode.value = newMode
    const dark = resolveIsDark(newMode)
    applyTheme(dark)

    // 持久化
    localStorage.setItem(THEME_STORAGE_KEY, newMode)
    try {
      await SetTheme(newMode)
    } catch {
      // 后端存储失败不影响前端体验
    }
  }

  async function initTheme() {
    // 优先从 localStorage 读取（快速响应），再与后端同步
    let savedMode: ThemeMode | null = null

    const local = localStorage.getItem(THEME_STORAGE_KEY)
    if (local === 'light' || local === 'dark' || local === 'system') {
      savedMode = local
    }

    // 从后端同步
    try {
      const backend = await GetTheme()
      if (backend === 'light' || backend === 'dark' || backend === 'system') {
        savedMode = backend
      }
    } catch {
      // 后端不可用，使用本地值
    }

    if (savedMode) {
      mode.value = savedMode
    }

    applyTheme(resolveIsDark(mode.value))
    startMediaListener()

    // 监听后端主题变更事件
    try {
      EventsOn('theme-changed', (theme: string) => {
        if (theme === 'light' || theme === 'dark' || theme === 'system') {
          mode.value = theme
          applyTheme(resolveIsDark(theme))
          localStorage.setItem(THEME_STORAGE_KEY, theme)
        }
      })
    } catch {
      // Wails runtime 可能不可用
    }
  }

  onMounted(() => {
    initTheme()
  })

  onUnmounted(() => {
    stopMediaListener()
  })

  return {
    mode: currentMode,
    isDark,
    setTheme,
    initTheme,
  }
}
