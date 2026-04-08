<template>
  <Teleport to="body">
    <!-- 侧边栏遮罩 -->
    <Transition name="sidebar-overlay">
      <div v-if="visible" class="sidebar-overlay" @click.self="$emit('close')" />
    </Transition>

    <!-- 侧边栏面板 -->
    <Transition name="sidebar-slide">
      <div v-if="visible" class="snapshot-sidebar">
        <!-- 头部 -->
        <div class="sidebar-header">
          <span class="sidebar-title">版本快照</span>
          <button class="close-btn" @click="$emit('close')" title="关闭">
            <span>&times;</span>
          </button>
        </div>

        <!-- 操作栏 -->
        <div class="sidebar-actions">
          <button class="action-btn action-btn-primary" @click="openCreateDialog">
            <span class="btn-icon">+</span>
            创建快照
          </button>
          <button
            v-if="selectedSnapshots.length === 2"
            class="action-btn action-btn-compare"
            @click="handleCompare"
          >
            对比
          </button>
        </div>

        <!-- 快照列表 -->
        <div class="snapshot-list">
          <!-- 加载状态 -->
          <div v-if="loading" class="loading-state">
            <span>加载中...</span>
          </div>

          <!-- 空状态 -->
          <div v-else-if="snapshots.length === 0" class="empty-state">
            <div class="empty-icon">&#128247;</div>
            <div class="empty-text">暂无版本快照</div>
            <div class="empty-hint">点击上方「创建快照」按钮保存当前版本</div>
          </div>

          <!-- 快照列表 -->
          <div v-else class="snapshot-items">
            <div
              v-for="snap in snapshots"
              :key="snap.id"
              class="snapshot-item"
              :class="{ selected: selectedSnapshots.includes(snap.id) }"
              @click="toggleSelection(snap.id)"
            >
              <div class="item-checkbox">
                <input
                  type="checkbox"
                  :checked="selectedSnapshots.includes(snap.id)"
                  @click.stop
                  @change="toggleSelection(snap.id)"
                />
              </div>
              <div class="item-content">
                <div class="item-header">
                  <span class="item-label">{{ snap.label || '未命名快照' }}</span>
                  <span class="trigger-badge" :class="getTriggerClass(snap.triggerType)">
                    {{ getTriggerLabel(snap.triggerType) }}
                  </span>
                </div>
                <div v-if="snap.note" class="item-note">{{ snap.note }}</div>
                <div class="item-time">{{ formatTime(snap.createdAt) }}</div>
              </div>
              <div class="item-actions" @click.stop>
                <button class="item-btn" @click="handleView(snap.id)" title="查看详情">查看</button>
                <button class="item-btn" @click="handleRollbackConfirm(snap.id)" title="回滚到此版本">回滚</button>
                <button class="item-btn item-btn-danger" @click="handleDelete(snap.id)" title="删除快照">删除</button>
              </div>
            </div>
          </div>
        </div>

        <!-- 创建快照对话框 -->
        <div v-if="showCreateDialog" class="dialog-overlay" @click.self="closeCreateDialog">
          <div class="create-dialog">
            <div class="dialog-header">
              <span class="dialog-title">创建快照</span>
              <button class="close-btn" @click="closeCreateDialog">&times;</button>
            </div>
            <div class="dialog-body">
              <div class="form-group">
                <label for="snap-label">标签 <span class="required">*</span></label>
                <input
                  id="snap-label"
                  v-model="createForm.label"
                  type="text"
                  placeholder="例如：面试前最终版"
                  maxlength="20"
                />
                <span class="char-count">{{ createForm.label.length }}/20</span>
              </div>
              <div class="form-group">
                <label for="snap-note">备注</label>
                <textarea
                  id="snap-note"
                  v-model="createForm.note"
                  placeholder="可选备注信息"
                  maxlength="100"
                  rows="3"
                />
                <span class="char-count">{{ createForm.note.length }}/100</span>
              </div>
            </div>
            <div class="dialog-footer">
              <button class="dialog-btn dialog-btn-cancel" @click="closeCreateDialog">取消</button>
              <button
                class="dialog-btn dialog-btn-confirm"
                :disabled="!createForm.label.trim()"
                @click="handleCreate"
              >
                确定
              </button>
            </div>
          </div>
        </div>

        <!-- 回滚确认对话框 -->
        <div v-if="showRollbackDialog" class="dialog-overlay" @click.self="closeRollbackDialog">
          <div class="create-dialog">
            <div class="dialog-header">
              <span class="dialog-title">确认回滚</span>
              <button class="close-btn" @click="closeRollbackDialog">&times;</button>
            </div>
            <div class="dialog-body">
              <div class="rollback-warning">
                <div class="warning-icon">&#9888;</div>
                <div class="warning-text">
                  确定要回滚到版本「{{ rollbackTargetLabel }}」吗？
                </div>
                <div class="warning-hint">回滚后将自动创建当前版本的备份快照</div>
              </div>
            </div>
            <div class="dialog-footer">
              <button class="dialog-btn dialog-btn-cancel" @click="closeRollbackDialog">取消</button>
              <button class="dialog-btn dialog-btn-confirm rollback-confirm" @click="handleRollback">确认回滚</button>
            </div>
          </div>
        </div>

        <!-- 删除确认对话框 -->
        <div v-if="showDeleteDialog" class="dialog-overlay" @click.self="closeDeleteDialog">
          <div class="create-dialog">
            <div class="dialog-header">
              <span class="dialog-title">确认删除</span>
              <button class="close-btn" @click="closeDeleteDialog">&times;</button>
            </div>
            <div class="dialog-body">
              <div class="rollback-warning">
                <div class="warning-icon danger">&#10060;</div>
                <div class="warning-text">确定要删除此快照吗？此操作不可恢复。</div>
              </div>
            </div>
            <div class="dialog-footer">
              <button class="dialog-btn dialog-btn-cancel" @click="closeDeleteDialog">取消</button>
              <button class="dialog-btn dialog-btn-danger" @click="handleDeleteConfirm">删除</button>
            </div>
          </div>
        </div>

        <!-- 对比视图 -->
        <div v-if="showDiffView" class="dialog-overlay" @click.self="closeDiffView">
          <div class="diff-dialog">
            <div class="dialog-header">
              <span class="dialog-title">版本对比</span>
              <button class="close-btn" @click="closeDiffView">&times;</button>
            </div>
            <div class="diff-content" ref="diffContentRef">
              <div class="diff-stats">
                <span class="stat-added">+{{ diffStats.addedLines }} 行</span>
                <span class="stat-removed">-{{ diffStats.removedLines }} 行</span>
              </div>
              <div class="diff-lines">
                <div
                  v-for="(line, index) in diffLines"
                  :key="index"
                  :class="getDiffLineClass(line.type)"
                >
                  <span class="line-prefix">{{ getLinePrefix(line.type) }}</span>
                  <span class="line-text">{{ line.content }}</span>
                </div>
                <div v-if="diffLines.length === 0" class="diff-empty">两个版本内容相同</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 查看快照详情对话框 -->
        <div v-if="showViewDialog" class="dialog-overlay" @click.self="closeViewDialog">
          <div class="view-dialog">
            <div class="dialog-header">
              <span class="dialog-title">{{ viewSnapshot?.label || '快照详情' }}</span>
              <button class="close-btn" @click="closeViewDialog">&times;</button>
            </div>
            <div class="dialog-body view-body">
              <div class="view-meta">
                <span class="trigger-badge" :class="getTriggerClass(viewSnapshot?.triggerType || '')">
                  {{ getTriggerLabel(viewSnapshot?.triggerType || '') }}
                </span>
                <span class="view-time">{{ formatTime(viewSnapshot?.createdAt || '') }}</span>
              </div>
              <div v-if="viewSnapshot?.note" class="view-note">
                <strong>备注：</strong>{{ viewSnapshot.note }}
              </div>
              <div class="view-preview">
                <pre class="markdown-preview">{{ viewSnapshot?.markdownContent }}</pre>
              </div>
            </div>
            <div class="dialog-footer">
              <button class="dialog-btn dialog-btn-cancel" @click="closeViewDialog">关闭</button>
              <button class="dialog-btn dialog-btn-confirm" @click="handleRollbackFromView">回滚到此版本</button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { diffLines as computeDiffLines, Change } from 'diff'
