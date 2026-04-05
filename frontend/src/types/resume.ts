// Resume TypeScript type definitions

export interface Resume {
  id: string
  title: string
  basicInfo: BasicInfo
  modules: Module[]
  templateId: string
  customCss: string
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
}
