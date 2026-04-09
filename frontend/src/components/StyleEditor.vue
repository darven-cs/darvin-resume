<template>
  <div class="style-editor" :class="{ visible: modelValue }">
    <div class="style-editor-header">
      <span class="style-editor-title">样式调整</span>
      <button class="close-btn" @click="emit('update:modelValue', false)" title="关闭">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18" />
          <line x1="6" y1="6" x2="18" y2="18" />
        </svg>
      </button>
    </div>

    <div class="style-editor-body">
      <!-- 1. 主色调选择器 -->
      <div class="control-group">
        <div class="control-label">
          <span>主色调</span>
          <span class="control-value">{{ primaryColor }}</span>
        </div>
        <div class="color-controls">
          <input type="color" v-model="primaryColor" class="color-picker" @input="onColorChange" />
          <div class="color-presets">
            <button
              v-for="preset in colorPresets"
              :key="preset"
              class="color-preset"
              :style="{ backgroundColor: preset }"
              :class="{ active: primaryColor === preset }"
              @click="setPrimaryColor(preset)"
              :title="preset"
            />
          </div>
        </div>
      </div>

      <!-- 2. 字号滑块 -->
      <div class="control-group">
        <div class="control-label">
          <span>字号</span>
          <span class="control-value">{{ fontSize }}pt</span>
        </div>
        <input
          type="range"
          v-model.number="fontSize"
          min="8"
          max="14"
          step="0.5"
          class="range-slider"
          @input="onFontSizeChange"
        />
        <div class="range-labels">
          <span>8pt</span>
          <span>11pt</span>
          <span>14pt</span>
        </div>
      </div>

      <!-- 3. 行高滑块 -->
      <div class="control-group">
        <div class="control-label">
          <span>行高</span>
          <span class="control-value">{{ lineHeight }}</span>
        </div>
        <input
          type="range"
          v-model.number="lineHeight"
          min="1.2"
          max="2.0"
          step="0.1"
          class="range-slider"
          @input="onLineHeightChange"
        />
        <div class="range-labels">
          <span>1.2</span>
          <span>1.6</span>
          <span>2.0</span>
        </div>
      </div>

      <!-- 4. 页面边距滑块 -->
      <div class="control-group">
        <div class="control-label">
          <span>页面边距</span>
          <span class="control-value">{{ pagePadding }}mm</span>
        </div>
        <input
          type="range"
          v-model.number="pagePadding"
          min="15"
          max="30"
          step="1"
          class="range-slider"
          @input="onPagePaddingChange"
        />
        <div class="range-labels">
          <span>15mm</span>
          <span>22mm</span>
          <span>30mm</span>
        </div>
      </div>

      <!-- 5. 字体选择器 -->
      <div class="control-group">
        <div class="control-label">
          <span>字体</span>
        </div>
        <select v-model="fontFamily" class="font-select" @change="onFontFamilyChange">
          <option v-for="f in fontOptions" :key="f.value" :value="f.value">
            {{ f.label }}
          </option>
        </select>
      </div>

      <!-- 6. 自定义 CSS -->
      <div class="control-group custom-css-group">
        <div class="control-label">
          <span>自定义 CSS</span>
          <span class="css-status" :class="cssStatusClass">
            {{ cssStatusText }}
          </span>
        </div>
        <textarea
          v-model="customCSSInput"
          class="custom-css-textarea"
          placeholder=".page-content h1 { font-size: 20pt; }"
          maxlength="2000"
          @input="onCustomCSSChange"
        />
        <div class="css-meta">
          <span class="char-count">{{ customCSSInput.length }}/2000</span>
          <span v-if="removedCount > 0" class="removed-hint">
            已移除 {{ removedCount }} 条危险规则
          </span>
        </div>
      </div>

      <!-- 一键重置按钮 -->
      <div class="control-group reset-group">
        <button class="reset-btn" @click="handleReset">
          重置为默认样式
        </button>
      </div>
    </div>

    <!-- 重置确认弹窗 -->
    <div v-if="showResetConfirm" class="reset-confirm-overlay" @click.self="showResetConfirm = false">
      <div class="reset-confirm-dialog">
        <div class="dialog-icon">⚠️</div>
        <div class="dialog-message">确定要重置为默认样式吗？</div>
        <div class="dialog-sub">当前自定义样式将丢失。</div>
        <div class="dialog-actions">
          <button class="dialog-btn cancel" @click="showResetConfirm = false">取消</button>
          <button class="dialog-btn confirm" @click="confirmReset">确定重置</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { sanitizeCustomCSS, buildCSSVars, DEFAULT_CSS_VARS } from '../utils/sanitizeCSS'
