<template>
  <Transition name="wizard-sidebar">
    <div v-if="visible" class="wizard-sidebar">
      <!-- 顶部标题栏 -->
      <div class="wizard-header">
        <span class="wizard-title">创建简历向导</span>
        <button class="close-btn" @click="handleClose" title="关闭">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>

      <!-- 步骤指示器 -->
      <div class="step-indicator">
        <div
          v-for="(step, idx) in stepLabels"
          :key="idx"
          :class="['step-dot', {
            active: currentStepIndex === idx,
            completed: currentStepIndex > idx
          }]"
          @click="idx < currentStepIndex && goToStep(idx)"
        >
          <span class="step-number">{{ idx < currentStepIndex ? '✓' : idx + 1 }}</span>
          <span class="step-label">{{ step }}</span>
        </div>
      </div>

      <!-- 内容区 -->
      <div class="wizard-content">
        <!-- 步骤一：模块勾选 -->
        <div v-if="currentStep === 'select'" class="step-content">
          <WizardModuleSelect
            v-model:selected="selectedModules"
            @start-fill="startFill"
          />
        </div>

        <!-- 步骤二：逐步填写 -->
        <div v-if="currentStep === 'fill'" class="step-content">
          <WizardStepForm
            :modules="selectedModules"
            :current-index="currentModuleIndex"
            :form-data="formData"
            :job-target="jobTarget"
            @update:form-data="updateFormData"
            @next="nextModule"
            @skip="skipModule"
            @back="prevModule"
            @finish="goToGenerate"
          />
        </div>

        <!-- 步骤三：生成简历 -->
        <div v-if="currentStep === 'generate'" class="step-content">
          <WizardGenerate
            :selected-modules="selectedModules"
            :form-data="formData"
            :is-generating="isGenerating"
            @back="currentStep = 'fill'"
            @generate="handleGenerate"
          />
        </div>
      </div>

      <!-- 中途退出确认框 -->
      <Teleport to="body">
        <div v-if="showCloseConfirm" class="confirm-overlay" @click.self="showCloseConfirm = false">
          <div class="confirm-dialog">
            <p class="confirm-text">退出将保存已填写内容为草稿，是否继续？</p>
            <div class="confirm-actions">
              <button class="btn-secondary" @click="showCloseConfirm = false">继续填写</button>
              <button class="btn-primary" @click="confirmClose">保存草稿并退出</button>
            </div>
          </div>
        </div>
      </Teleport>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { UpdateResume } from '../wailsjs/wailsjs/go/main/App'
import { sendMessageSync } from '../services/ai'
import WizardModuleSelect from './wizard/WizardModuleSelect.vue'
import WizardStepForm from './wizard/WizardStepForm.vue'
import WizardGenerate from './wizard/WizardGenerate.vue'

// ============================================================
// Props & Emits
// ============================================================

interface Props {
  visible: boolean
  resumeId: string
  jobTarget: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'complete'): void
  (e: 'save-draft'): void
}>()

// ============================================================
// 向导状态
// ============================================================

type WizardStep = 'select' | 'fill' | 'generate'

const currentStep = ref<WizardStep>('select')
const selectedModules = ref<string[]>([])
const currentModuleIndex = ref(0)
const formData = ref<Record<string, any>>({})
const showCloseConfirm = ref(false)
const isGenerating = ref(false)

// 步骤标签
const stepLabels = ['选择模块', '填写内容', '生成简历']

// 当前步骤索引
const currentStepIndex = computed(() => {
  switch (currentStep.value) {
    case 'select': return 0
    case 'fill': return 1
    case 'generate': return 2
    default: return 0
  }
})

// ============================================================
// 步骤切换
// ============================================================

function goToStep(idx: number) {
  if (idx < currentStepIndex.value) {
    const steps: WizardStep[] = ['select', 'fill', 'generate']
    currentStep.value = steps[idx]
  }
}

function startFill() {
  currentStep.value = 'fill'
  currentModuleIndex.value = 0
  // 初始化每个模块的表单数据
  for (const mod of selectedModules.value) {
    if (!formData.value[mod]) {
      formData.value[mod] = getInitialFormData(mod)
    }
  }
}

function updateFormData(moduleType: string, data: any) {
  formData.value[moduleType] = data
}

function nextModule() {
  if (currentModuleIndex.value < selectedModules.value.length - 1) {
    currentModuleIndex.value++
  } else {
    // 所有模块填写完成，进入生成步骤
    currentStep.value = 'generate'
  }
}

function skipModule() {
  // 跳过当前模块（保留空数据）
  nextModule()
}

function prevModule() {
  if (currentModuleIndex.value > 0) {
    currentModuleIndex.value--
  } else {
    // 返回模块勾选步骤
    currentStep.value = 'select'
  }
}

function goToGenerate() {
  currentStep.value = 'generate'
}

// ============================================================
// 生成简历
// ============================================================

