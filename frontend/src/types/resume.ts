// Resume TypeScript type definitions

export interface Resume {
  id: string
  title: string
  basicInfo: BasicInfo
  modules: Module[]
  templateId: string
  customCss: string
  markdownContent: string
  jobTarget: string
  createdAt: string
  updatedAt: string
}

export interface BasicInfo {
  name: string
  phone: string
  email: string
  avatar: string
  website: string
  github: string
  address: string
  summary: string
}

export interface Module {
  type: string
  title: string
  order: number
  items: unknown
  visible: boolean
}

export interface EducationItem {
  school: string
  degree: string
  major: string
  startDate: string
  endDate: string
  gpa: string
  highlights: string[]
}

export interface SkillItem {
  category: string
  skills: string[]
}

export interface ProjectItem {
  name: string
  role: string
  startDate: string
  endDate: string
  description: string
  highlights: string[]
  techStack: string[]
}

export interface InternshipItem {
  company: string
  position: string
  startDate: string
  endDate: string
  description: string
  highlights: string[]
}

export interface CampusItem {
  name: string
  role: string
  startDate: string
  endDate: string
  description: string
  highlights: string[]
}

export interface AwardItem {
  name: string
  level: string
  date: string
  description: string
}

export interface CertificateItem {
  name: string
  issuer: string
  date: string
  score: string
}

export interface EvaluationItem {
  content: string
}

export interface ResumeListItem {
  id: string
  title: string
  updatedAt: string
  deletedAt?: string  // 回收站列表需要
}

/**
 * Parsed resume data structure from AI resume parsing.
 * This is the expected JSON output format from the resume parser AI.
 */
export interface ParsedResume {
  name: string | null
  email: string | null
  phone: string | null
  education: EducationItem[] | null
  experience: ExperienceItem[] | null
  projects: ProjectItem[] | null
  skills: string[] | null
  summary: string | null
}

/**
 * Experience item for parsed resume data.
 */
export interface ExperienceItem {
  company: string | null
  position: string | null
  startDate: string | null
  endDate: string | null
  description: string | null
}

/**
 * Validates that a parsed JSON object matches the expected ParsedResume schema.
 * Used to verify AI parser output before importing.
 */
export function validateParsedResume(data: unknown): data is ParsedResume {
  if (data === null || data === undefined || typeof data !== 'object') {
    return false
  }
  const r = data as Record<string, unknown>
  if (typeof r.name !== 'string' && r.name !== null) {
    return false
  }
  if (r.education !== null && !Array.isArray(r.education)) {
    return false
  }
  if (r.experience !== null && !Array.isArray(r.experience)) {
    return false
  }
  if (r.projects !== null && !Array.isArray(r.projects)) {
    return false
  }
  if (r.skills !== null && !Array.isArray(r.skills)) {
    return false
  }
  if (r.email !== undefined && typeof r.email !== 'string' && r.email !== null) {
    return false
  }
  if (r.phone !== undefined && typeof r.phone !== 'string' && r.phone !== null) {
    return false
  }
  if (r.summary !== undefined && typeof r.summary !== 'string' && r.summary !== null) {
    return false
  }
  return true
}
