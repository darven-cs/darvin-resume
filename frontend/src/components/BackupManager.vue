<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="handleClose">
        <div class="modal-box" role="dialog" aria-modal="true" aria-labelledby="backup-modal-title">
          <div class="modal-header">
            <h2 id="backup-modal-title" class="modal-title">数据备份与恢复</h2>
            <button class="close-btn" @click="handleClose" aria-label="关闭">
              <span>&times;</span>
            </button>
          </div>

          <div class="modal-body">
            <!-- 导出区域 -->
            <div class="section">
              <div class="section-title">导出备份</div>
              <div class="section-desc">
                将所有简历数据导出为加密备份文件（.darvin-backup）
              </div>

              <div class="form-group">
                <label for="export-password">密码保护（可选）</label>
                <input
                  id="export-password"
                  v-model="exportPassword"
                  type="password"
                  placeholder="留空则不加密"
                  autocomplete="new-password"
                />
              </div>

              <button
                class="btn btn-primary"
                :disabled="isExporting"
                @click="handleExport"
              >
                <span v-if="isExporting" class="btn-spinner"></span>
                {{ isExporting ? '导出中...' : '导出备份' }}
              </button>
            </div>

            <div class="divider"></div>

            <!-- 导入区域 -->
            <div class="section">
              <div class="section-title">导入恢复</div>
              <div class="section-desc warning">
                恢复将替换当前所有数据，恢复前会自动创建备份
              </div>

              <div class="form-group">
                <label for="import-password">备份密码（如有）</label>
                <input
                  id="import-password"
                  v-model="importPassword"
                  type="password"
                  placeholder="输入备份密码"
                  autocomplete="new-password"
                />
              </div>

              <button
                class="btn btn-danger"
                :disabled="isImporting"
                @click="handleImport"
              >
                <span v-if="isImporting" class="btn-spinner"></span>
                {{ isImporting ? '恢复中...' : '选择备份文件并恢复' }}
              </button>
            </div>

            <div class="divider"></div>

            <!-- 自动备份设置 -->
            <div class="section">
              <div class="section-title-row">
                <span class="section-title">自动备份</span>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="autoEnabled" @change="handleAutoBackupChange" />
                  <span class="toggle-slider"></span>
                </label>
              </div>

              <div class="form-group" v-if="autoEnabled">
                <label for="auto-interval">备份周期</label>
                <select
                  id="auto-interval"
                  v-model="autoInterval"
                  @change="handleAutoBackupChange"
                  class="select-input"
                >
                  <option value="daily">每日</option>
                  <option value="weekly">每周</option>
                  <option value="monthly">每月</option>
                </select>
              </div>

              <div class="section-desc">
                {{ autoEnabled ? `自动备份已启用 — ${intervalLabel}` : '自动备份已禁用' }}
              </div>
              <div class="section-desc muted">
                自动备份使用设备密钥加密，无需手动输入密码。最多保留 10 个备份文件。
              </div>
            </div>

            <div class="divider"></div>

            <!-- 本地备份列表 -->
            <div class="section">
              <div class="section-title-row">
                <span class="section-title">本地备份</span>
                <button class="btn btn-sm" :disabled="isLoading" @click="handleRefresh">
                  {{ isLoading ? '加载中...' : '刷新' }}
                </button>
              </div>
              <div class="section-desc">
                自动保留最多 10 个备份文件，超出自动删除最旧的
              </div>

              <div v-if="backupList.length === 0" class="empty-state">
                暂无本地备份
              </div>

              <div v-else class="backup-list">
                <div v-for="backup in backupList" :key="backup.filename" class="backup-item">
                  <div class="backup-info">
                    <div class="backup-name">
                      {{ backup.filename }}
                    </div>
                    <div class="backup-meta">
                      {{ formatTime(backup.createdAt) }}
                      <span class="meta-dot">·</span>
                      {{ formatSize(backup.dbSize) }}
                      <span v-if="backup.encrypted" class="meta-dot">·</span>
                      <span v-if="backup.encrypted" class="badge-encrypted">已加密</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 消息提示 -->
            <div v-if="errorMsg" class="message message-error">
              {{ errorMsg }}
            </div>
            <div v-if="successMsg" class="message message-success">
              {{ successMsg }}
            </div>
          </div>

          <div class="modal-footer">
            <button class="btn" @click="handleClose">关闭</button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 恢复确认弹窗 -->
    <Transition name="modal">
      <div v-if="confirmVisible" class="modal-overlay" @click.self="confirmVisible = false">
        <div class="modal-box confirm-box">
          <div class="modal-header">
            <h2 class="modal-title">确认恢复</h2>
          </div>
          <div class="modal-body">
            <p class="confirm-text">
              恢复将替换当前所有简历数据。恢复前会自动创建一个备份文件，以防万一。
            </p>
            <p class="confirm-text warn">确定要继续吗？</p>
          </div>
          <div class="modal-footer">
            <button class="btn" @click="confirmVisible = false">取消</button>
            <button class="btn btn-danger" :disabled="isImporting" @click="confirmImport">
              确认恢复
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useBackup } from '../composables/useBackup'
import { useToast } from '../composables/useToast'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

