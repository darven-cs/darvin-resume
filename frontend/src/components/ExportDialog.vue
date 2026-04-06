<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="handleClose">
        <div class="modal-box" role="dialog" aria-modal="true" aria-labelledby="export-modal-title">
          <div class="modal-header">
            <h2 id="export-modal-title" class="modal-title">导出 PDF</h2>
            <button class="close-btn" @click="handleClose" aria-label="关闭">
              <span>&times;</span>
            </button>
          </div>

          <div class="modal-body">
            <!-- 导出模式选择 -->
            <div class="section-title">导出方式</div>

            <div
              class="mode-option"
              :class="{ selected: exportMode === 'system' }"
              @click="exportMode = 'system'"
            >
              <div class="mode-radio">
                <input
                  type="radio"
                  id="mode-system"
                  value="system"
                  v-model="exportMode"
                />
                <label for="mode-system"></label>
              </div>
              <div class="mode-info">
                <div class="mode-name">系统打印（推荐）</div>
                <div class="mode-desc">无需额外下载，使用系统自带 PDF 打印机，适合大多数场景</div>
              </div>
              <div class="mode-badge">默认</div>
            </div>

            <div
              class="mode-option"
              :class="{ selected: exportMode === 'chromedp' }"
              @click="exportMode = 'chromedp'"
            >
              <div class="mode-radio">
                <input
                  type="radio"
                  id="mode-chromedp"
                  value="chromedp"
                  v-model="exportMode"
                />
                <label for="mode-chromedp"></label>
              </div>
              <div class="mode-info">
                <div class="mode-name">高级导出（Chromedp）</div>
                <div class="mode-desc">精确控制，可静默导出到指定路径，适合批量导出</div>
              </div>
              <div v-if="exportMode === 'chromedp'" class="mode-badge advanced">高级</div>
            </div>

            <!-- Chromedp 模式选项 -->
            <Transition name="slide-fade">
              <div v-if="exportMode === 'chromedp'" class="advanced-options">
                <div class="form-group">
                  <label for="output-path">输出路径</label>
                  <div class="input-row">
                    <input
                      id="output-path"
                      v-model="outputPath"
                      type="text"
                      placeholder="/home/user/resume.pdf"
                    />
                    <button class="btn btn-small" @click="selectOutputPath">选择</button>
                  </div>
                </div>

                <div class="form-group">
                  <label for="dpi">DPI（导出分辨率）</label>
                  <input
                    id="dpi"
                    v-model.number="dpi"
                    type="number"
                    min="72"
                    max="300"
                    step="1"
                    placeholder="96"
                  />
                  <span class="field-hint">范围 72~300，推荐 96</span>
                </div>
              </div>
            </Transition>

            <!-- 通用导出参数 -->
            <div class="section-title" style="margin-top: 8px;">导出参数</div>

            <div class="form-group">
              <label for="page-range">页码范围</label>
              <input
                id="page-range"
                v-model="pageRange"
                type="text"
                placeholder="all 或 1-3"
              />
              <span class="field-hint">留空表示全部页面，支持范围如 "1-3"</span>
            </div>

            <div class="form-group">
              <label class="toggle-label">
                <input
                  v-model="hidePageBreaks"
                  type="checkbox"
                  class="toggle-checkbox"
                />
                <span class="toggle-text">显示分页线</span>
                <span class="toggle-hint">在 PDF 中显示蓝色虚线分页参考线</span>
              </label>
            </div>

            <!-- 错误提示 -->
            <div v-if="errorMessage" class="error-message">
              {{ errorMessage }}
            </div>
          </div>

          <div class="modal-footer">
            <button class="btn btn-secondary" @click="handleClose" :disabled="isExporting">
              取消
            </button>
            <button
              class="btn btn-primary"
              @click="handleExport"
              :disabled="isExporting || (exportMode === 'chromedp' && !outputPath.trim())"
            >
              {{ isExporting ? '导出中...' : '导出 PDF' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ExportPDFFromHTML, ShowSaveDialog } from '../wailsjs/wailsjs/go/main/App'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  exported: []
}>()

// 导出模式
const exportMode = ref<'system' | 'chromedp'>('system')

// 导出参数
const pageRange = ref('all')
const hidePageBreaks = ref(false)
const dpi = ref(96)
const outputPath = ref('')

// 状态
const isExporting = ref(false)
const errorMessage = ref('')

function handleClose() {
  if (!isExporting.value) {
    emit('close')
  }
}

/**
 * 系统打印模式导出
 * 使用 window.print() 依赖系统原生打印对话框
 */
async function exportWithSystemPrint() {
  // 1. 应用分页线显示/隐藏
  document.body.classList.toggle('hide-page-breaks', !hidePageBreaks.value)

  // 2. 将预览区包装到 .print-container 中
  const previewContainer = document.querySelector('.preview-container')
  previewContainer?.classList.add('print-container')

  // 3. 触发打印
  window.print()

  // 4. 恢复 UI（打印对话框关闭后）
  window.addEventListener('afterprint', () => {
    previewContainer?.classList.remove('print-container')
    document.body.classList.remove('hide-page-breaks')
  }, { once: true })

  emit('exported')
}

/**
 * Chromedp 高级模式导出
 * 通过 Wails bridge 调用后端无头浏览器导出
 */