import { useSnapshot } from '@/composables/useSnapshot'
import { useToast } from '@/composables/useToast'
import type { Snapshot } from '@/types/snapshot'

const toast = useToast()

interface Props {
  visible: boolean
  resumeId: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  rollback: [snapshotId: string]
}>()

const {
  snapshots,
  isLoading: loading,
  loadSnapshots,
  createSnapshot,
  rollback,
  deleteSnapshot,
  getSnapshot,
  diffSnapshots,
} = useSnapshot()

// 加载快照列表
const loadList = async () => {
  if (props.resumeId) {
    await loadSnapshots(props.resumeId)
  }
}

// 初始加载
if (props.visible) {
  loadList()
}

// 监听 visible 变化
watch(() => props.visible, (val) => {
  if (val) {
    loadList()
  }
})

// 对比相关
const selectedSnapshots = ref<string[]>([])
const showDiffView = ref(false)
const diffLines = ref<{ type: 'added' | 'removed' | 'unchanged'; content: string }[]>([])
const diffStats = ref({ addedLines: 0, removedLines: 0 })

// 创建快照相关
const showCreateDialog = ref(false)
const createForm = ref({ label: '', note: '' })

// 回滚相关
const showRollbackDialog = ref(false)
const rollbackTargetId = ref('')
const rollbackTargetLabel = ref('')

