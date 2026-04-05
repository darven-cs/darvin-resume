<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-overlay" @click.self="handleClose">
      <div class="modal-container">
        <!-- Header -->
        <div class="modal-header">
          <span class="modal-title">📥 导入旧简历</span>
          <button class="close-btn" @click="handleClose" title="关闭">×</button>
        </div>

        <!-- Input Section -->
        <div class="modal-body">
          <!-- Paste Area -->
          <div class="paste-section">
            <div class="section-label">请粘贴您的旧简历内容（纯文本或 Markdown）</div>
            <textarea
              v-model="pasteContent"
              class="paste-textarea"
              placeholder="# 张三
邮箱: zhangsan@email.com
手机: 138xxxx8888

## 教育背景
xxx大学 | 计算机科学 | 2020-2024

## 工作经历
xxx公司 | 前端开发工程师 | 2024-至今
..."
              :disabled="isParsing"
            />
          </div>

          <!-- Job Target Input -->
          <div class="input-row">
            <label class="input-label">目标岗位:</label>
            <input
              v-model="jobTarget"
              class="text-input"
              type="text"
              placeholder="例如：前端开发工程师"
              :disabled="isParsing"
            />
          </div>

          <!-- Full Context Toggle -->
          <div class="input-row">
            <label class="checkbox-label">
              <input
                v-model="includeFullContext"
                type="checkbox"
                :disabled="isParsing"
              />
              <span>包含全文参考（AI 会参考完整简历优化解析）</span>
            </label>
          </div>

          <!-- Parse Button -->
          <div class="action-row">
            <button
              class="btn-primary"
              :disabled="!canParse || isParsing"
              @click="handleParse"
            >
              <span v-if="isParsing" class="spinner" />
              <span v-else>开始解析 ▶</span>
            </button>
          </div>

          <!-- Error Message -->
          <div v-if="errorMessage" class="error-message">
            {{ errorMessage }}
          </div>

          <!-- Parsing Status -->
          <div v-if="isParsing" class="parsing-status">
            正在解析简历内容，请稍候...
          </div>

          <!-- Preview Section -->
          <div v-if="previewMode && parsedData" class="preview-section">
            <div class="section-label">解析结果预览:</div>
            <div class="preview-content">
              <div class="preview-item">
                <span class="preview-label">姓名:</span>
                <span class="preview-value">{{ parsedData.name ?? '未提取' }}</span>
              </div>
              <div class="preview-item">
                <span class="preview-label">邮箱:</span>
                <span class="preview-value">{{ parsedData.email ?? '未提取' }}</span>
              </div>
              <div class="preview-item">
                <span class="preview-label">手机:</span>
                <span class="preview-value">{{ parsedData.phone ?? '未提取' }}</span>
              </div>
              <div class="preview-item">
                <span class="preview-label">教育经历:</span>
                <span class="preview-value">
                  {{ parsedData.education?.length ?? 0 }} 项
                  <template v-if="parsedData.education?.[0]">
                    — {{ parsedData.education[0].school }}
                  </template>
                </span>
              </div>
              <div class="preview-item">
                <span class="preview-label">工作经历:</span>
                <span class="preview-value">
                  {{ parsedData.experience?.length ?? 0 }} 项
                  <template v-if="parsedData.experience?.[0]">
                    — {{ parsedData.experience[0].company }}
                  </template>
                </span>
              </div>
              <div class="preview-item">
                <span class="preview-label">项目经历:</span>
                <span class="preview-value">
                  {{ parsedData.projects?.length ?? 0 }} 项
                </span>
              </div>
              <div class="preview-item">
                <span class="preview-label">技能:</span>
                <span class="preview-value">
                  {{ parsedData.skills?.length ?? 0 }} 项
                  <template v-if="parsedData.skills?.length">
                    — {{ parsedData.skills.slice(0, 3).join(', ') }}
                    <template v-if="(parsedData.skills?.length ?? 0) > 3">...</template>
                  </template>
                </span>
              </div>
              <div v-if="parsedData.summary" class="preview-item preview-item--full">
                <span class="preview-label">自我评价:</span>
                <span class="preview-value">{{ parsedData.summary }}</span>
              </div>

              <!-- Schema Validation Warning -->
              <div v-if="validationWarnings.length > 0" class="validation-warnings">
                <div class="warning-title">⚠ 部分字段解析缺失:</div>
                <div v-for="warn in validationWarnings" :key="warn" class="warning-item">
                  - {{ warn }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="modal-footer">
          <button
            v-if="previewMode"
            class="btn-primary"
            :disabled="!parsedData || isParsing"
            @click="handleConfirm"
          >
            确认导入
          </button>
          <button
            class="btn-secondary"
            @click="handleClose"
          >
            取消
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { sendMessageSync } from '../services/ai'
import { jsonToMarkdown } from '../utils/resume'
import type { ParsedResume } from '../types/resume'
import { validateParsedResume } from '../types/resume'

const props = defineProps<{
  visible: boolean
  resumeId: string
}>()

const emit = defineEmits<{
  close: []
  import: [markdown: string, jobTarget: string]
}>()

// State
const pasteContent = ref('')
const jobTarget = ref('')
const includeFullContext = ref(false)
const parsedData = ref<ParsedResume | null>(null)
const parsedMarkdown = ref('')
const isParsing = ref(false)
const previewMode = ref(false)
const errorMessage = ref('')

// Validation warnings for missing fields
const validationWarnings = computed<string[]>(() => {
  if (!parsedData.value) return []
  const warnings: string[] = []
  if (!parsedData.value.name) warnings.push('姓名')
  if (!parsedData.value.email) warnings.push('邮箱')
  if (!parsedData.value.education?.length) warnings.push('教育经历')
  if (!parsedData.value.experience?.length && !parsedData.value.projects?.length) {
    warnings.push('工作经历或项目经历')
  }
  if (!parsedData.value.skills?.length) warnings.push('技能')
  return warnings
})

// Can parse: non-empty content
const canParse = computed(() => pasteContent.value.trim().length > 0)

// System prompt for resume parsing
const PARSER_SYSTEM_PROMPT = `你是一个简历解析专家。请从用户提供的简历文本中提取结构化信息，并以 JSON 格式输出。

**输出格式（严格遵循）：**
{
  "name": "姓名",
  "email": "邮箱",
  "phone": "手机号",
  "education": [
    {
      "school": "学校名称",
      "major": "专业",
      "degree": "学历",
      "startDate": "开始时间",
      "endDate": "结束时间",
      "description": "其他描述（如GPA、荣誉等）"
    }
  ],
  "experience": [
    {
      "company": "公司名称",
      "position": "职位",
      "startDate": "开始时间",
      "endDate": "结束时间",
      "description": "工作内容和成就"
    }
  ],
  "projects": [
    {
      "name": "项目名称",
      "role": "角色",
      "techStack": ["技术栈"],
      "description": "项目描述和你的贡献"
    }
  ],
  "skills": ["技能列表"],
  "summary": "自我评价"
}

**规则：**
- 只输出 JSON，不要任何解释或 markdown 代码块标记
- 字段为空或不适用则输出 null
- 日期格式：YYYY-MM
- 提取尽可能多的信息，但不要臆造`

/**
 * Extracts JSON from AI response text, handling common formats.
 * Tries to find JSON between code blocks or as plain text.
 */
function extractJSON(text: string): string {
  // Try to find JSON in code blocks first
  const codeBlockMatch = text.match(/```(?:json)?\s*\n?([\s\S]*?)\n?```/)
  if (codeBlockMatch) {
    return codeBlockMatch[1].trim()
  }

  // Try to find bare JSON object
  const jsonStart = text.indexOf('{')
  const jsonEnd = text.lastIndexOf('}') + 1
  if (jsonStart >= 0 && jsonEnd > jsonStart) {
    return text.substring(jsonStart, jsonEnd)
  }

  return text.trim()
}

// Parse resume
async function handleParse() {
  if (!canParse.value || isParsing.value) return

  errorMessage.value = ''
  isParsing.value = true
  previewMode.value = false
  parsedData.value = null
  parsedMarkdown.value = ''

  const operationId = `parse-${crypto.randomUUID()}`
  const userPrompt = `请解析以下简历内容：

${pasteContent.value.trim()}`

  try {
    const response = await sendMessageSync(
      operationId,
      userPrompt,
      jobTarget.value,
      includeFullContext.value
    )

    // Extract JSON from response
    const jsonText = extractJSON(response)

    // Parse JSON
    let parsed: unknown
    try {
      parsed = JSON.parse(jsonText)
    } catch {
      errorMessage.value = '解析失败：AI 返回的内容无法解析为 JSON 格式。请检查简历内容或稍后重试。'
      isParsing.value = false
      return
    }

    // Validate schema
    if (!validateParsedResume(parsed)) {
      errorMessage.value = '解析失败：返回的数据结构不符合预期。允许跳过缺失字段。'
      // Still show the data so user can review
    }

    parsedData.value = parsed as ParsedResume

    // Generate Markdown from parsed data
    parsedMarkdown.value = jsonToMarkdown(parsed as ParsedResume)
    previewMode.value = true

  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err)
    if (msg.includes('API key') || msg.includes('auth') || msg.includes('认证')) {
      errorMessage.value = 'AI 未配置或 API Key 无效。请先在设置中配置 AI。'
    } else if (msg.includes('timeout') || msg.includes('超时')) {
      errorMessage.value = '解析超时，请稍后重试或减少简历内容。'
    } else {
      errorMessage.value = `解析失败: ${msg}`
    }
  } finally {
    isParsing.value = false
  }
}

