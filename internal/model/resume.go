package model

import (
	"encoding/json"
	"time"
)

// Resume represents a resume document
type Resume struct {
	ID              string       `json:"id"`
	Title           string       `json:"title"`
	BasicInfo       BasicInfo    `json:"basicInfo"`
	Modules         []Module     `json:"modules"`
	TemplateID      string       `json:"templateId"`
	CustomCSS       string       `json:"customCss"`
	MarkdownContent string       `json:"markdownContent"`
	JobTarget       string       `json:"jobTarget"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
}

// BasicInfo contains personal information
type BasicInfo struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Website  string `json:"website"`
	GitHub   string `json:"github"`
	Address  string `json:"address"`
	Summary  string `json:"summary"`
}

// Module represents a resume section
type Module struct {
	Type     string          `json:"type"`
	Title    string          `json:"title"`
	Order    int             `json:"order"`
	Items    json.RawMessage `json:"items"`
	Visible  bool            `json:"visible"`
}

// EducationItem represents an education entry
type EducationItem struct {
	School     string   `json:"school"`
	Degree     string   `json:"degree"`
	Major      string   `json:"major"`
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate"`
	GPA        string   `json:"gpa"`
	Highlights []string `json:"highlights"`
}

// SkillItem represents a skill entry
type SkillItem struct {
	Category string   `json:"category"`
	Skills   []string `json:"skills"`
}

// ProjectItem represents a project entry
type ProjectItem struct {
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Description string   `json:"description"`
	Highlights  []string `json:"highlights"`
	TechStack   []string `json:"techStack"`
}

// InternshipItem represents an internship entry
type InternshipItem struct {
	Company     string   `json:"company"`
	Position    string   `json:"position"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Description string   `json:"description"`
	Highlights  []string `json:"highlights"`
}

// CampusItem represents a campus activity entry
type CampusItem struct {
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Description string   `json:"description"`
	Highlights  []string `json:"highlights"`
}

// AwardItem represents an award entry
type AwardItem struct {
	Name        string `json:"name"`
	Level       string `json:"level"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

// CertificateItem represents a certificate entry
type CertificateItem struct {
	Name    string `json:"name"`
	Issuer  string `json:"issuer"`
	Date    string `json:"date"`
	Score   string `json:"score"`
}

// EvaluationItem represents a self-evaluation entry
type EvaluationItem struct {
	Content string `json:"content"`
}

// ResumeListItem represents a resume item in list view
type ResumeListItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
}