import { UpdateResumeCustomCSS } from '../wailsjs/wailsjs/go/main/App'

const props = defineProps<{
  /** 控制面板显示/隐藏 */
  modelValue: boolean
  /** 当前简历 ID */
  resumeId: string
  /** 当前模板 ID */
  templateId: string
  /** 外部传入的初始 CSS（来自 useTemplate） */
  initialCss?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [visible: boolean]
}>()

// ============================================================
// 样式状态（与 CSS 变量一一对应）
// ============================================================

/** 主色调 */
const primaryColor = ref(DEFAULT_CSS_VARS['--resume-primary-color'])
/** 标题颜色跟随主色调 */
const headingColor = ref(DEFAULT_CSS_VARS['--resume-heading-color'])
/** 链接颜色 */
const linkColor = ref(DEFAULT_CSS_VARS['--resume-link-color'])
/** 背景色 */
const bgColor = ref(DEFAULT_CSS_VARS['--resume-bg-color'])
/** 字号（pt） */
const fontSize = ref(parseFloat(DEFAULT_CSS_VARS['--resume-font-size']))
/** 行高 */
const lineHeight = ref(parseFloat(DEFAULT_CSS_VARS['--resume-line-height']))
/** 页面边距（mm） */
const pagePadding = ref(parseInt(DEFAULT_CSS_VARS['--resume-padding']))
/** 字体栈 */
const fontFamily = ref(DEFAULT_CSS_VARS['--resume-font-family'])
/** 自定义 CSS 输入 */
const customCSSInput = ref('')
/** 重置确认弹窗 */
const showResetConfirm = ref(false)

// ============================================================
// 常量
// ============================================================

/** 预设颜色 */
const colorPresets = ['#1a1a1a', '#333333', '#0066cc', '#2c5f8a', '#5c3317', '#1a4d1a']

/** 字体选项 */
const fontOptions = [
  { label: '系统默认', value: '-apple-system, BlinkMacSystemFont, "PingFang SC", "Microsoft YaHei", sans-serif' },
  { label: '微软雅黑', value: '"Microsoft YaHei", "PingFang SC", sans-serif' },
  { label: '宋体', value: 'SimSun, "Songti SC", serif' },
  { label: '黑体', value: '"Microsoft YaHei", "Heiti SC", sans-serif' },
  { label: 'Georgia', value: 'Georgia, serif' },
  { label: 'Times New Roman', value: '"Times New Roman", Times, serif' },
]

// ============================================================
// CSS 校验状态
// ============================================================

/** 被移除的危险规则数量 */
const removedCount = ref(0)

/** 校验状态文本 */
const cssStatusText = computed(() => {
  if (!customCSSInput.value.trim()) return '未输入'
  if (removedCount.value > 0) return `已过滤 ${removedCount.value} 条`
  return '通过'
})

/** 校验状态样式 */
const cssStatusClass = computed(() => {
  if (!customCSSInput.value.trim()) return 'status-idle'
  if (removedCount.value > 0) return 'status-filtered'
  return 'status-pass'
})

// ============================================================
// 核心：构建并应用 CSS 变量
// ============================================================

/** 当前样式参数的 CSS 变量值 */
function buildCurrentCSSVars(): Partial<Record<string, string>> {
  return {
    '--resume-primary-color': primaryColor.value,
    '--resume-heading-color': headingColor.value,
    '--resume-link-color': linkColor.value,
    '--resume-bg-color': bgColor.value,
    '--resume-font-size': `${fontSize.value}pt`,
    '--resume-line-height': String(lineHeight.value),
    '--resume-padding': `${pagePadding.value}mm`,
    '--resume-font-family': fontFamily.value,
  }
}

