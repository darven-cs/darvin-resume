<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="handleClose">
        <div class="settings-box" role="dialog" aria-modal="true" aria-labelledby="settings-title">
          <div class="settings-header">
            <h2 id="settings-title" class="settings-title">设置</h2>
            <button class="close-btn" @click="handleClose" aria-label="关闭">
              <span>&times;</span>
            </button>
          </div>

          <!-- Tab 栏 -->
          <div class="settings-tabs">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              :class="['tab-btn', { active: activeTab === tab.id }]"
              @click="activeTab = tab.id"
            >
              {{ tab.label }}
            </button>
          </div>

          <!-- Tab 内容 -->
          <div class="settings-body">
            <!-- AI 配置 Tab -->
            <div v-if="activeTab === 'ai'" class="tab-panel">
              <div class="form-group">
                <label for="api-key">API Key</label>
                <input
                  id="api-key"
                  v-model="localConfig.apiKey"
                  type="password"
                  placeholder="sk-ant-..."
                  autocomplete="off"
                  :class="{ 'input-error': errors.apiKey }"
                />
                <span v-if="errors.apiKey" class="error-text">{{ errors.apiKey }}</span>
              </div>

              <div class="form-group">
                <label for="base-url">Base URL</label>
                <input
                  id="base-url"
                  v-model="localConfig.baseURL"
                  type="text"
                  placeholder="https://api.anthropic.com"
                  :class="{ 'input-error': errors.baseURL }"
                />
                <span v-if="errors.baseURL" class="error-text">{{ errors.baseURL }}</span>
              </div>

              <div class="form-group">
                <label for="default-model">默认模型</label>
                <input
                  id="default-model"
                  v-model="localConfig.defaultModel"
                  type="text"
                  placeholder="claude-sonnet-4-20250514"
                  :class="{ 'input-error': errors.defaultModel }"
                />
                <span v-if="errors.defaultModel" class="error-text">{{ errors.defaultModel }}</span>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label for="max-tokens">最大 Token 数</label>
                  <input
                    id="max-tokens"
                    v-model.number="localConfig.maxTokens"
                    type="number"
                    min="1"
                    max="8192"
                    :class="{ 'input-error': errors.maxTokens }"
                  />
                  <span v-if="errors.maxTokens" class="error-text">{{ errors.maxTokens }}</span>
                </div>

                <div class="form-group">
                  <label for="timeout">超时时间（秒）</label>
                  <input
                    id="timeout"
                    v-model.number="localConfig.timeoutSeconds"
                    type="number"
                    min="1"
                    max="300"
                    :class="{ 'input-error': errors.timeout }"
                  />
                  <span v-if="errors.timeout" class="error-text">{{ errors.timeout }}</span>
                </div>
              </div>

              <div class="form-group">
                <label class="toggle-label">
                  <input
                    v-model="includeFullContext"
                    type="checkbox"
                    class="toggle-checkbox"
                  />
                  <span class="toggle-text">包含全文参考</span>
                  <span class="toggle-hint">发送给 AI 时附上完整简历内容</span>
                </label>
              </div>

              <div v-if="validationError" class="validation-error">
                {{ validationError }}
              </div>
            </div>

            <!-- 外观 Tab -->
            <div v-if="activeTab === 'appearance'" class="tab-panel">
              <div class="form-group">
                <label>主题模式</label>
                <div class="theme-options">
                  <label
                    v-for="opt in themeOptions"
                    :key="opt.value"
                    :class="['theme-option', { active: themeMode === opt.value }]"
                  >
                    <input
                      type="radio"
                      :value="opt.value"
                      v-model="themeMode"
                      class="theme-radio"
                    />
                    <span class="theme-icon">{{ opt.icon }}</span>
                    <span class="theme-label">{{ opt.label }}</span>
                  </label>
                </div>
              </div>
              <p class="theme-status">
                当前: {{ themeMode === 'system' ? (isDark ? '深色模式（跟随系统）' : '浅色模式（跟随系统）') : themeMode === 'dark' ? '深色模式' : '浅色模式' }}
              </p>
            </div>

            <!-- 快捷键 Tab -->
            <ShortcutSettingsPanel v-if="activeTab === 'shortcuts'" />
          </div>

          <!-- 底部按钮 -->
          <div class="settings-footer">
            <button class="btn btn-secondary" @click="handleClose" :disabled="isSaving">
              关闭
            </button>
            <button
              v-if="activeTab === 'ai'"
              class="btn btn-secondary"
              @click="handleValidate"
              :disabled="isValidating || isSaving"
            >
              {{ isValidating ? '验证中...' : '验证连接' }}
            </button>
            <button
              v-if="activeTab === 'ai'"
              class="btn btn-primary"
              @click="handleSave"
              :disabled="isSaving || !isFormValid"
            >
              {{ isSaving ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useAIConfig } from '../composables/useAIConfig'
import { useTheme } from '../composables/useTheme'
import ShortcutSettingsPanel from './ShortcutSettingsPanel.vue'
import type { AIConfig } from '../types/ai'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const tabs = [
  { id: 'ai', label: 'AI 配置' },
  { id: 'appearance', label: '外观' },
  { id: 'shortcuts', label: '快捷键' },
]

const activeTab = ref('ai')

// AI Config
const {
  config,
  isSaving,
  isValidating,
  validationError,
  validationErrors,
  isValid: isFormValid,
  loadConfig,
  persistConfig,
  validateKey,
} = useAIConfig()

const localConfig = ref<AIConfig>({
  apiKey: '',
  baseURL: 'https://api.anthropic.com',
  defaultModel: 'claude-sonnet-4-20250514',
  maxTokens: 4096,
  timeoutSeconds: 60,
})

const includeFullContext = ref(false)

const errors = computed(() => {
  const errs: Record<string, string> = {}
  if (!localConfig.value.apiKey.trim()) {
    errs.apiKey = 'API Key 不能为空'
  }
  if (!localConfig.value.baseURL.trim()) {
    errs.baseURL = 'Base URL 不能为空'
  } else {
    try {
      new URL(localConfig.value.baseURL)
    } catch {
      errs.baseURL = '无效的 URL 格式'
    }
  }
  if (!localConfig.value.defaultModel.trim()) {
    errs.defaultModel = '模型名称不能为空'
  }
  if (localConfig.value.maxTokens <= 0) {
    errs.maxTokens = '必须大于 0'
  }
  if (localConfig.value.maxTokens > 8192) {
    errs.maxTokens = '建议不超过 8192'
  }
  if (localConfig.value.timeoutSeconds <= 0) {
    errs.timeout = '必须大于 0'
  }
  if (localConfig.value.timeoutSeconds > 300) {
    errs.timeout = '建议不超过 300 秒'
  }
  return errs
})

// Theme
const { mode: themeMode, isDark, setTheme } = useTheme()

const themeOptions = [
  { value: 'light', label: '浅色', icon: '☀️' },
  { value: 'dark', label: '深色', icon: '🌙' },
  { value: 'system', label: '跟随系统', icon: '💻' },
]

// Theme changes apply immediately
watch(themeMode, (newMode) => {
  setTheme(newMode as 'light' | 'dark' | 'system')
})

// Load config when modal opens
watch(() => props.visible, async (isVisible) => {
  if (isVisible) {
    await loadConfig()
    localConfig.value = { ...config.value }
  }
})

watch(config, (newConfig) => {
  localConfig.value = { ...newConfig }
}, { deep: true })

async function handleValidate() {
  Object.assign(config.value, localConfig.value)
  await validateKey()
}

async function handleSave() {
  Object.assign(config.value, localConfig.value)
  try {
    await persistConfig()
    emit('saved')
    handleClose()
  } catch {
    // Error already set in composable
  }
}

function handleClose() {
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 9998;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--ui-overlay-bg);
  backdrop-filter: blur(2px);
}

