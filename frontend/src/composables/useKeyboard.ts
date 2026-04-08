import { ref, reactive, readonly, onMounted, onUnmounted, type Ref } from 'vue'
import { GetShortcuts, SetShortcuts } from '../wailsjs/go/main/App'

// ============================================================
// Types
// ============================================================

export interface ShortcutBinding {
  id: string
  label: string
  defaultKey: string
  customKey?: string
  handler: () => void
  scope: 'global' | 'editor' | 'chat'
  enabled: boolean
}

export interface ShortcutDefinition {
  id: string
  label: string
  defaultKey: string
  scope: 'global' | 'editor' | 'chat'
}

export type ShortcutOverrides = Record<string, string>

// ============================================================
// Default shortcut definitions
// ============================================================

export const DEFAULT_SHORTCUTS: ShortcutDefinition[] = [
  // AI operations
  { id: 'ai.polish', label: 'AI润色', defaultKey: 'Ctrl+Shift+R', scope: 'editor' },
  { id: 'ai.translate', label: 'AI翻译', defaultKey: 'Ctrl+Shift+T', scope: 'editor' },
  { id: 'ai.shorten', label: 'AI缩写', defaultKey: 'Ctrl+Shift+D', scope: 'editor' },

  // View
  { id: 'view.togglePreview', label: '切换预览', defaultKey: 'Ctrl+Shift+V', scope: 'global' },
  { id: 'view.toggleChat', label: 'AI对话', defaultKey: 'Ctrl+Shift+A', scope: 'global' },

  // Save
  { id: 'file.save', label: '保存', defaultKey: 'Ctrl+S', scope: 'global' },
]

// ============================================================
// Key format utilities
// ============================================================

export function normalizeKey(key: string): string {
  return key
    .replace(/Command/g, 'Ctrl')
    .replace(/Cmd/g, 'Ctrl')
    .replace(/Control/g, 'Ctrl')
}

export function parseKeyCombo(keyStr: string): { ctrl: boolean; shift: boolean; alt: boolean; key: string } {
  const parts = normalizeKey(keyStr).split('+')
  return {
    ctrl: parts.includes('Ctrl'),
    shift: parts.includes('Shift'),
    alt: parts.includes('Alt'),
    key: parts[parts.length - 1] || '',
  }
}

export function formatKeyCombo(e: KeyboardEvent): string {
  const parts: string[] = []
  if (e.ctrlKey || e.metaKey) parts.push('Ctrl')
  if (e.shiftKey) parts.push('Shift')
  if (e.altKey) parts.push('Alt')
  parts.push(e.key.length === 1 ? e.key.toUpperCase() : e.key)
  return parts.join('+')
}

export function formatKeyDisplay(keyStr: string): string {
  return normalizeKey(keyStr)
    .replace(/\+/g, ' + ')
}

export function eventMatchesKey(e: KeyboardEvent, keyStr: string): boolean {
  const combo = parseKeyCombo(keyStr)
  const ctrlMatch = combo.ctrl ? (e.ctrlKey || e.metaKey) : (!e.ctrlKey && !e.metaKey)
  const shiftMatch = combo.shift ? e.shiftKey : !e.shiftKey
  const altMatch = combo.alt ? e.altKey : !e.altKey
  const keyMatch = e.key.length === 1
    ? e.key.toUpperCase() === combo.key.toUpperCase()
    : e.key === combo.key
  return ctrlMatch && shiftMatch && altMatch && keyMatch
}

// ============================================================
// Global state (singleton)
// ============================================================

const bindings = reactive<Map<string, ShortcutBinding>>(new Map())
const overrides = ref<ShortcutOverrides>({})
const isListening = ref(false)
const isMonacoFocused = ref(false)
const currentScope = ref<'global' | 'editor' | 'chat'>('global')

let keydownHandler: ((e: KeyboardEvent) => void) | null = null

// ============================================================
// Core keyboard handler
// ============================================================

function handleGlobalKeydown(e: KeyboardEvent) {
  // Skip if target is an input/textarea (but not Monaco's special areas)
  const target = e.target as HTMLElement
  if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.tagName === 'SELECT') {
    return
  }

  // Determine if Monaco editor is focused
  const monacoFocused = target.closest('.monaco-editor') !== null || target.closest('.overflow-guard') !== null

  for (const binding of bindings.values()) {
    if (!binding.enabled) continue

    const activeKey = binding.customKey || binding.defaultKey

    // Scope filtering
    if (binding.scope === 'editor' && !monacoFocused) continue
    if (binding.scope === 'chat' && currentScope.value !== 'chat') continue

    if (eventMatchesKey(e, activeKey)) {
      // For global shortcuts, allow in Monaco unless it's Ctrl+S (which Monaco doesn't use)
      // For editor-scoped shortcuts, always intercept
      e.preventDefault()
      e.stopPropagation()
      binding.handler()
      return
    }
  }
}

// ============================================================
// Persistence
// ============================================================