// 删除相关
const showDeleteDialog = ref(false)
const deleteTargetId = ref('')

// 查看相关
const showViewDialog = ref(false)
const viewSnapshot = ref<Snapshot | null>(null)

const diffContentRef = ref<HTMLElement | null>(null)

// 触发类型标签
function getTriggerClass(type_: string): string {
  switch (type_) {
    case 'manual': return 'badge-manual'
    case 'auto_pdf_export': return 'badge-auto'
    case 'rollback': return 'badge-rollback'
    default: return 'badge-manual'
  }
}

function getTriggerLabel(type_: string): string {
  switch (type_) {
    case 'manual': return '手动'
    case 'auto_pdf_export': return '自动导出'
    case 'rollback': return '回滚备份'
    default: return '手动'
  }
}

// 格式化时间
function formatTime(timeStr: string): string {
  if (!timeStr) return ''
  try {
    const date = new Date(timeStr)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))

    if (days === 0) {
      const hours = Math.floor(diff / (1000 * 60 * 60))
      if (hours === 0) {
        const minutes = Math.floor(diff / (1000 * 60))
        if (minutes === 0) return '刚刚'
        return `${minutes}分钟前`
      }
      return `${hours}小时前`
    } else if (days === 1) {
      return '昨天'
    } else if (days < 7) {
      return `${days}天前`
    } else {
      return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
    }
  } catch {
    return timeStr
  }
}

// 选择快照
function toggleSelection(id: string) {
  const index = selectedSnapshots.value.indexOf(id)
  if (index > -1) {
    selectedSnapshots.value.splice(index, 1)
  } else {
    if (selectedSnapshots.value.length < 2) {
      selectedSnapshots.value.push(id)
    } else {
      selectedSnapshots.value = [selectedSnapshots.value[1], id]
    }
  }
}

// 创建快照
function openCreateDialog() {
  createForm.value = { label: '', note: '' }
  showCreateDialog.value = true
}

function closeCreateDialog() {
  showCreateDialog.value = false
}

async function handleCreate() {
  if (!createForm.value.label.trim()) return
  try {
    await createSnapshot(props.resumeId, createForm.value.label.trim(), createForm.value.note.trim(), 'manual')
    closeCreateDialog()
    toast.success('快照已创建')
  } catch (err) {
    toast.error('创建快照失败', String(err))
  }
}

// 查看详情
async function handleView(id: string) {
  try {
    viewSnapshot.value = await getSnapshot(id)
    showViewDialog.value = true
  } catch (err) {
    console.error('获取快照详情失败:', err)
  }
}

function closeViewDialog() {
  showViewDialog.value = false
  viewSnapshot.value = null
}

function handleRollbackFromView() {
  if (!viewSnapshot.value) return
  rollbackTargetId.value = viewSnapshot.value.id
  rollbackTargetLabel.value = viewSnapshot.value.label
  closeViewDialog()
  showRollbackDialog.value = true
}