/**
 * 应用 CSS 变量到页面（实时预览，不保存）
 */
function applyCSSVarsToPage(vars: Partial<Record<string, string>>) {
  const root = document.documentElement
  for (const [key, val] of Object.entries(vars)) {
    if (val !== undefined) {
      root.style.setProperty(key, val)
    }
  }
}

/**
 * 生成完整的自定义 CSS 字符串（含变量 + 用户自定义）
 * 格式: CSS Variables :root { ... } Custom CSS ...
 */
function buildFullCSS(): string {
  const cssVars = buildCSSVars(buildCurrentCSSVars())
  const sanitized = sanitizeCustomCSS(customCSSInput.value)
  return [cssVars, sanitized].filter(Boolean).join('\n')
}

// ============================================================
// 事件处理：每个控件变更后立即预览 + 保存
// ============================================================

function onColorChange() {
  applyCSSVarsToPage({
    '--resume-primary-color': primaryColor.value,
    '--resume-heading-color': primaryColor.value,
  })
  saveToBackend()
}

function setPrimaryColor(color: string) {
  primaryColor.value = color
  onColorChange()
}

function onFontSizeChange() {
  applyCSSVarsToPage({ '--resume-font-size': `${fontSize.value}pt` })
  saveToBackend()
}

function onLineHeightChange() {
  applyCSSVarsToPage({ '--resume-line-height': String(lineHeight.value) })
  saveToBackend()
}

function onPagePaddingChange() {
  applyCSSVarsToPage({ '--resume-padding': `${pagePadding.value}mm` })
  saveToBackend()
}

function onFontFamilyChange() {
  applyCSSVarsToPage({ '--resume-font-family': fontFamily.value })
  saveToBackend()
}

function onCustomCSSChange() {
  // 校验并获取被移除的规则数
  const original = customCSSInput.value
  const sanitized = sanitizeCustomCSS(original)
  removedCount.value = countRemoved(original, sanitized)

  // 实时预览校验后的 CSS
  applySanitizedCSS(sanitized)

  // 延迟保存（用户停止输入后）
  saveTimer = Date.now()
}

let saveTimer = 0
let saveTimerHandle: ReturnType<typeof setTimeout> | null = null

// ============================================================
// 保存机制
// ============================================================

/**
 * 向后端保存完整 CSS
 */
function saveToBackend() {
  // 简单防抖：500ms 内合并
  if (saveTimerHandle) clearTimeout(saveTimerHandle)
  saveTimer = Date.now()
  saveTimerHandle = setTimeout(async () => {
    const fullCSS = buildFullCSS()
    try {
      await UpdateResumeCustomCSS(props.resumeId, fullCSS)
    } catch (err) {
      console.error('[StyleEditor] 保存样式失败:', err)
    }
  }, 500)
}

/**
 * 应用校验后的自定义 CSS（实时预览用）
 */
function applySanitizedCSS(sanitized: string) {
  // 合并：CSS 变量 + 校验后的自定义 CSS
  const fullCSS = buildFullCSS()
  applyCSSToPreview(fullCSS)
}

// ============================================================
// 预览区 CSS 注入
// ============================================================

const previewStyleId = 'styleeditor-preview-css'

function applyCSSToPreview(css: string) {
  let styleEl = document.getElementById(previewStyleId) as HTMLStyleElement | null
  if (!styleEl) {
    styleEl = document.createElement('style')
    styleEl.id = previewStyleId
    document.head.appendChild(styleEl)
  }
  styleEl.textContent = css
}

// ============================================================
// 重置功能
// ============================================================

function handleReset() {
  showResetConfirm.value = true
}