// Confirm import
function handleConfirm() {
  if (!parsedMarkdown.value) return
  emit('import', parsedMarkdown.value, jobTarget.value)
  handleClose()
}

// Close and reset
function handleClose() {
  emit('close')
  // Reset state after close animation
  setTimeout(() => {
    pasteContent.value = ''
    jobTarget.value = ''
    includeFullContext.value = false
    parsedData.value = null
    parsedMarkdown.value = ''
    isParsing.value = false
    previewMode.value = false
    errorMessage.value = ''
  }, 200)
}

// Reset when modal becomes invisible
watch(() => props.visible, (v) => {
  if (!v) {
    handleClose()
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9990;
}

.modal-container {
  width: 680px;
  max-width: 95vw;
  max-height: 90vh;
  background: #1e1e1e;
  border: 1px solid #3c3c3c;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #252526;
  border-bottom: 1px solid #3c3c3c;
}

.modal-title {
  font-size: 14px;
  font-weight: 600;
  color: #e0e0e0;
}

.close-btn {
  background: transparent;
  border: none;
  color: #8b949e;
  font-size: 20px;
  cursor: pointer;
  padding: 0 4px;
  line-height: 1;
  transition: color 0.15s;
}

.close-btn:hover {
  color: #e0e0e0;
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid #3c3c3c;
  background: #252526;
}

/* Paste Section */
.paste-section {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.section-label {
  font-size: 12px;
  color: #8b949e;
  font-weight: 500;
}

.paste-textarea {
  width: 100%;
  min-height: 180px;
  max-height: 300px;
  padding: 10px 12px;
  background: #2d2d2d;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  color: #e0e0e0;
  font-family: 'SF Mono', 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  line-height: 1.5;
  resize: vertical;
  box-sizing: border-box;
  transition: border-color 0.15s;
}

.paste-textarea:focus {
  outline: none;
  border-color: #0078d4;
}

.paste-textarea:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.paste-textarea::placeholder {
  color: #5a5a5a;
}

/* Input Row */
.input-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.input-label {
  font-size: 13px;
  color: #8b949e;
  min-width: 70px;
  flex-shrink: 0;
}

.text-input {
  flex: 1;
  padding: 6px 10px;
  background: #2d2d2d;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  color: #e0e0e0;
  font-size: 13px;
  transition: border-color 0.15s;
}

.text-input:focus {
  outline: none;
  border-color: #0078d4;
}

.text-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.text-input::placeholder {
  color: #5a5a5a;
}

/* Checkbox */
.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #8b949e;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 15px;
  height: 15px;
  accent-color: #0078d4;
  cursor: pointer;
}

.checkbox-label input:disabled {
  cursor: not-allowed;
}

/* Action Row */
.action-row {
  display: flex;
  justify-content: flex-end;
}

/* Buttons */
.btn-primary,
.btn-secondary {
  padding: 8px 20px;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s, opacity 0.15s;
  display: flex;
  align-items: center;
  gap: 6px;
  border: none;
}

.btn-primary {
  background: #0e639c;
  color: #ffffff;
}

.btn-primary:hover:not(:disabled) {
  background: #1177bb;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: #3c3c3c;
  color: #cccccc;
}

.btn-secondary:hover {
  background: #4a4a4a;
}

/* Spinner */
.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Error Message */
.error-message {
  padding: 8px 12px;
  background: rgba(218, 54, 51, 0.15);
  border: 1px solid rgba(218, 54, 51, 0.4);
  border-radius: 4px;
  color: #f48771;
  font-size: 12px;
}

/* Parsing Status */
.parsing-status {
  padding: 8px 12px;
  background: rgba(14, 99, 156, 0.15);
  border: 1px solid rgba(0, 120, 212, 0.3);
  border-radius: 4px;
  color: #79c0ff;
  font-size: 12px;
  text-align: center;
}

/* Preview Section */
.preview-section {
  border-top: 1px solid #3c3c3c;
  padding-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.preview-content {
  background: #252526;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.preview-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 12px;
}

.preview-item--full {
  flex-direction: column;
  gap: 2px;
}

.preview-label {
  color: #8b949e;
  min-width: 60px;
  flex-shrink: 0;
}

.preview-value {
  color: #e0e0e0;
  word-break: break-all;
}

/* Validation Warnings */
.validation-warnings {
  margin-top: 8px;
  padding: 8px;
  background: rgba(203, 166, 47, 0.1);
  border: 1px solid rgba(203, 166, 47, 0.3);
  border-radius: 4px;
}

.warning-title {
  font-size: 12px;
  color: #e3b341;
  margin-bottom: 4px;
}

.warning-item {
  font-size: 11px;
  color: #8b949e;
  padding-left: 8px;
}
</style>
