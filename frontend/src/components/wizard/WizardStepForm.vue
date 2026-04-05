<template>
  <div class="step-form">
    <!-- 当前模块进度 -->
    <div class="form-header">
      <button v-if="currentIndex > 0" class="back-btn" @click="$emit('back')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="15 18 9 12 15 6" />
        </svg>
        返回
      </button>
      <div class="form-progress">
        <span class="progress-text">{{ currentIndex + 1 }}/{{ modules.length }} {{ currentModuleLabel }}</span>
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
        </div>
      </div>
    </div>

    <!-- 动态表单 -->
    <div class="form-body">
      <template v-if="currentModuleType === 'basicInfo'">
        <div v-for="field in currentFields" :key="field.key" class="field-group">
          <label class="field-label">{{ field.label }}</label>
          <textarea
            v-if="field.type === 'textarea'"
            v-model="localData[field.key]"
            class="field-textarea"
            :placeholder="field.placeholder"
            rows="3"
          ></textarea>
          <input
            v-else
            v-model="localData[field.key]"
            class="field-input"
            :type="field.type || 'text'"
            :placeholder="field.placeholder"
          />
        </div>
      </template>

      <template v-else-if="currentModuleType === 'skills'">
        <div class="items-section">
          <div v-for="(cat, catIdx) in localData.categories" :key="catIdx" class="skill-category">
            <div class="category-header">
              <span class="category-label">技能类别 {{ catIdx + 1 }}</span>
              <button v-if="localData.categories.length > 1" class="remove-btn" @click="removeSkillCategory(catIdx)" title="删除">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            </div>
            <input v-model="cat.category" class="field-input" placeholder="例如：前端开发、后端开发" />
            <textarea v-model="cat.skillsText" class="field-textarea" placeholder="输入技能，用逗号分隔，例如：Vue, React, TypeScript" rows="2"></textarea>
          </div>
          <button class="add-btn" @click="addSkillCategory">+ 添加技能类别</button>
        </div>
      </template>

      <template v-else-if="currentModuleType === 'evaluation'">
        <div class="field-group">
          <label class="field-label">自我评价</label>
          <textarea v-model="localData.content" class="field-textarea" placeholder="简要描述你的核心优势、技术热情和职业目标（建议3-5句话）" rows="6"></textarea>
        </div>
      </template>

      <!-- 通用列表模块：projects, campus, internship, awards, certificates -->
      <template v-else>
        <div class="items-section">
          <div v-for="(item, itemIdx) in localData.items" :key="itemIdx" class="list-item-card">
            <div class="item-card-header">
              <span class="item-card-title">{{ getItemTitle(item, currentModuleType) }} {{ itemIdx + 1 }}</span>
              <button v-if="localData.items.length > 1" class="remove-btn" @click="removeItem(itemIdx)" title="删除">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            </div>
            <div v-for="field in currentFields" :key="field.key" class="field-group">
              <label class="field-label">{{ field.label }}</label>
              <textarea
                v-if="field.type === 'textarea'"
                v-model="item[field.key]"
                class="field-textarea"
                :placeholder="field.placeholder"
                :rows="field.rows || 3"
              ></textarea>
              <input
                v-else
                v-model="item[field.key]"
                class="field-input"
                :type="field.type || 'text'"
                :placeholder="field.placeholder"
              />
            </div>
          </div>
          <button class="add-btn" @click="addItem">{{ addItemLabel }}</button>
        </div>
      </template>
    </div>

    <!-- AI 润色结果对比 -->
    <div v-if="polishResult" class="polish-section">
      <div class="polish-header">
        <span class="polish-title">AI 润色结果</span>
        <button class="polish-close" @click="polishResult = null">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>
      <div class="polish-compare">
        <div class="polish-column">
          <span class="polish-label">原文</span>
          <pre class="polish-content original">{{ polishResult.original }}</pre>
        </div>
        <div class="polish-column">
          <span class="polish-label">润色后</span>
          <pre class="polish-content polished">{{ polishResult.polished }}</pre>
        </div>
      </div>
      <div class="polish-actions">
        <button class="btn-secondary btn-sm" @click="polishResult = null">保持原样</button>
        <button class="btn-secondary btn-sm" @click="handleRePolish" :disabled="isPolishing">
          {{ isPolishing ? '润色中...' : '重新润色' }}
        </button>
        <button class="btn-primary btn-sm" @click="acceptPolish">接受润色</button>
      </div>
    </div>

    <!-- 底部操作按钮 -->
    <div class="form-footer">
      <button
        class="btn-polish"
        :disabled="isPolishing"
        @click="handlePolish"
      >
        <span v-if="isPolishing" class="spinner"></span>
        <span v-else>AI 润色</span>
      </button>
      <div class="footer-right">
        <button class="btn-secondary" @click="$emit('skip')">跳过</button>
        <button class="btn-primary" @click="handleNext">
          {{ isLastModule ? '预览生成' : '下一步' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { sendMessageSync } from '../../services/ai'

// ============================================================
// Props & Emits
// ============================================================

interface Props {
  modules: string[]
  currentIndex: number
  formData: Record<string, any>
  jobTarget: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:formData', moduleType: string, data: any): void
  (e: 'next'): void
  (e: 'skip'): void
  (e: 'back'): void
  (e: 'finish'): void
}>()

// ============================================================
// 表单字段配置（配置驱动）
// ============================================================

interface FieldConfig {
  key: string
  label: string
  type?: string
  placeholder: string
  rows?: number
}

const moduleFieldConfigs: Record<string, FieldConfig[]> = {
  basicInfo: [
    { key: 'name', label: '姓名', placeholder: '张三' },
    { key: 'phone', label: '电话', placeholder: '138-xxxx-xxxx' },
    { key: 'email', label: '邮箱', placeholder: 'zhangsan@email.com' },
    { key: 'github', label: 'GitHub', placeholder: 'https://github.com/username' },
    { key: 'website', label: '个人网站', placeholder: 'https://yourwebsite.com' },
    { key: 'address', label: '地址', placeholder: '北京市海淀区' },
    { key: 'summary', label: '一句话简介', type: 'textarea', placeholder: '3年经验的 Web 全栈开发者，专注 Vue/React 生态' },
  ],
  projects: [
    { key: 'name', label: '项目名称', placeholder: 'Darvin-Resume' },
    { key: 'role', label: '角色', placeholder: '前端负责人' },
    { key: 'startDate', label: '开始时间', placeholder: '2025-01' },
    { key: 'endDate', label: '结束时间', placeholder: '2025-06 或 至今' },
    { key: 'techStack', label: '技术栈', placeholder: 'Vue3, TypeScript, Go' },
    { key: 'description', label: '项目描述', type: 'textarea', placeholder: '简要描述项目和你负责的部分', rows: 3 },
    { key: 'highlights', label: '亮点（每行一条）', type: 'textarea', placeholder: '重构了XX模块，性能提升30%\n独立完成了XX功能', rows: 3 },
  ],
  campus: [
    { key: 'name', label: '活动名称', placeholder: 'ACM 算法竞赛社团' },
    { key: 'role', label: '角色', placeholder: '社长' },
    { key: 'startDate', label: '开始时间', placeholder: '2022-09' },
    { key: 'endDate', label: '结束时间', placeholder: '2023-06' },
    { key: 'description', label: '描述', type: 'textarea', placeholder: '描述你在活动中的职责和贡献', rows: 3 },
    { key: 'highlights', label: '亮点（每行一条）', type: 'textarea', placeholder: '组织了XX比赛，参与人数XX人\n获得了XX奖项', rows: 2 },
  ],
  internship: [
    { key: 'company', label: '公司', placeholder: 'XX科技有限公司' },
    { key: 'position', label: '职位', placeholder: '前端开发实习生' },
    { key: 'startDate', label: '开始时间', placeholder: '2024-03' },
    { key: 'endDate', label: '结束时间', placeholder: '2024-06' },
    { key: 'description', label: '描述', type: 'textarea', placeholder: '描述你的工作职责和成果', rows: 3 },
    { key: 'highlights', label: '亮点（每行一条）', type: 'textarea', placeholder: '独立完成了XX功能\n优化了XX，效率提升XX', rows: 2 },
  ],
  awards: [
    { key: 'name', label: '奖项名称', placeholder: '全国大学生数学建模竞赛一等奖' },
    { key: 'level', label: '级别', placeholder: '国家级 / 省级 / 校级' },
    { key: 'date', label: '日期', placeholder: '2024-03' },
    { key: 'description', label: '描述', type: 'textarea', placeholder: '简要说明获奖背景和意义', rows: 2 },
  ],
  certificates: [
    { key: 'name', label: '证书名称', placeholder: 'CET-6' },
    { key: 'issuer', label: '颁发机构', placeholder: '教育部考试中心' },
    { key: 'date', label: '日期', placeholder: '2023-06' },
    { key: 'score', label: '成绩', placeholder: '580' },
  ],
}

const moduleLabels: Record<string, string> = {
  basicInfo: '基础信息',
  skills: '专业技能',
  projects: '项目经历',
  evaluation: '自我评价',
  campus: '校园经历',
  internship: '实习经历',
  awards: '获奖',
  certificates: '证书',
}

const addItemLabels: Record<string, string> = {
  projects: '+ 添加项目',
  campus: '+ 添加校园经历',
  internship: '+ 添加实习经历',
  awards: '+ 添加获奖',
  certificates: '+ 添加证书',
}

// ============================================================
// 计算属性
// ============================================================

const currentModuleType = computed(() => props.modules[props.currentIndex] || '')
const currentModuleLabel = computed(() => moduleLabels[currentModuleType.value] || '')
const currentFields = computed(() => moduleFieldConfigs[currentModuleType.value] || [])
const isLastModule = computed(() => props.currentIndex === props.modules.length - 1)
const progressPercent = computed(() => ((props.currentIndex + 1) / props.modules.length) * 100)
const addItemLabel = computed(() => addItemLabels[currentModuleType.value] || '+ 添加')

// 本地表单数据（深拷贝避免直接修改 props）
const localData = ref<any>({})

// AI 润色状态
const isPolishing = ref(false)
const polishResult = ref<{ original: string; polished: string } | null>(null)

// 初始化本地数据
function initLocalData() {
  const moduleType = currentModuleType.value
  const source = props.formData[moduleType]
  if (source) {
    // 深拷贝
    localData.value = JSON.parse(JSON.stringify(source))
    // 为 skills 类别的 skillsText 做转换
    if (moduleType === 'skills' && localData.value.categories) {
      for (const cat of localData.value.categories) {
        if (!cat.skillsText && cat.skills) {
          cat.skillsText = Array.isArray(cat.skills) ? cat.skills.join(', ') : ''
        }
      }
    }
  } else {
    localData.value = {}
  }
  polishResult.value = null
}

// 监听当前模块变化
watch(() => [props.currentIndex, props.modules], () => {
  initLocalData()
}, { immediate: true })

// ============================================================
// 列表操作
// ============================================================

function addItem() {
  if (!localData.value.items) return
  const template = getItemTemplate(currentModuleType.value)
  localData.value.items.push({ ...template })
}

function removeItem(idx: number) {
  if (!localData.value.items) return
  localData.value.items.splice(idx, 1)
}

function addSkillCategory() {
  if (!localData.value.categories) return
  localData.value.categories.push({ category: '', skillsText: '' })
}

function removeSkillCategory(idx: number) {
  if (!localData.value.categories) return
  localData.value.categories.splice(idx, 1)
}

function getItemTemplate(moduleType: string): Record<string, string> {
  switch (moduleType) {
    case 'projects':
      return { name: '', role: '', startDate: '', endDate: '', techStack: '', description: '', highlights: '' }
    case 'campus':
      return { name: '', role: '', startDate: '', endDate: '', description: '', highlights: '' }
    case 'internship':
      return { company: '', position: '', startDate: '', endDate: '', description: '', highlights: '' }
    case 'awards':
      return { name: '', level: '', date: '', description: '' }
    case 'certificates':
      return { name: '', issuer: '', date: '', score: '' }
    default:
      return {}
  }
}

function getItemTitle(_item: any, moduleType: string): string {
  const titles: Record<string, string> = {
    projects: '项目',
    campus: '活动',
    internship: '实习',
    awards: '奖项',
    certificates: '证书',
  }
  return titles[moduleType] || '项目'
}

// ============================================================
// 数据同步
// ============================================================

function emitData() {
  const moduleType = currentModuleType.value
  const data = JSON.parse(JSON.stringify(localData.value))

  // skills: 将 skillsText 转回 skills 数组
  if (moduleType === 'skills' && data.categories) {
    for (const cat of data.categories) {
      if (cat.skillsText) {
        cat.skills = cat.skillsText.split(/[,，、]/).map((s: string) => s.trim()).filter(Boolean)
      } else {
        cat.skills = []
      }
    }
  }

  emit('update:formData', moduleType, data)
}

// ============================================================
// 下一步
// ============================================================

function handleNext() {
  emitData()
  if (isLastModule.value) {
    emit('finish')
  } else {
    emit('next')
  }
}

// ============================================================
// AI 润色
// ============================================================

const POLISH_SYSTEM_PROMPT = '你是专业的简历优化助手。请优化以下{模块名}内容，使其更加专业和有吸引力。只返回优化后的内容，不要添加解释或格式标记。'

function serializeFormData(): string {
  const moduleType = currentModuleType.value
  const data = localData.value

  // 将表单数据序列化为可读文本
  if (moduleType === 'basicInfo') {
    return Object.entries(data)
      .filter(([, v]) => typeof v === 'string' && v.trim())
      .map(([k, v]) => `${moduleFieldConfigs.basicInfo.find(f => f.key === k)?.label || k}: ${v}`)
      .join('\n')
  }

  if (moduleType === 'evaluation') {
    return data.content || ''
  }

  if (moduleType === 'skills') {
    return (data.categories || [])
      .filter((cat: any) => cat.category || cat.skillsText)
      .map((cat: any) => `${cat.category}: ${cat.skillsText}`)
      .join('\n')
  }

  // 列表类模块
  if (data.items) {
    const fields = currentFields.value
    return data.items
      .map((item: any, idx: number) => {
        const header = getItemTitle(item, moduleType) + ` ${idx + 1}`
        const body = fields
          .map(f => `${f.label}: ${item[f.key] || ''}`)
          .filter(line => !line.endsWith(': '))
          .join('\n')
        return header + '\n' + body
      })
      .join('\n\n')
  }

  return ''
}

async function handlePolish() {
  const content = serializeFormData()
  if (!content.trim()) return

  isPolishing.value = true
  polishResult.value = null

  try {
    const prompt = POLISH_SYSTEM_PROMPT.replace('{模块名}', currentModuleLabel.value)
      + '\n\n' + content

    const operationId = `polish-${currentModuleType.value}-${crypto.randomUUID()}`
    const result = await sendMessageSync(operationId, prompt, props.jobTarget, false)

    if (result) {
      polishResult.value = {
        original: content,
        polished: result.trim(),
      }
    }
  } catch (err) {
    console.error('AI 润色失败:', err)
  } finally {
    isPolishing.value = false
  }
}

async function handleRePolish() {
  if (polishResult.value) {
    polishResult.value = null
  }
  await handlePolish()
}

function acceptPolish() {
  if (!polishResult.value) return

  const polished = polishResult.value.polished
  const moduleType = currentModuleType.value

  // 将润色结果回填到表单
  // 对于 evaluation，直接设置 content
  if (moduleType === 'evaluation') {
    localData.value.content = polished
  }
  // 对于 basicInfo，尝试逐行解析回填
  else if (moduleType === 'basicInfo') {
    const lines = polished.split('\n').filter(Boolean)
    for (const line of lines) {
      const colonIdx = line.indexOf(':') !== -1 ? line.indexOf(':') : line.indexOf('：')
      if (colonIdx > 0) {
        const label = line.substring(0, colonIdx).trim()
        const value = line.substring(colonIdx + 1).trim()
        const field = currentFields.value.find(f => f.label === label)
        if (field) {
          localData.value[field.key] = value
        }
      }
    }
  }

  polishResult.value = null
}
</script>

<style scoped>
.step-form {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* 进度头 */
.form-header {
  padding: 12px 20px;
  border-bottom: 1px solid #3c3c3c;
  display: flex;
  align-items: center;
  gap: 12px;
  background: #252526;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: transparent;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  color: #8b949e;
  font-size: 12px;
  cursor: pointer;
  transition: color 0.15s, border-color 0.15s;
}

.back-btn:hover {
  color: #e0e0e0;
  border-color: #4c4c4c;
}

.form-progress {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.progress-text {
  font-size: 13px;
  font-weight: 500;
  color: #e0e0e0;
}

.progress-bar {
  height: 3px;
  background: #3c3c3c;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #0078d4;
  border-radius: 2px;
  transition: width 0.3s ease;
}

/* 表单主体 */
.form-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 字段组 */
.field-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.field-label {
  font-size: 12px;
  font-weight: 500;
  color: #8b949e;
}

.field-input {
  padding: 8px 10px;
  background: #2d2d2d;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  color: #e0e0e0;
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
}

.field-input:focus {
  border-color: #0078d4;
}

.field-input::placeholder {
  color: #5a5a5a;
}

.field-textarea {
  padding: 8px 10px;
  background: #2d2d2d;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  color: #e0e0e0;
  font-size: 13px;
  font-family: inherit;
  line-height: 1.5;
  outline: none;
  resize: vertical;
  transition: border-color 0.15s;
}

.field-textarea:focus {
  border-color: #0078d4;
}

.field-textarea::placeholder {
  color: #5a5a5a;
}

/* 列表项卡片 */
.items-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.list-item-card {
  background: #2a2a2a;
  border: 1px solid #3c3c3c;
  border-radius: 6px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.item-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.item-card-title {
  font-size: 13px;
  font-weight: 600;
  color: #e0e0e0;
}

.remove-btn {
  width: 22px;
  height: 22px;
  border: none;
  background: transparent;
  color: #8b949e;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background-color 0.15s, color 0.15s;
}

.remove-btn:hover {
  background: rgba(220, 38, 38, 0.15);
  color: #f87171;
}

.add-btn {
  padding: 8px;
  background: transparent;
  border: 1px dashed #3c3c3c;
  border-radius: 4px;
  color: #8b949e;
  font-size: 13px;
  cursor: pointer;
  transition: border-color 0.15s, color 0.15s;
}

.add-btn:hover {
  border-color: #0078d4;
  color: #58a6ff;
}

/* 技能类别 */
.skill-category {
  background: #2a2a2a;
  border: 1px solid #3c3c3c;
  border-radius: 6px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.category-label {
  font-size: 13px;
  font-weight: 600;
  color: #e0e0e0;
}

/* AI 润色结果 */
.polish-section {
  border-top: 1px solid #3c3c3c;
  padding: 12px 20px;
  background: #252526;
}

.polish-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.polish-title {
  font-size: 13px;
  font-weight: 600;
  color: #58a6ff;
}

.polish-close {
  width: 22px;
  height: 22px;
  border: none;
  background: transparent;
  color: #8b949e;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.polish-close:hover {
  background: rgba(220, 38, 38, 0.15);
  color: #f87171;
}

.polish-compare {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.polish-column {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.polish-label {
  font-size: 11px;
  color: #8b949e;
  font-weight: 500;
}

.polish-content {
  margin: 0;
  padding: 8px;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.4;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 120px;
  overflow-y: auto;
}

.polish-content.original {
  background: rgba(220, 38, 38, 0.08);
  border: 1px solid rgba(220, 38, 38, 0.2);
  color: #e0e0e0;
}

.polish-content.polished {
  background: rgba(34, 197, 94, 0.08);
  border: 1px solid rgba(34, 197, 94, 0.2);
  color: #e0e0e0;
}

.polish-actions {
  display: flex;
  gap: 6px;
  justify-content: flex-end;
}

/* 底部操作 */
.form-footer {
  padding: 12px 20px;
  border-top: 1px solid #3c3c3c;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #252526;
}

.btn-polish {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: transparent;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  color: #58a6ff;
  font-size: 13px;
  cursor: pointer;
  transition: border-color 0.15s, background-color 0.15s;
}

.btn-polish:hover:not(:disabled) {
  border-color: #58a6ff;
  background: rgba(88, 166, 255, 0.08);
}

.btn-polish:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.footer-right {
  display: flex;
  gap: 8px;
}

.btn-primary {
  padding: 8px 16px;
  background: #0e639c;
  color: #ffffff;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s;
}

.btn-primary:hover {
  background: #1177bb;
}

.btn-secondary {
  padding: 8px 16px;
  background: #3c3c3c;
  color: #cccccc;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: background-color 0.15s;
}

.btn-secondary:hover {
  background: #4a4a4a;
}

.btn-sm {
  padding: 4px 10px;
  font-size: 12px;
}

/* 加载动画 */
.spinner {
  display: inline-block;
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
</style>