async function handleGenerate() {
  isGenerating.value = true
  try {
    // 将 formData 组装为后端期望的 modules JSON 格式
    const modules = buildModulesJSON(selectedModules.value, formData.value)

    // 调用 UpdateResume 保存，后端自动执行 JSON->Markdown 转换
    await UpdateResume(props.resumeId, JSON.stringify(modules))

    // 触发完成事件
    emit('complete')
  } catch (err) {
    console.error('生成简历失败:', err)
  } finally {
    isGenerating.value = false
  }
}

// ============================================================
// 关闭处理（中途退出）
// ============================================================

function handleClose() {
  // 如果已经填写了内容，弹出确认框
  const hasData = Object.values(formData.value).some(
    (data) => data && Object.values(data).some((v) => {
      if (Array.isArray(v)) return v.length > 0
      if (typeof v === 'string') return v.trim() !== ''
      return v != null && v !== ''
    })
  )

  if (hasData && currentStep.value !== 'select') {
    showCloseConfirm.value = true
  } else {
    emit('close')
  }
}

async function confirmClose() {
  showCloseConfirm.value = false

  // 保存已填写的草稿
  const hasFilledModules = selectedModules.value.length > 0 &&
    Object.keys(formData.value).length > 0

  if (hasFilledModules) {
    try {
      const modules = buildModulesJSON(selectedModules.value, formData.value)
      await UpdateResume(props.resumeId, JSON.stringify(modules))
    } catch (err) {
      console.error('保存草稿失败:', err)
    }
  }

  emit('save-draft')
}

// ============================================================
// 数据构建
// ============================================================

/**
 * 获取模块的初始表单数据
 */
function getInitialFormData(moduleType: string): any {
  switch (moduleType) {
    case 'basicInfo':
      return { name: '', phone: '', email: '', github: '', website: '', address: '', summary: '' }
    case 'skills':
      return { categories: [{ category: '', skills: [] }] }
    case 'projects':
      return { items: [{ name: '', role: '', startDate: '', endDate: '', description: '', techStack: '', highlights: '' }] }
    case 'evaluation':
      return { content: '' }
    case 'campus':
      return { items: [{ name: '', role: '', startDate: '', endDate: '', description: '', highlights: '' }] }
    case 'internship':
      return { items: [{ company: '', position: '', startDate: '', endDate: '', description: '', highlights: '' }] }
    case 'awards':
      return { items: [{ name: '', level: '', date: '', description: '' }] }
    case 'certificates':
      return { items: [{ name: '', issuer: '', date: '', score: '' }] }
    default:
      return {}
  }
}

/**
 * 将表单数据构建为后端期望的 modules JSON 数组
 */
function buildModulesJSON(selected: string[], data: Record<string, any>): any[] {
  const modules: any[] = []
  let order = 0

  for (const moduleType of selected) {
    const moduleData = data[moduleType]
    if (!moduleData) continue

    const moduleTitle = getModuleTitle(moduleType)

    switch (moduleType) {
      case 'basicInfo': {
        // basicInfo 作为特殊模块，items 是 BasicInfo JSON
        modules.push({
          type: 'basicInfo',
          title: '基本信息',
          order: order++,
          visible: true,
          items: {
            name: moduleData.name || '',
            phone: moduleData.phone || '',
            email: moduleData.email || '',
            avatar: '',
            website: moduleData.website || '',
            github: moduleData.github || '',
            address: moduleData.address || '',
            summary: moduleData.summary || '',
          },
        })
        break
      }

      case 'skills': {
        const items = (moduleData.categories || [])
          .filter((cat: any) => cat.category || (cat.skills && cat.skills.length > 0))
          .map((cat: any) => ({
            category: cat.category || '',
            skills: Array.isArray(cat.skills) ? cat.skills : (cat.skills ? cat.skills.split(/[,，、]/).map((s: string) => s.trim()).filter(Boolean) : []),
          }))
        if (items.length > 0) {
          modules.push({ type: 'skills', title: moduleTitle, order: order++, visible: true, items })
        }
        break
      }

      case 'projects': {
        const items = (moduleData.items || [])
          .filter((item: any) => item.name)
          .map((item: any) => ({
            name: item.name,
            role: item.role || '',
            startDate: item.startDate || '',
            endDate: item.endDate || '',
            description: item.description || '',
            highlights: item.highlights ? item.highlights.split('\n').filter((h: string) => h.trim()) : [],
            techStack: item.techStack ? item.techStack.split(/[,，、]/).map((s: string) => s.trim()).filter(Boolean) : [],
          }))
        if (items.length > 0) {
          modules.push({ type: 'projects', title: moduleTitle, order: order++, visible: true, items })
        }
        break
      }

      case 'evaluation': {
        if (moduleData.content) {
          modules.push({
            type: 'evaluation',
            title: moduleTitle,
            order: order++,
            visible: true,
            items: [{ content: moduleData.content }],
          })
        }
        break
      }

      case 'campus': {
        const items = (moduleData.items || [])
          .filter((item: any) => item.name)
          .map((item: any) => ({
            name: item.name,
            role: item.role || '',
            startDate: item.startDate || '',
            endDate: item.endDate || '',
            description: item.description || '',
            highlights: item.highlights ? item.highlights.split('\n').filter((h: string) => h.trim()) : [],
          }))
        if (items.length > 0) {
          modules.push({ type: 'campus', title: moduleTitle, order: order++, visible: true, items })
        }
        break
      }

      case 'internship': {
        const items = (moduleData.items || [])
          .filter((item: any) => item.company)
          .map((item: any) => ({
            company: item.company,
            position: item.position || '',
            startDate: item.startDate || '',
            endDate: item.endDate || '',
            description: item.description || '',
            highlights: item.highlights ? item.highlights.split('\n').filter((h: string) => h.trim()) : [],
          }))
        if (items.length > 0) {
          modules.push({ type: 'internship', title: moduleTitle, order: order++, visible: true, items })
        }
        break
      }

      case 'awards': {
        const items = (moduleData.items || [])
          .filter((item: any) => item.name)
          .map((item: any) => ({
            name: item.name,
            level: item.level || '',
            date: item.date || '',
            description: item.description || '',
          }))
        if (items.length > 0) {
          modules.push({ type: 'awards', title: moduleTitle, order: order++, visible: true, items })
        }
        break
      }

      case 'certificates': {
        const items = (moduleData.items || [])
          .filter((item: any) => item.name)
          .map((item: any) => ({
            name: item.name,
            issuer: item.issuer || '',
            date: item.date || '',
            score: item.score || '',
          }))
        if (items.length > 0) {
          modules.push({ type: 'certificates', title: moduleTitle, order: order++, visible: true, items })
        }
        break
      }
    }
  }

  return modules
}