.settings-box {
  background: var(--ui-bg-secondary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-lg);
  box-shadow: var(--ui-shadow-lg);
  width: 90%;
  max-width: 520px;
  max-height: 90vh;
  overflow-y: auto;
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--ui-spacing-lg) var(--ui-spacing-lg) 0;
}

.settings-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--ui-text-primary);
  margin: 0;
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--ui-text-tertiary);
  font-size: 20px;
  cursor: pointer;
  padding: var(--ui-spacing-xs) var(--ui-spacing-sm);
  border-radius: var(--ui-radius-sm);
  line-height: 1;
}

.close-btn:hover {
  color: var(--ui-text-primary);
  background: var(--ui-bg-hover);
}

/* Tabs */
.settings-tabs {
  display: flex;
  gap: 0;
  border-bottom: 1px solid var(--ui-border);
  padding: var(--ui-spacing-md) var(--ui-spacing-lg) 0;
}

.tab-btn {
  padding: var(--ui-spacing-sm) var(--ui-spacing-md);
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  color: var(--ui-text-secondary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: color var(--ui-transition-fast), border-color var(--ui-transition-fast);
}

.tab-btn:hover {
  color: var(--ui-text-primary);
}

.tab-btn.active {
  color: var(--ui-accent);
  border-bottom-color: var(--ui-accent);
}

/* Body */
.settings-body {
  padding: var(--ui-spacing-md) var(--ui-spacing-lg);
  min-height: 200px;
}

.tab-panel {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* Form */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.form-group label {
  font-size: 12px;
  font-weight: 500;
  color: var(--ui-text-secondary);
}

.form-group input[type="text"],
.form-group input[type="password"],
.form-group input[type="number"] {
  padding: var(--ui-spacing-sm) 10px;
  background: var(--ui-input-bg);
  border: 1px solid var(--ui-input-border);
  border-radius: var(--ui-radius-sm);
  color: var(--ui-text-primary);
  font-size: 13px;
  outline: none;
  transition: border-color var(--ui-transition-fast);
}

.form-group input:focus {
  border-color: var(--ui-border-focus);
}

.form-group input.input-error {
  border-color: var(--ui-danger);
}

.error-text {
  font-size: 11px;
  color: var(--ui-danger);
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: var(--ui-spacing-sm);
  cursor: pointer;
}

.toggle-checkbox {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.toggle-text {
  font-size: 13px;
  color: var(--ui-text-primary);
  font-weight: 500;
}

.toggle-hint {
  font-size: 11px;
  color: var(--ui-text-tertiary);
  margin-left: auto;
}

.validation-error {
  padding: var(--ui-spacing-sm) var(--ui-spacing-md);
  background: var(--ui-bg-active);
  border: 1px solid var(--ui-danger);
  border-radius: var(--ui-radius-sm);
  color: var(--ui-danger);
  font-size: 12px;
}

/* Theme Options */
.theme-options {
  display: flex;
  gap: var(--ui-spacing-sm);
  margin-top: var(--ui-spacing-xs);
}

.theme-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--ui-spacing-xs);
  padding: var(--ui-spacing-md) var(--ui-spacing-lg);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-md);
  cursor: pointer;
  flex: 1;
  transition: border-color var(--ui-transition-fast), background var(--ui-transition-fast);
}

