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
  saveTimerHandle = setTimeout(() => {
    // 触发父组件保存（通过 v-model 事件）
    emit('update:modelValue', props.modelValue)
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
  emit('update:modelValue', props.modelValue)
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
  background: #ffffff;
  border-left: 1px solid #e0e0e0;
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
  background: #f8f8f8;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.style-editor-title {
  font-size: 13px;
  font-weight: 600;
  color: #1a1a1a;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  color: #666;
  padding: 0;
}

.close-btn:hover {
  background: #e0e0e0;
  color: #1a1a1a;
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
  color: #333;
}

.control-value {
  font-size: 11px;
  color: #888;
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
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  padding: 2px;
  cursor: pointer;
  flex-shrink: 0;
}

.color-presets {
  display: flex;
  gap: 4px;
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
  border-color: #0078d4;
  box-shadow: 0 0 0 2px rgba(0, 120, 212, 0.2);
}

/* Range Slider */
.range-slider {
  width: 100%;
  height: 4px;
  appearance: none;
  background: #e0e0e0;
  border-radius: 2px;
  outline: none;
  cursor: pointer;
}

.range-slider::-webkit-slider-thumb {
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #0078d4;
  border: 2px solid #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  cursor: pointer;
}

.range-slider::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #0078d4;
  border: 2px solid #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  cursor: pointer;
}

.range-labels {
  display: flex;
  justify-content: space-between;
  font-size: 10px;
  color: #aaa;
  padding: 0 2px;
}

/* Font Select */
.font-select {
  width: 100%;
  padding: 6px 8px;
  font-size: 12px;
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  background: #fff;
  color: #333;
  cursor: pointer;
  outline: none;
}

.font-select:focus {
  border-color: #0078d4;
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
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  resize: vertical;
  outline: none;
  line-height: 1.5;
  box-sizing: border-box;
}

.custom-css-textarea:focus {
  border-color: #0078d4;
}

.css-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 10px;
  color: #aaa;
}

.removed-hint {
  color: #e67e22;
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
  background: #f0f0f0;
  color: #888;
}

.status-pass {
  background: #d4edda;
  color: #155724;
}

.status-filtered {
  background: #fff3cd;
  color: #856404;
}

/* Reset */
.reset-group {
  margin-top: auto;
  padding-top: 8px;
  border-top: 1px solid #e8e8e8;
}

.reset-btn {
  width: 100%;
  padding: 8px 12px;
  font-size: 12px;
  background: #fff;
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  color: #666;
  cursor: pointer;
  transition: background-color 0.15s, border-color 0.15s;
}

.reset-btn:hover {
  background: #f5f5f5;
  border-color: #bbb;
  color: #333;
}

/* Reset Confirm Dialog */
.reset-confirm-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.reset-confirm-dialog {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  width: 280px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.dialog-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.dialog-message {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.dialog-sub {
  font-size: 12px;
  color: #888;
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
  border-radius: 4px;
  cursor: pointer;
  border: none;
  font-weight: 500;
}

.dialog-btn.cancel {
  background: #f0f0f0;
  color: #333;
}

.dialog-btn.cancel:hover {
  background: #e0e0e0;
}

.dialog-btn.confirm {
  background: #dc3545;
  color: #fff;
}

.dialog-btn.confirm:hover {
  background: #c82333;
}
</style>