/**
 * 获取模块显示标题
 */
function getModuleTitle(moduleType: string): string {
  const titles: Record<string, string> = {
    basicInfo: '基本信息',
    skills: '专业技能',
    projects: '项目经历',
    evaluation: '自我评价',
    campus: '校园经历',
    internship: '实习经历',
    awards: '获奖',
    certificates: '证书',
  }
  return titles[moduleType] || moduleType
}

// ============================================================
// 重置
// ============================================================

watch(() => props.visible, (v) => {
  if (v) {
    // 打开时重置状态
    currentStep.value = 'select'
    selectedModules.value = []
    currentModuleIndex.value = 0
    formData.value = {}
    isGenerating.value = false
  }
})
</script>

<style scoped>
.wizard-sidebar {
  position: fixed;
  top: 0;
  right: 0;
  width: 420px;
  height: 100vh;
  background: #1e1e1e;
  border-left: 1px solid #3c3c3c;
  z-index: 100;
  display: flex;
  flex-direction: column;
  box-shadow: -4px 0 20px rgba(0, 0, 0, 0.3);
}

/* 滑入动画 */
.wizard-sidebar-enter-active,
.wizard-sidebar-leave-active {
  transition: transform 300ms ease-out;
}

.wizard-sidebar-enter-from,
.wizard-sidebar-leave-to {
  transform: translateX(100%);
}

/* 顶部标题栏 */
.wizard-header {
  height: 48px;
  min-height: 48px;
  padding: 0 16px;
  background: #252526;
  border-bottom: 1px solid #3c3c3c;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.wizard-title {
  font-size: 14px;
  font-weight: 600;
  color: #e0e0e0;
}

.close-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #8b949e;
  transition: background-color 0.15s, color 0.15s;
}

.close-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}

/* 步骤指示器 */
.step-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px 16px;
  gap: 4px;
  background: #252526;
  border-bottom: 1px solid #3c3c3c;
}

.step-dot {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: #8b949e;
  cursor: default;
  transition: color 0.2s, background-color 0.2s;
}

.step-dot.completed {
  color: #4fc3f7;
  cursor: pointer;
}

.step-dot.completed:hover {
  background: rgba(79, 195, 247, 0.1);
}

.step-dot.active {
  color: #e0e0e0;
  font-weight: 600;
}

.step-number {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  background: #3c3c3c;
  color: #8b949e;
}

.step-dot.active .step-number {
  background: #0078d4;
  color: #ffffff;
}

.step-dot.completed .step-number {
  background: #2ea043;
  color: #ffffff;
}

.step-label {
  font-size: 12px;
}

/* 内容区 */
.wizard-content {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.step-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 确认弹窗 */
.confirm-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.confirm-dialog {
  background: #2d2d2d;
  border: 1px solid #3c3c3c;
  border-radius: 8px;
  padding: 24px;
  max-width: 360px;
  width: 90%;
}

.confirm-text {
  font-size: 14px;
  color: #e0e0e0;
  margin: 0 0 20px 0;
  line-height: 1.5;
}

.confirm-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

/* 通用按钮 */
.btn-primary {
  padding: 8px 16px;
  background: #0e639c;
  color: #ffffff;
  border: none;
  border-radius: 4px;
  font-size: 13px;
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
</style>