// 对比
async function handleCompare() {
  if (selectedSnapshots.value.length !== 2) return
  try {
    const result = await diffSnapshots(selectedSnapshots.value[0], selectedSnapshots.value[1])
    diffStats.value = result.stats

    // 从 delta 还原 diff（使用 diff 包计算可读的 diffLines）
    const snap1 = await getSnapshot(selectedSnapshots.value[0])
    const snap2 = await getSnapshot(selectedSnapshots.value[1])

    const changes: Change[] = computeDiffLines(snap1.markdownContent, snap2.markdownContent)
    const lines: { type: 'added' | 'removed' | 'unchanged'; content: string }[] = []
    for (const change of changes) {
      const parts = change.value.replace(/\n$/, '').split('\n')
      for (const part of parts) {
        if (part === '' && !change.value) continue
        lines.push({
          type: change.added ? 'added' : change.removed ? 'removed' : 'unchanged',
          content: part,
        })
      }
    }
    diffLines.value = lines
    showDiffView.value = true
  } catch (err) {
    console.error('对比失败:', err)
  }
}

function closeDiffView() {
  showDiffView.value = false
  diffLines.value = []
}

function getDiffLineClass(type_: string): string {
  switch (type_) {
    case 'added': return 'diff-line diff-line-added'
    case 'removed': return 'diff-line diff-line-removed'
    case 'unchanged': return 'diff-line diff-line-unchanged'
    default: return 'diff-line'
  }
}

function getLinePrefix(type_: string): string {
  switch (type_) {
    case 'added': return '+'
    case 'removed': return '-'
    case 'unchanged': return ' '
    default: return ' '
  }
}

// 回滚
function handleRollbackConfirm(id: string) {
  const snap = snapshots.value.find(s => s.id === id)
  rollbackTargetId.value = id
  rollbackTargetLabel.value = snap?.label || '未命名快照'
  showRollbackDialog.value = true
}

function closeRollbackDialog() {
  showRollbackDialog.value = false
  rollbackTargetId.value = ''
  rollbackTargetLabel.value = ''
}

async function handleRollback() {
  if (!rollbackTargetId.value) return
  try {
    await rollback(props.resumeId, rollbackTargetId.value)
    closeRollbackDialog()
    emit('rollback', rollbackTargetId.value)
    toast.success('已回滚到选定版本')
  } catch (err) {
    toast.error('回滚失败', String(err))
  }
}

// 删除
function handleDelete(id: string) {
  deleteTargetId.value = id
  showDeleteDialog.value = true
}

function closeDeleteDialog() {
  showDeleteDialog.value = false
  deleteTargetId.value = ''
}

async function handleDeleteConfirm() {
  if (!deleteTargetId.value) return
  try {
    await deleteSnapshot(deleteTargetId.value, props.resumeId)
    closeDeleteDialog()
    // 如果被删除的快照在选中列表中，移除
    const idx = selectedSnapshots.value.indexOf(deleteTargetId.value)
    if (idx > -1) {
      selectedSnapshots.value.splice(idx, 1)
    }
    toast.success('快照已删除')
  } catch (err) {
    toast.error('删除快照失败', String(err))
  }
}
</script>

<style scoped>
/* 遮罩 */
.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 9990;
}

/* 侧边栏 */
.snapshot-sidebar {
  position: fixed;
  top: 0;
  right: 0;
  width: 360px;
  height: 100vh;
  background: var(--ui-bg-primary);
  border-left: 1px solid var(--ui-border);
  z-index: 9991;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: var(--ui-shadow-lg);
}

/* 头部 */
.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--ui-bg-secondary);
  border-bottom: 1px solid var(--ui-border);
  flex-shrink: 0;
}

.sidebar-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--ui-text-primary);
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--ui-text-tertiary);
  font-size: 20px;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: var(--ui-radius-sm);
  line-height: 1;
}

.close-btn:hover {
  color: var(--ui-text-primary);
  background: var(--ui-bg-hover);
}

/* 操作栏 */
.sidebar-actions {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--ui-border);
  flex-shrink: 0;
}

.action-btn {
  padding: 6px 12px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  background: var(--ui-bg-tertiary);
  color: var(--ui-text-secondary);
  font-size: 12px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: all var(--ui-transition-fast);
}

.action-btn:hover {
  background: var(--ui-border);
  border-color: var(--ui-border-hover);
}

