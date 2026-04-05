/**
 * Resume JSON to Markdown conversion utility.
 * Converts structured ParsedResume data to the app's Markdown format.
 */

import type { ParsedResume, EducationItem, ProjectItem, ExperienceItem } from '../types/resume'

/**
 * Converts ParsedResume data into Markdown content.
 * Produces a well-formatted resume in Markdown format.
 */
export function jsonToMarkdown(data: ParsedResume): string {
  const parts: string[] = []

  // Basic Info Section
  const basicInfoParts: string[] = []
  if (data.name) basicInfoParts.push(data.name)
  if (data.email) basicInfoParts.push(data.email)
  if (data.phone) basicInfoParts.push(data.phone)

  if (basicInfoParts.length > 0) {
    parts.push(basicInfoParts.join(' | ') + '\n')
  }

  // Summary / 自我评价
  if (data.summary) {
    parts.push(`## 自我评价\n\n${data.summary}\n\n`)
  }

  // Education / 教育背景
  if (data.education && data.education.length > 0) {
    parts.push('## 教育背景\n\n')
    for (const edu of data.education) {
      const line = formatEducationItem(edu)
      parts.push(line + '\n')
    }
    parts.push('\n')
  }

  // Experience / 工作经历
  if (data.experience && data.experience.length > 0) {
    parts.push('## 工作经历\n\n')
    for (const exp of data.experience) {
      const section = formatExperienceItem(exp)
      parts.push(section + '\n')
    }
    parts.push('\n')
  }

  // Projects / 项目经历
  if (data.projects && data.projects.length > 0) {
    parts.push('## 项目经历\n\n')
    for (const proj of data.projects) {
      const section = formatProjectItem(proj)
      parts.push(section + '\n')
    }
    parts.push('\n')
  }

  // Skills / 技能证书
  if (data.skills && data.skills.length > 0) {
    parts.push('## 技能证书\n\n')
    parts.push(data.skills.map(s => `- ${s}`).join('\n') + '\n')
  }

  return parts.join('').trim() + '\n'
}

function formatEducationItem(edu: EducationItem): string {
  const parts: string[] = []
  if (edu.school) parts.push(edu.school)
  if (edu.major) parts.push(edu.major)
  if (edu.degree) parts.push(edu.degree)

  let line = parts.join(' | ')
  if (edu.startDate || edu.endDate) {
    const dateRange = [edu.startDate, edu.endDate].filter(Boolean).join(' - ')
    line += ` | ${dateRange}`
  }

  if (edu.gpa) {
    line += ` | GPA: ${edu.gpa}`
  }

  if (edu.highlights && edu.highlights.length > 0) {
    line += '\n' + edu.highlights.map(h => `  - ${h}`).join('\n')
  }

  return line
}

function formatExperienceItem(exp: ExperienceItem): string {
  const parts: string[] = []
  if (exp.company) parts.push(exp.company)
  if (exp.position) parts.push(`(${exp.position})`)

  let line = parts.join(' ')
  if (exp.startDate || exp.endDate) {
    const dateRange = [exp.startDate, exp.endDate].filter(Boolean).join(' - ')
    line += ` | ${dateRange}`
  }

  if (exp.description) {
    line += '\n' + exp.description
  }

  return line
}

function formatProjectItem(proj: ProjectItem): string {
  const parts: string[] = []
  if (proj.name) parts.push(proj.name)
  if (proj.role) parts.push(`(${proj.role})`)

  let line = parts.join(' ')
  if (proj.startDate || proj.endDate) {
    const dateRange = [proj.startDate, proj.endDate].filter(Boolean).join(' - ')
    line += ` | ${dateRange}`
  }

  if (proj.techStack && proj.techStack.length > 0) {
    line += `\n技术栈: ${proj.techStack.join(', ')}`
  }

  if (proj.description) {
    line += '\n' + proj.description
  }

  if (proj.highlights && proj.highlights.length > 0) {
    line += '\n' + proj.highlights.map(h => `- ${h}`).join('\n')
  }

  return line
}

/**
 * Extracts and merges contact information from parsed resume into basicInfo format.
 * Returns a partial basicInfo object that can be merged with existing data.
 */
export function extractBasicInfo(data: ParsedResume): Record<string, string> {
  const info: Record<string, string> = {}
  if (data.name) info.name = data.name
  if (data.email) info.email = data.email
  if (data.phone) info.phone = data.phone
  if (data.summary) info.summary = data.summary
  return info
}
