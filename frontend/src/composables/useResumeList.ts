import { ref, computed } from 'vue'
import type { ResumeListItem } from '../types/resume'
import {
  ListResumes,
  RenameResume,
  DuplicateResume as DuplicateResumeAPI,
  DeleteResume
} from '../wailsjs/wailsjs/go/main/App'

/**
 * 简历列表状态管理 composable
 * 提供响应式简历列表、搜索过滤、排序、CRUD 操作
 */
export function useResumeList() {
  const resumes = ref<ResumeListItem[]>([])
  const loading = ref(false)
  const searchQuery = ref('')
  const sortOrder = ref<'newest' | 'oldest'>('newest')
  const error = ref<string | null>(null)

  // 根据搜索词过滤并按排序规则排序
  const filteredResumes = computed(() => {
    let list = [...resumes.value]
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      list = list.filter(r => r.title.toLowerCase().includes(query))
    }
    return list.sort((a, b) => {
      const dateA = new Date(a.updatedAt).getTime()
      const dateB = new Date(b.updatedAt).getTime()
      return sortOrder.value === 'newest' ? dateB - dateA : dateA - dateB
    })
  })

  // 加载简历列表
  async function fetchResumes(): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const list = await ListResumes()
      resumes.value = (list || []).map(item => ({
        id: item.id,
        title: item.title,
        updatedAt: item.updatedAt
      }))
    } catch (err) {
      error.value = String(err)
      console.error('加载简历列表失败:', err)
    } finally {
      loading.value = false
    }
  }

  // 重命名简历
  async function renameResume(id: string, title: string): Promise<void> {
    try {
      await RenameResume(id, title)
      // 更新本地列表中的标题
      const item = resumes.value.find(r => r.id === id)
      if (item) {
        item.title = title
      }
    } catch (err) {
      error.value = String(err)
      console.error('重命名失败:', err)
      throw err
    }
  }

  // 复制简历
  async function duplicateResume(id: string): Promise<void> {
    try {
      await DuplicateResumeAPI(id)
      // 重新加载列表以获取新的副本
      await fetchResumes()
    } catch (err) {
      error.value = String(err)
      console.error('复制简历失败:', err)
      throw err
    }
  }

  // 删除简历（软删除）
  async function deleteResume(id: string): Promise<void> {
    try {
      await DeleteResume(id)
      // 从本地列表中移除
      resumes.value = resumes.value.filter(r => r.id !== id)
    } catch (err) {
      error.value = String(err)
      console.error('删除简历失败:', err)
      throw err
    }
  }

  return {
    resumes,
    loading,
    searchQuery,
    sortOrder,
    filteredResumes,
    error,
    fetchResumes,
    renameResume,
    duplicateResume,
    deleteResume
  }
}