.action-btn-primary {
  background: var(--ui-accent);
  border-color: var(--ui-accent);
  color: var(--ui-text-inverse);
}

.action-btn-primary:hover {
  background: var(--ui-accent-hover);
  border-color: var(--ui-accent-hover);
}

.action-btn-compare {
  background: #7c3aed;
  border-color: #7c3aed;
  color: var(--ui-text-inverse);
}

.action-btn-compare:hover {
  background: #8b5cf6;
  border-color: #8b5cf6;
}

/* 快照列表 */
.snapshot-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: var(--ui-text-tertiary);
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
  opacity: 0.5;
}

.empty-text {
  font-size: 14px;
  color: var(--ui-text-secondary);
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 12px;
  color: var(--ui-text-tertiary);
}

/* 快照项 */
.snapshot-items {
  display: flex;
  flex-direction: column;
}

.snapshot-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--ui-border);
  cursor: pointer;
  transition: background-color var(--ui-transition-fast);
}

.snapshot-item:hover {
  background: var(--ui-bg-hover);
}

.snapshot-item.selected {
  background: rgba(64, 158, 255, 0.1);
}

.item-checkbox {
  flex-shrink: 0;
  margin-top: 2px;
}

.item-checkbox input {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.item-content {
  flex: 1;
  min-width: 0;
}

.item-header {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.item-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--ui-text-primary);
}

.trigger-badge {
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 3px;
  flex-shrink: 0;
}

.badge-manual {
  background: rgba(64, 158, 255, 0.2);
  color: var(--ui-accent);
}

.badge-auto {
  background: rgba(82, 196, 26, 0.2);
  color: var(--ui-success);
}

.badge-rollback {
  background: rgba(250, 173, 20, 0.2);
  color: var(--ui-warning);
}

.item-note {
  font-size: 11px;
  color: var(--ui-text-tertiary);
  font-style: italic;
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-time {
  font-size: 11px;
  color: var(--ui-text-tertiary);
  margin-top: 4px;
}

.item-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.item-btn {
  padding: 2px 6px;
  border: 1px solid var(--ui-border);
  border-radius: 3px;
  background: transparent;
  color: var(--ui-text-tertiary);
  font-size: 11px;
  cursor: pointer;
  transition: all var(--ui-transition-fast);
}

.item-btn:hover {
  background: var(--ui-border);
  color: var(--ui-text-secondary);
}

.item-btn-danger:hover {
  background: rgba(229, 57, 53, 0.2);
  border-color: var(--ui-danger);
  color: var(--ui-danger);
}

/* 对话框 */
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: var(--ui-overlay-bg);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.create-dialog,
.diff-dialog,
.view-dialog {
  background: var(--ui-bg-secondary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-lg);
  width: 90%;
  max-width: 480px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: var(--ui-shadow-lg);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--ui-bg-tertiary);
  border-bottom: 1px solid var(--ui-border);
}

.dialog-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--ui-text-primary);
}

.dialog-body {
  padding: 16px;
  overflow-y: auto;
  flex: 1;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--ui-border);
  background: var(--ui-bg-tertiary);
}

/* 表单 */
.form-group {
  margin-bottom: 16px;
  position: relative;
}

.form-group label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: var(--ui-text-secondary);
  margin-bottom: 6px;
}

.form-group .required {
  color: var(--ui-danger);
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 8px 10px;
  background: var(--ui-input-bg);
  border: 1px solid var(--ui-input-border);
  border-radius: var(--ui-radius-sm);
  color: var(--ui-text-primary);
  font-size: 13px;
  font-family: inherit;
  outline: none;
  transition: border-color var(--ui-transition-fast);
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: var(--ui-accent);
}

.form-group textarea {
  resize: vertical;
  min-height: 60px;
}

.char-count {
  position: absolute;
  right: 8px;
  bottom: -18px;
  font-size: 10px;
  color: var(--ui-text-tertiary);
}

/* 回滚警告 */
.rollback-warning {
  text-align: center;
  padding: 12px;
}

.warning-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.warning-icon.danger {
  color: var(--ui-danger);
}

.warning-text {
  font-size: 14px;
  color: var(--ui-text-primary);
  margin-bottom: 8px;
}

.warning-hint {
  font-size: 12px;
  color: var(--ui-text-tertiary);
}