async function loadOverrides(): Promise<ShortcutOverrides> {
  try {
    const json = await GetShortcuts()
    if (json) {
      const parsed = JSON.parse(json)
      if (typeof parsed === 'object' && parsed !== null) {
        return parsed as ShortcutOverrides
      }
    }
  } catch {
    // Settings not found or parse error — use defaults
  }
  return {}
}

async function saveOverrides(newOverrides: ShortcutOverrides): Promise<void> {
  try {
    const json = JSON.stringify(newOverrides)
    await SetShortcuts(json)
  } catch (err) {
    console.error('Failed to save shortcuts:', err)
  }
}

// ============================================================
// Composable
// ============================================================

export function useKeyboard() {
  let started = false

  /**
   * Start the global keyboard listener and load saved overrides.
   */
  async function start() {
    if (isListening.value) return

    // Load persisted overrides
    overrides.value = await loadOverrides()

    // Apply overrides to registered bindings
    for (const [id, customKey] of Object.entries(overrides.value)) {
      const binding = bindings.get(id)
      if (binding) {
        binding.customKey = customKey
      }
    }

    // Attach global keydown listener
    keydownHandler = handleGlobalKeydown
    window.addEventListener('keydown', keydownHandler, true)
    isListening.value = true
    started = true
  }

  /**
   * Stop the global keyboard listener.
   */
  function stop() {
    if (keydownHandler) {
      window.removeEventListener('keydown', keydownHandler, true)
      keydownHandler = null
    }
    isListening.value = false
    started = false
  }

  /**
   * Register a shortcut binding with its handler.
   */
  function register(id: string, handler: () => void) {
    const def = DEFAULT_SHORTCUTS.find(s => s.id === id)
    if (!def) {
      console.warn(`[useKeyboard] Unknown shortcut id: ${id}`)
      return
    }

    const binding: ShortcutBinding = {
      id: def.id,
      label: def.label,
      defaultKey: def.defaultKey,
      customKey: overrides.value[id],
      handler,
      scope: def.scope,
      enabled: true,
    }

    bindings.set(id, binding)
  }

  /**
   * Unregister a shortcut binding.
   */
  function unregister(id: string) {
    bindings.delete(id)
  }

  /**
   * Update a shortcut's key binding (user customization).
   */
  async function setCustomKey(id: string, newKey: string) {
    const binding = bindings.get(id)
    if (!binding) return

    // Validate: check for conflicts
    const conflict = findConflict(id, newKey)
    if (conflict) {
      console.warn(`[useKeyboard] Shortcut conflict: ${newKey} already used by ${conflict.label}`)
    }

    binding.customKey = newKey
    overrides.value = { ...overrides.value, [id]: newKey }
    await saveOverrides(overrides.value)
  }

  /**
   * Reset a shortcut to its default key binding.
   */
  async function resetToDefault(id: string) {
    const binding = bindings.get(id)
    if (!binding) return

    delete binding.customKey
    const newOverrides = { ...overrides.value }
    delete newOverrides[id]
    overrides.value = newOverrides
    await saveOverrides(overrides.value)
  }

  /**
   * Reset all shortcuts to defaults.
   */
  async function resetAll() {
    for (const binding of bindings.values()) {
      delete binding.customKey
    }
    overrides.value = {}
    await saveOverrides({})
  }

  /**
   * Find if a key is already used by another binding.
   */
  function findConflict(excludeId: string, key: string): ShortcutBinding | null {
    for (const binding of bindings.values()) {
      if (binding.id === excludeId) continue
      const activeKey = binding.customKey || binding.defaultKey
      if (normalizeKey(activeKey) === normalizeKey(key)) {
        return binding
      }
    }
    return null
  }

  /**
   * Get the effective key for a shortcut (custom or default).
   */
  function getEffectiveKey(id: string): string {
    const binding = bindings.get(id)
    if (!binding) return ''
    return binding.customKey || binding.defaultKey
  }

  /**
   * Get all registered bindings as an array.
   */
  function getAllBindings(): ShortcutBinding[] {
    return Array.from(bindings.values())
  }

  /**
   * Check if a shortcut has a custom override.
   */
  function isCustom(id: string): boolean {
    return !!bindings.get(id)?.customKey
  }

  /**
   * Set the Monaco focus state (called from editor component).
   */
  function setMonacoFocused(focused: boolean) {
    isMonacoFocused.value = focused
  }

  /**
   * Set the current scope context.
   */
  function setScope(scope: 'global' | 'editor' | 'chat') {
    currentScope.value = scope
  }

  // Auto-cleanup
  onMounted(() => {
    start()
  })

  onUnmounted(() => {
    if (started) {
      stop()
    }
  })

  return {
    // State
    isListening: readonly(isListening),
    isMonacoFocused: readonly(isMonacoFocused),
    overrides: readonly(overrides),

    // Methods
    start,
    stop,
    register,
    unregister,
    setCustomKey,
    resetToDefault,
    resetAll,
    findConflict,
    getEffectiveKey,
    getAllBindings,
    isCustom,
    setMonacoFocused,
    setScope,

    // Constants
    defaultShortcuts: DEFAULT_SHORTCUTS,
  }
}