async function exportWithChromedp() {
  // 获取 A4Page 的完整 HTML 内容
  const a4PageEl = document.querySelector('.preview-container') as HTMLElement | null
  if (!a4PageEl) {
    errorMessage.value = '无法获取预览内容'
    return
  }

  // 克隆节点并移除分页线装饰
  const clone = a4PageEl.cloneNode(true) as HTMLElement
  const pageBreaks = clone.querySelectorAll('.a4-page::before, .a4-page::after')
  pageBreaks.forEach(el => el.remove())

  // 获取所有内联样式表内容
  const allStyles = Array.from(document.styleSheets)
    .filter(sheet => !sheet.href || sheet.href.startsWith(window.location.origin))
    .flatMap(sheet => {
      try {
        return Array.from(sheet.cssRules).map(rule => rule.cssText)
      } catch {
        return []
      }
    })
    .join('\n')

  // 构建完整 HTML 文档
  const htmlContent = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
${allStyles}
/* 打印样式 */
@media print {
  body > *:not(.print-container) { display: none !important; }
  .print-container { display: block !important; width: 100%; }
  .a4-page {
    width: 210mm !important;
    height: auto !important;
    min-height: 297mm;
    padding: 20mm !important;
    margin: 0 !important;
    box-shadow: none !important;
    page-break-after: always;
    break-after: page;
  }
  .a4-page:last-child { page-break-after: auto; break-after: auto; }
  .page-content > * { break-inside: avoid !important; page-break-inside: avoid !important; }
  .page-content h1, .page-content h2, .page-content h3 {
    break-after: avoid !important; page-break-after: avoid !important;
  }
  .page-content table, .page-content thead, .page-content tbody,
  .page-content tr, .page-content th, .page-content td {
    break-inside: avoid !important; page-break-inside: avoid !important;
  }
  .page-content li { break-inside: avoid !important; page-break-inside: avoid !important; }
  * { -webkit-print-color-adjust: exact !important; print-color-adjust: exact !important; }
  .page-content { color: #1a1a1a !important; }
  html, body { font-size: 10.5pt !important; }
  .hide-page-breaks .a4-page::before,
  .hide-page-breaks .a4-page::after { display: none !important; }
}
</style>
</head>
<body>
<div class="print-container">
${clone.outerHTML}
</div>
</body>
</html>`

  // 5. 调用后端 Chromedp 导出
  isExporting.value = true
  errorMessage.value = ''

  try {
    await ExportPDFFromHTML(htmlContent, outputPath.value)
    emit('exported')
    emit('close')
  } catch (err) {
    errorMessage.value = `导出失败: ${err}`
  } finally {
    isExporting.value = false
  }
}

/**
 * 选择输出路径（Chromedp 模式）
 * 通过 Wails 系统对话框选择文件保存位置
 */
async function selectOutputPath() {
  try {
    // 使用 Wails 系统文件保存对话框
    const result = await ShowSaveDialog({
      title: '保存 PDF 文件',
      defaultPath: 'resume.pdf',
      filters: [{ name: 'PDF 文件', extensions: ['pdf'] }],
    })
    if (result && result.filePath) {
      outputPath.value = result.filePath
    }
  } catch (err) {
    console.error('选择路径失败:', err)
    // fallback: 使用默认路径
    if (!outputPath.value) {
      outputPath.value = `${document.title || 'resume'}.pdf`
    }
  }
}

async function handleExport() {
  errorMessage.value = ''

  if (exportMode.value === 'system') {
    await exportWithSystemPrint()
    emit('close')
  } else {
    await exportWithChromedp()
  }
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
  max-width: 520px;
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
  gap: 12px;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  color: #888;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 4px;
}

/* 模式选择卡片 */
.mode-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 1px solid #3c3c3c;
  border-radius: 6px;
  cursor: pointer;
  transition: border-color 0.15s, background-color 0.15s;
}

.mode-option:hover {
  border-color: #555;
  background: rgba(255, 255, 255, 0.03);
}

.mode-option.selected {
  border-color: #3d8bfd;
  background: rgba(61, 139, 253, 0.08);
}

.mode-radio input {
  display: none;
}

.mode-radio label {
  display: block;
  width: 18px;
  height: 18px;
  border: 2px solid #555;
  border-radius: 50%;
  cursor: pointer;
  position: relative;
  flex-shrink: 0;
  transition: border-color 0.15s;
}

.mode-option.selected .mode-radio label {
  border-color: #3d8bfd;
}

.mode-option.selected .mode-radio label::after {
  content: '';
  position: absolute;
  top: 3px;
  left: 3px;
  width: 8px;
  height: 8px;
  background: #3d8bfd;
  border-radius: 50%;
}

.mode-info {
  flex: 1;
}

.mode-name {
  font-size: 13px;
  font-weight: 500;
  color: #e0e0e0;
}

.mode-desc {
  font-size: 11px;
  color: #888;
  margin-top: 2px;
}

.mode-badge {
  font-size: 10px;
  padding: 2px 6px;
  background: #3d8bfd;
  color: #fff;
  border-radius: 3px;
  flex-shrink: 0;
}

.mode-badge.advanced {
  background: #7c3aed;
}

/* 高级选项 */
.advanced-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid #3c3c3c;
  border-radius: 6px;
}

/* 表单 */
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

.field-hint {
  font-size: 11px;
  color: #888;
}

.input-row {
  display: flex;
  gap: 8px;
}

.input-row input {
  flex: 1;
}

.btn-small {
  padding: 6px 12px;
  font-size: 12px;
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

.error-message {
  padding: 8px 12px;
  background: #2d1f1f;
  border: 1px solid #e5484d;
  border-radius: 4px;
  color: #ff7b7b;
  font-size: 12px;
}

/* 底部按钮 */
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
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}

@keyframes modal-out {
  from { opacity: 1; transform: scale(1); }
  to { opacity: 0; transform: scale(0.95); }
}

/* 高级选项滑入动画 */
.slide-fade-enter-active {
  transition: all 0.2s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.15s ease-in;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