.theme-option:hover {
  border-color: var(--ui-border-hover);
}

.theme-option.active {
  border-color: var(--ui-accent);
  background: var(--ui-bg-active);
}

.theme-radio {
  display: none;
}

.theme-icon {
  font-size: 20px;
}

.theme-label {
  font-size: 12px;
  color: var(--ui-text-primary);
  font-weight: 500;
}

.theme-status {
  font-size: 12px;
  color: var(--ui-text-tertiary);
  margin: var(--ui-spacing-sm) 0 0;
}

.placeholder-text {
  color: var(--ui-text-tertiary);
  font-size: 13px;
  text-align: center;
  padding: var(--ui-spacing-xl) 0;
}

/* Footer */
.settings-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--ui-spacing-sm);
  padding: 0 var(--ui-spacing-lg) var(--ui-spacing-lg);
}

.btn {
  padding: 6px 16px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  background: var(--ui-bg-tertiary);
  color: var(--ui-text-primary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color var(--ui-transition-fast), border-color var(--ui-transition-fast), opacity var(--ui-transition-fast);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn:hover:not(:disabled) {
  background: var(--ui-bg-active);
}

.btn-primary {
  background: var(--ui-accent);
  border-color: var(--ui-accent);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: var(--ui-accent-hover);
  border-color: var(--ui-accent-hover);
}

/* Transition */
.modal-enter-active {
  animation: modal-in 0.2s ease-out;
}

.modal-leave-active {
  animation: modal-out 0.15s ease-in;
}

@keyframes modal-in {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes modal-out {
  from {
    opacity: 1;
    transform: scale(1);
  }
  to {
    opacity: 0;
    transform: scale(0.95);
  }
}
</style>