const toast = useToast()

const {
  isExporting,
  isImporting,
  isLoading,
  backupList,
  error: compError,
  success: compSuccess,
  exportBackup,
  importBackup,
  loadBackupList,
  getAutoBackupSettings,
  setAutoBackupSettings,
  formatSize,
  formatTime,
} = useBackup()

const exportPassword = ref('')
const importPassword = ref('')
const confirmVisible = ref(false)
const autoEnabled = ref(false)
const autoInterval = ref('daily')
const isSavingAuto = ref(false)

const errorMsg = ref('')
const successMsg = ref('')

const intervalLabels: Record<string, string> = {
  daily: '每日',
  weekly: '每周',
  monthly: '每月',
}

const intervalLabel = computed(() => intervalLabels[autoInterval.value] || autoInterval.value)

function handleRefresh() {
  errorMsg.value = ''
  successMsg.value = ''
  loadBackupList()
}

// 加载备份列表和自动备份设置
watch(() => props.visible, async (visible) => {
  if (visible) {
    handleRefresh()
    await loadAutoBackupSettings()
  }
})

async function loadAutoBackupSettings() {
  try {
    const settings = await getAutoBackupSettings()
    autoEnabled.value = settings.enabled === 'true'
    autoInterval.value = settings.interval || 'daily'
  } catch {
    // ignore
  }
}

async function handleAutoBackupChange() {
  isSavingAuto.value = true
  errorMsg.value = ''
  successMsg.value = ''
  try {
    await setAutoBackupSettings(autoEnabled.value, autoInterval.value)
    successMsg.value = `自动备份${autoEnabled.value ? '已启用' : '已禁用'}`
    toast.success(autoEnabled.value ? '自动备份已启用' : '自动备份已禁用')
  } catch (err) {
    errorMsg.value = String(err)
    toast.error('保存设置失败', String(err))
  } finally {
    isSavingAuto.value = false
  }
}

async function handleExport() {
  errorMsg.value = ''
  successMsg.value = ''
  try {
    const result = await exportBackup(exportPassword.value)
    if (result) {
      successMsg.value = `备份导出成功`
      toast.success('备份已创建', result)
      exportPassword.value = ''
    }
  } catch (err) {
    errorMsg.value = String(err)
    toast.error('备份导出失败', String(err))
  }
}

function handleImport() {
  errorMsg.value = ''
  successMsg.value = ''
  confirmVisible.value = true
}

async function confirmImport() {
  confirmVisible.value = false
  try {
    await importBackup(undefined, importPassword.value)
    toast.success('数据已恢复')
    // Page will reload on success
  } catch (err) {
    toast.error('数据恢复失败', String(err))
  }
}

function handleClose() {
  errorMsg.value = ''
  successMsg.value = ''
  emit('close')
}

// 同步错误/成功消息
watch(compError, (e) => {
  if (e) errorMsg.value = e
})
watch(compSuccess, (s) => {
  if (s) successMsg.value = s
})
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

.modal-box {
  background: var(--ui-bg-secondary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-lg);
  box-shadow: var(--ui-shadow-lg);
  width: 90%;
  max-width: 520px;
  max-height: 90vh;
  overflow-y: auto;
}

.confirm-box {
  max-width: 400px;
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
  color: var(--ui-text-primary);
  margin: 0;
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--ui-text-tertiary);
  font-size: 20px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: var(--ui-radius-sm);
  line-height: 1;
}

.close-btn:hover {
  color: var(--ui-text-primary);
  background: var(--ui-bg-hover);
}

.modal-body {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--ui-text-primary);
}