function confirmReset() {
  showResetConfirm.value = false

  // 恢复默认值
  primaryColor.value = DEFAULT_CSS_VARS['--resume-primary-color']
  headingColor.value = DEFAULT_CSS_VARS['--resume-heading-color']
  linkColor.value = DEFAULT_CSS_VARS['--resume-link-color']
  bgColor.value = DEFAULT_CSS_VARS['--resume-bg-color']
  fontSize.value = parseFloat(DEFAULT_CSS_VARS['--resume-font-size'])
  lineHeight.value = parseFloat(DEFAULT_CSS_VARS['--resume-line-height'])
  pagePadding.value = parseInt(DEFAULT_CSS_VARS['--resume-padding'])
  fontFamily.value = DEFAULT_CSS_VARS['--resume-font-family']
  customCSSInput.value = ''
  removedCount.value = 0

  // 应用默认 CSS 变量
  applyCSSVarsToPage(buildCurrentCSSVars())

  // 清除预览区的自定义 CSS
  const styleEl = document.getElementById(previewStyleId)
  if (styleEl) styleEl.textContent = ''

  // 保存到后端
  const defaultCSS = buildCSSVars(buildCurrentCSSVars())
  saveTimer = Date.now()
  saveTimerHandle = setTimeout(async () => {
    try {
      await UpdateResumeCustomCSS(props.resumeId, defaultCSS)
    } catch (err) {
      console.error('[StyleEditor] 保存默认样式失败:', err)
    }
  }, 500)
}

// ============================================================
// 初始化：从初始 CSS 中解析当前值
// ============================================================

function parseInitialCSS(css: string) {
  if (!css) return
  // 简单解析 --resume-* 变量值
  const parser = /--resume-([\w-]+):\s*([^;]+);/g
  let match
  while ((match = parser.exec(css)) !== null) {
    const [, prop, value] = match
    switch (prop) {
      case 'primary-color': primaryColor.value = value.trim(); break
      case 'heading-color': headingColor.value = value.trim(); break
      case 'link-color': linkColor.value = value.trim(); break
      case 'bg-color': bgColor.value = value.trim(); break
      case 'font-size': fontSize.value = parseFloat(value.trim()); break
      case 'line-height': lineHeight.value = parseFloat(value.trim()); break
      case 'padding': pagePadding.value = parseInt(value.trim()); break
      case 'font-family': fontFamily.value = value.trim(); break
    }
  }

  // 提取自定义 CSS（变量块后的内容）
  const customPart = css.split('/* Custom CSS */')[1]
  if (customPart) {
    customCSSInput.value = customPart.trim()
  }
}

// ============================================================
// 工具函数
// ============================================================

/** 粗略估算被移除的规则数（通过行数差） */
function countRemoved(original: string, sanitized: string): number {
  const originalLines = original.split('\n').filter(l => l.trim())
  const sanitizedLines = sanitized.split('\n').filter(l => l.trim())
  return Math.max(0, originalLines.length - sanitizedLines.length)
}

// ============================================================
// 生命周期
// ============================================================

onMounted(() => {
  if (props.initialCss) {
    parseInitialCSS(props.initialCss)
  }
  // 初始化时应用当前样式
  applyCSSVarsToPage(buildCurrentCSSVars())
  if (props.initialCss) {
    applyCSSToPreview(props.initialCss)
  }
})
</script>

<style scoped>
.style-editor {
  width: 280px;
  height: 100%;
  background: var(--ui-bg-primary);
  border-left: 1px solid var(--ui-border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-shrink: 0;
}

/* Header */
.style-editor-header {
  height: 40px;
  min-height: 40px;
  padding: 0 12px;
  background: var(--ui-bg-secondary);
  border-bottom: 1px solid var(--ui-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.style-editor-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--ui-text-primary);
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: transparent;
  border: none;
  border-radius: var(--ui-radius-sm);
  cursor: pointer;
  color: var(--ui-text-tertiary);
  padding: 0;
}

.close-btn:hover {
  background: var(--ui-border);
  color: var(--ui-text-primary);
}

/* Body */
.style-editor-body {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Control Group */
.control-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.control-label {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  font-weight: 500;
  color: var(--ui-text-primary);
}

.control-value {
  font-size: 11px;
  color: var(--ui-text-tertiary);
  font-weight: 400;
  font-family: monospace;
}

/* Color Picker */
.color-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-picker {
  width: 36px;
  height: 28px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  padding: 2px;
  cursor: pointer;
  flex-shrink: 0;
}

.color-presets {
  display: flex;
  gap: var(--ui-radius-sm);
  flex-wrap: wrap;
}

.color-preset {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: transform 0.1s, border-color 0.1s;
  padding: 0;
}

.color-preset:hover {
  transform: scale(1.15);
}

.color-preset.active {
  border-color: var(--ui-accent);
  box-shadow: 0 0 0 2px var(--ui-bg-active);
}

/* Range Slider */
.range-slider {
  width: 100%;
  height: 4px;
  appearance: none;
  background: var(--ui-border);
  border-radius: 2px;
  outline: none;
  cursor: pointer;
}

.range-slider::-webkit-slider-thumb {
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--ui-accent);
  border: 2px solid var(--ui-bg-primary);
  box-shadow: var(--ui-shadow-sm);
  cursor: pointer;
}

.range-slider::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--ui-accent);
  border: 2px solid var(--ui-bg-primary);
  box-shadow: var(--ui-shadow-sm);
  cursor: pointer;
}

