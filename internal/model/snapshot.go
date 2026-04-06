package model

import "time"

// Snapshot represents a version snapshot of a resume
type Snapshot struct {
	ID              string    `json:"id"`
	ResumeID        string    `json:"resumeId"`
	Label           string    `json:"label"`
	Note            string    `json:"note"`
	TriggerType     string    `json:"triggerType"` // 'manual' | 'auto_pdf_export' | 'rollback'
	JSONData        string    `json:"jsonData"`
	MarkdownContent string    `json:"markdownContent"`
	TemplateID      string    `json:"templateId"`
	CustomCSS       string    `json:"customCss"`
	CreatedAt       time.Time `json:"createdAt"`
}

// CreateSnapshotRequest 用于创建快照的请求参数
type CreateSnapshotRequest struct {
	ResumeID    string `json:"resumeId"`
	Label       string `json:"label"`
	Note        string `json:"note"`
	TriggerType string `json:"triggerType"` // 'manual' | 'auto_pdf_export' | 'rollback'
}

// SnapshotListItem 用于列表展示（不含完整数据）
type SnapshotListItem struct {
	ID          string    `json:"id"`
	ResumeID    string    `json:"resumeId"`
	Label       string    `json:"label"`
	Note        string    `json:"note"`
	TriggerType string    `json:"triggerType"`
	CreatedAt   time.Time `json:"createdAt"`
}

// DiffResult 快照对比结果
type DiffResult struct {
	Snapshot1ID    string     `json:"snapshot1Id"`
	Snapshot2ID    string     `json:"snapshot2Id"`
	Snapshot1Label string     `json:"snapshot1Label"`
	Snapshot2Label string     `json:"snapshot2Label"`
	ContentDiff    string     `json:"contentDiff"` // 完整 diff patch 字符串
	Stats          DiffStats  `json:"stats"`
}

// DiffStats 对比统计
type DiffStats struct {
	AddedLines   int `json:"addedLines"`
	RemovedLines int `json:"removedLines"`
	ChangedLines int `json:"changedLines"`
}