.section-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-desc {
  font-size: 12px;
  color: var(--ui-text-tertiary);
  line-height: 1.5;
}

.section-desc.muted {
  color: var(--ui-text-tertiary);
  font-size: 11px;
}

.section-desc.warning {
  color: var(--ui-warning);
}

.divider {
  height: 1px;
  background: var(--ui-border);
  margin: 4px 0;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-group label {
  font-size: 12px;
  font-weight: 500;
  color: var(--ui-text-secondary);
}

.form-group input[type="text"],
.form-group input[type="password"] {
  padding: 8px 10px;
  background: var(--ui-input-bg);
  border: 1px solid var(--ui-input-border);
  border-radius: var(--ui-radius-sm);
  color: var(--ui-text-primary);
  font-size: 13px;
  outline: none;
  transition: border-color var(--ui-transition-fast);
}

.form-group input:focus {
  border-color: var(--ui-accent);
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
  background: var(--ui-border);
  border-color: var(--ui-border-hover);
}

.btn-sm {
  padding: 4px 10px;
  font-size: 12px;
}

.btn-primary {
  background: var(--ui-accent);
  border-color: var(--ui-accent);
  color: var(--ui-text-inverse);
}

.btn-primary:hover:not(:disabled) {
  background: var(--ui-accent-hover);
  border-color: var(--ui-accent-hover);
}

.btn-danger {
  background: var(--ui-danger);
  border-color: var(--ui-danger);
  color: var(--ui-text-inverse);
}

.btn-danger:hover:not(:disabled) {
  background: var(--ui-danger-hover);
  border-color: var(--ui-danger-hover);
}

.empty-state {
  padding: 16px;
  text-align: center;
  color: var(--ui-text-tertiary);
  font-size: 13px;
  background: var(--ui-input-bg);
  border-radius: var(--ui-radius-sm);
  border: 1px solid var(--ui-border);
}

.backup-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 200px;
  overflow-y: auto;
}

.backup-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  background: var(--ui-input-bg);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
}

.backup-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.backup-name {
  font-size: 12px;
  color: var(--ui-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 320px;
}

.backup-meta {
  font-size: 11px;
  color: var(--ui-text-tertiary);
  display: flex;
  align-items: center;
  gap: 4px;
}

.meta-dot {
  color: var(--ui-text-tertiary);
}

.badge-encrypted {
  background: var(--ui-bg-active);
  color: var(--ui-success);
  padding: 1px 5px;
  border-radius: 3px;
  font-size: 10px;
  font-weight: 500;
}

.message {
  padding: 8px 12px;
  border-radius: var(--ui-radius-sm);
  font-size: 12px;
  word-break: break-all;
}

.message-error {
  background: var(--ui-bg-active);
  border: 1px solid var(--ui-danger);
  color: var(--ui-danger);
}

.message-success {
  background: var(--ui-bg-active);
  border: 1px solid var(--ui-success);
  color: var(--ui-success);
}

.confirm-text {
  font-size: 13px;
  color: var(--ui-text-secondary);
  margin: 0 0 8px;
  line-height: 1.6;
}

.confirm-text.warn {
  color: var(--ui-warning);
  font-weight: 500;
}

/* Toggle switch */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 36px;
  height: 20px;
  flex-shrink: 0;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  inset: 0;
  background-color: var(--ui-border);
  transition: var(--ui-transition-fast);
  border-radius: 20px;
}

.toggle-slider::before {
  position: absolute;
  content: "";
  height: 14px;
  width: 14px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: var(--ui-transition-fast);
  border-radius: 50%;
}

.toggle-switch input:checked + .toggle-slider {
  background-color: var(--ui-accent);
}

.toggle-switch input:checked + .toggle-slider::before {
  transform: translateX(16px);
}

/* Select input */
.select-input {
  padding: 8px 10px;
  background: var(--ui-input-bg);
  border: 1px solid var(--ui-input-border);
  border-radius: var(--ui-radius-sm);
  color: var(--ui-text-primary);
  font-size: 13px;
  outline: none;
  transition: border-color var(--ui-transition-fast);
  cursor: pointer;
  appearance: auto;
}

.select-input:focus {
  border-color: var(--ui-accent);
}

/* Transition */
.modal-enter-active {
  animation: modal-in var(--ui-transition-fast) ease-out;
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

.btn-spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  margin-right: 6px;
  vertical-align: middle;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