.range-labels {
  display: flex;
  justify-content: space-between;
  font-size: 10px;
  color: var(--ui-text-secondary);
  padding: 0 2px;
}

/* Font Select */
.font-select {
  width: 100%;
  padding: 6px 8px;
  font-size: 12px;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  background: var(--ui-bg-primary);
  color: var(--ui-text-primary);
  cursor: pointer;
  outline: none;
}

.font-select:focus {
  border-color: var(--ui-accent);
}

/* Custom CSS */
.custom-css-group {
  flex: 1;
}

.custom-css-textarea {
  width: 100%;
  min-height: 100px;
  padding: 8px;
  font-size: 11px;
  font-family: 'SF Mono', 'Consolas', monospace;
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  resize: vertical;
  outline: none;
  line-height: 1.5;
  box-sizing: border-box;
}

.custom-css-textarea:focus {
  border-color: var(--ui-accent);
}

.css-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 10px;
  color: var(--ui-text-secondary);
}

.removed-hint {
  color: var(--ui-warning);
}

.char-count {
  font-family: monospace;
}

/* CSS Status */
.css-status {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 10px;
  font-weight: 500;
}

.status-idle {
  background: var(--ui-bg-secondary);
  color: var(--ui-text-tertiary);
}

.status-pass {
  background: var(--ui-bg-active);
  color: var(--ui-success);
}

.status-filtered {
  background: var(--ui-bg-active);
  color: var(--ui-warning);
}

/* Reset */
.reset-group {
  margin-top: auto;
  padding-top: 8px;
  border-top: 1px solid var(--ui-border);
}

.reset-btn {
  width: 100%;
  padding: 8px 12px;
  font-size: 12px;
  background: var(--ui-bg-primary);
  border: 1px solid var(--ui-border);
  border-radius: var(--ui-radius-sm);
  color: var(--ui-text-tertiary);
  cursor: pointer;
  transition: background-color var(--ui-transition-fast), border-color var(--ui-transition-fast);
}

.reset-btn:hover {
  background: var(--ui-bg-secondary);
  border-color: var(--ui-border-hover);
  color: var(--ui-text-primary);
}

/* Reset Confirm Dialog */
.reset-confirm-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--ui-overlay-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.reset-confirm-dialog {
  background: var(--ui-bg-primary);
  border-radius: var(--ui-radius-lg);
  padding: 24px;
  width: 280px;
  text-align: center;
  box-shadow: var(--ui-shadow-md);
}

.dialog-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.dialog-message {
  font-size: 14px;
  font-weight: 600;
  color: var(--ui-text-primary);
  margin-bottom: 4px;
}

.dialog-sub {
  font-size: 12px;
  color: var(--ui-text-tertiary);
  margin-bottom: 16px;
}

.dialog-actions {
  display: flex;
  gap: 8px;
}

.dialog-btn {
  flex: 1;
  padding: 8px;
  font-size: 13px;
  border-radius: var(--ui-radius-sm);
  cursor: pointer;
  border: none;
  font-weight: 500;
}

.dialog-btn.cancel {
  background: var(--ui-bg-secondary);
  color: var(--ui-text-primary);
}

.dialog-btn.cancel:hover {
  background: var(--ui-border);
}

.dialog-btn.confirm {
  background: var(--ui-danger);
  color: var(--ui-text-inverse);
}

.dialog-btn.confirm:hover {
  background: var(--ui-danger-hover);
}
</style>