/* 对比视图 */
.diff-content {
  max-height: 400px;
  overflow-y: auto;
}

.diff-stats {
  display: flex;
  gap: 12px;
  padding: 8px 12px;
  background: var(--ui-bg-primary);
  border-bottom: 1px solid var(--ui-border);
  margin: -16px -16px 0;
  padding: 8px 16px;
}

.stat-added {
  color: var(--ui-success);
  font-size: 12px;
  font-weight: 600;
}

.stat-removed {
  color: var(--ui-danger);
  font-size: 12px;
  font-weight: 600;
}

.diff-lines {
  font-family: 'JetBrains Mono', 'Fira Code', Consolas, monospace;
  font-size: 12px;
  line-height: 1.5;
}

.diff-line {
  display: flex;
  align-items: flex-start;
  padding: 1px 12px;
  min-height: 20px;
}

.line-prefix {
  flex-shrink: 0;
  width: 14px;
  font-weight: 600;
  user-select: none;
  text-align: center;
}

.line-text {
  white-space: pre-wrap;
  word-break: break-all;
}

.diff-line-added {
  background: rgba(82, 196, 26, 0.08);
}

.diff-line-added .line-prefix {
  color: var(--ui-success);
}

.diff-line-added .line-text {
  color: var(--ui-success);
}

.diff-line-removed {
  background: rgba(229, 57, 53, 0.08);
}

.diff-line-removed .line-prefix {
  color: var(--ui-danger);
}

.diff-line-removed .line-text {
  color: var(--ui-danger);
}

.diff-line-unchanged {
  opacity: 0.5;
}

.diff-line-unchanged .line-prefix {
  color: var(--ui-text-tertiary);
}

.diff-line-unchanged .line-text {
  color: var(--ui-text-tertiary);
}

.diff-empty {
  text-align: center;
  color: var(--ui-text-tertiary);
  padding: 20px;
}

/* 查看对话框 */
.view-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.view-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.view-time {
  font-size: 12px;
  color: var(--ui-text-tertiary);
}

.view-note {
  font-size: 12px;
  color: var(--ui-text-secondary);
  padding: 8px;
  background: var(--ui-bg-hover);
  border-radius: var(--ui-radius-sm);
}

.view-preview {
  background: var(--ui-bg-primary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  padding: 12px;
  max-height: 300px;
  overflow-y: auto;
}

.markdown-preview {
  font-family: 'JetBrains Mono', 'Fira Code', Consolas, monospace;
  font-size: 12px;
  color: var(--ui-text-secondary);
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}

/* 对话框按钮 */
.dialog-btn {
  padding: 6px 16px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  background: var(--ui-bg-tertiary);
  color: var(--ui-text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: all var(--ui-transition-fast);
}

.dialog-btn:hover:not(:disabled) {
  background: var(--ui-border);
  border-color: var(--ui-border-hover);
}

.dialog-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.dialog-btn-confirm {
  background: var(--ui-accent);
  border-color: var(--ui-accent);
  color: var(--ui-text-inverse);
}

.dialog-btn-confirm:hover:not(:disabled) {
  background: var(--ui-accent-hover);
  border-color: var(--ui-accent-hover);
}

.dialog-btn-danger {
  background: var(--ui-danger);
  border-color: var(--ui-danger);
  color: var(--ui-text-inverse);
}

.dialog-btn-danger:hover:not(:disabled) {
  background: var(--ui-danger-hover);
  border-color: var(--ui-danger-hover);
}

.rollback-confirm {
  background: var(--ui-warning);
  border-color: var(--ui-warning);
  color: var(--ui-bg-primary);
}

.rollback-confirm:hover:not(:disabled) {
  background: var(--ui-warning);
  border-color: var(--ui-warning);
}

/* 过渡动画 */
.sidebar-overlay-enter-active,
.sidebar-overlay-leave-active {
  transition: opacity var(--ui-transition-fast) ease;
}

.sidebar-overlay-enter-from,
.sidebar-overlay-leave-to {
  opacity: 0;
}

.sidebar-slide-enter-active,
.sidebar-slide-leave-active {
  transition: transform var(--ui-transition-normal) ease;
}

.sidebar-slide-enter-from,
.sidebar-slide-leave-to {
  transform: translateX(100%);
}
</style>
