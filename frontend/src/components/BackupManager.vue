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

interface Props {
  visible: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

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
  } catch {
    // error already in compError
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
      exportPassword.value = ''
    }
  } catch {
    // error already in compError
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
    // Page will reload on success
  } catch {
    // error already in compError
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
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(2px);
}

.modal-box {
  background: #252526;
  border: 1px solid #3c3c3c;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
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

.section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #e0e0e0;
}

.section-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-desc {
  font-size: 12px;
  color: #888;
  line-height: 1.5;
}

.section-desc.muted {
  color: #666;
  font-size: 11px;
}

.section-desc.warning {
  color: #e5a050;
}

.divider {
  height: 1px;
  background: #3c3c3c;
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
  color: #ccc;
}

.form-group input[type="text"],
.form-group input[type="password"] {
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

.btn-sm {
  padding: 4px 10px;
  font-size: 12px;
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

.btn-danger {
  background: #c42b1c;
  border-color: #c42b1c;
  color: #fff;
}

.btn-danger:hover:not(:disabled) {
  background: #d93d30;
  border-color: #d93d30;
}

.empty-state {
  padding: 16px;
  text-align: center;
  color: #666;
  font-size: 13px;
  background: #1e1e1e;
  border-radius: 4px;
  border: 1px solid #3c3c3c;
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
  background: #1e1e1e;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
}

.backup-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.backup-name {
  font-size: 12px;
  color: #e0e0e0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 320px;
}

.backup-meta {
  font-size: 11px;
  color: #888;
  display: flex;
  align-items: center;
  gap: 4px;
}

.meta-dot {
  color: #555;
}

.badge-encrypted {
  background: #2d4a1e;
  color: #7ec850;
  padding: 1px 5px;
  border-radius: 3px;
  font-size: 10px;
  font-weight: 500;
}

.message {
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 12px;
  word-break: break-all;
}

.message-error {
  background: #2d1f1f;
  border: 1px solid #e5484d;
  color: #ff7b7b;
}

.message-success {
  background: #1f2d1f;
  border: 1px solid #30a46c;
  color: #7ec850;
}

.confirm-text {
  font-size: 13px;
  color: #ccc;
  margin: 0 0 8px;
  line-height: 1.6;
}

.confirm-text.warn {
  color: #e5a050;
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
  background-color: #3c3c3c;
  transition: 0.2s;
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
  transition: 0.2s;
  border-radius: 50%;
}

.toggle-switch input:checked + .toggle-slider {
  background-color: #3d8bfd;
}

.toggle-switch input:checked + .toggle-slider::before {
  transform: translateX(16px);
}

/* Select input */
.select-input {
  padding: 8px 10px;
  background: #1e1e1e;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  color: #e0e0e0;
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
  cursor: pointer;
  appearance: auto;
}

.select-input:focus {
  border-color: #3d8bfd;
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
