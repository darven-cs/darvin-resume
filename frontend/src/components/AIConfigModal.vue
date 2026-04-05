<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="handleClose">
        <div class="modal-box" role="dialog" aria-modal="true" aria-labelledby="config-modal-title">
          <div class="modal-header">
            <h2 id="config-modal-title" class="modal-title">AI 设置</h2>
            <button class="close-btn" @click="handleClose" aria-label="关闭">
              <span>&times;</span>
            </button>
          </div>

          <div class="modal-body">
            <!-- API Key -->
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

            <!-- Base URL -->
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

            <!-- Default Model -->
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

            <!-- Max Tokens -->
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

              <!-- Timeout -->
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

            <!-- Include Full Context -->
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

            <!-- Validation Error -->
            <div v-if="validationError" class="validation-error">
              {{ validationError }}
            </div>
          </div>

          <div class="modal-footer">
            <button class="btn btn-secondary" @click="handleClose" :disabled="isSaving">
              取消
            </button>
            <button class="btn btn-secondary" @click="handleValidate" :disabled="isValidating || isSaving">
              {{ isValidating ? '验证中...' : '验证连接' }}
            </button>
            <button class="btn btn-primary" @click="handleSave" :disabled="isSaving || !isFormValid">
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
import type { AIConfig } from '../types/ai'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

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

// Local copy of config for the form
const localConfig = ref<AIConfig>({
  apiKey: '',
  baseURL: 'https://api.anthropic.com',
  defaultModel: 'claude-sonnet-4-20250514',
  maxTokens: 4096,
  timeoutSeconds: 60,
})

const includeFullContext = ref(false)

// Errors from form validation
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

// Load config when modal opens
watch(() => props.visible, async (isVisible) => {
  if (isVisible) {
    await loadConfig()
    localConfig.value = { ...config.value }
  }
})

// Sync config changes back
watch(config, (newConfig) => {
  localConfig.value = { ...newConfig }
}, { deep: true })

async function handleValidate() {
  // Update config before validating
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
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(2px);
}

.modal-box {
  background: #252526;
  border: 1px solid #3c3c3c;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
  width: 90%;
  max-width: 480px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 20px 0;
}

.modal-title {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.close-btn {
  background: transparent;
  border: none;
  color: #888;
  font-size: 20px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  line-height: 1;
}

.close-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

.modal-body {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

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
  color: #ccc;
}

.form-group input[type="text"],
.form-group input[type="password"],
.form-group input[type="number"] {
  padding: 8px 10px;
  background: #1e1e1e;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  color: #e0e0e0;
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
}

.form-group input:focus {
  border-color: #3d8bfd;
}

.form-group input.input-error {
  border-color: #e5484d;
}

.error-text {
  font-size: 11px;
  color: #e5484d;
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.toggle-checkbox {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.toggle-text {
  font-size: 13px;
  color: #e0e0e0;
  font-weight: 500;
}

.toggle-hint {
  font-size: 11px;
  color: #888;
  margin-left: auto;
}

.validation-error {
  padding: 8px 12px;
  background: #2d1f1f;
  border: 1px solid #e5484d;
  border-radius: 4px;
  color: #ff7b7b;
  font-size: 12px;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding: 0 20px 20px;
}

.btn {
  padding: 6px 16px;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  background: #333;
  color: #e0e0e0;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s, border-color 0.15s, opacity 0.15s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn:hover:not(:disabled) {
  background: #3c3c3c;
  border-color: #555;
}

.btn-primary {
  background: #3d8bfd;
  border-color: #3d8bfd;
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: #5a9bff;
  border-color: #5a9bff;
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
