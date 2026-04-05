package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"Darvin-Resume/internal/database"
	"Darvin-Resume/internal/model"

	_ "modernc.org/sqlite"
)

var testDBPath = filepath.Join(os.TempDir(), "test_resume.db")

func setupTestDB(t *testing.T) func() {
	// Remove existing test database
	os.Remove(testDBPath)

	// Create a new test database
	db, err := sql.Open("sqlite", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables
	schema := `
	CREATE TABLE IF NOT EXISTS resumes (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL DEFAULT '',
		json_data TEXT NOT NULL,
		markdown_content TEXT NOT NULL DEFAULT '',
		template_id TEXT NOT NULL DEFAULT 'default',
		custom_css TEXT NOT NULL DEFAULT '',
		module_order TEXT NOT NULL DEFAULT '[]',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
		deleted_at DATETIME
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Set the global DB
	database.DB = db

	// Return cleanup function
	return func() {
		db.Close()
		os.Remove(testDBPath)
	}
}

func TestCreateResume(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	svc := NewResumeService()
	ctx := context.Background()

	// Create a resume with basic info
	resume := &model.Resume{
		Title: "测试简历",
		BasicInfo: model.BasicInfo{
			Name:  "张三",
			Phone: "13800138000",
			Email: "zhangsan@example.com",
		},
	}

	err := svc.Create(ctx, resume)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify UUID was generated
	if resume.ID == "" {
		t.Error("UUID should be generated")
	}

	// Verify timestamps
	if resume.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if resume.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}

	// Verify resume can be retrieved
	retrieved, err := svc.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrieved.Title != "测试简历" {
		t.Errorf("Expected title '测试简历', got '%s'", retrieved.Title)
	}
	if retrieved.BasicInfo.Name != "张三" {
		t.Errorf("Expected name '张三', got '%s'", retrieved.BasicInfo.Name)
	}
}

func TestGetResume(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	svc := NewResumeService()
	ctx := context.Background()

	// Create a resume
	resume := &model.Resume{
		Title: "测试简历2",
		Modules: []model.Module{
			{
				Type:    "education",
				Title:   "教育经历",
				Order:   1,
				Visible: true,
				Items:   json.RawMessage(`[{"school":"清华大学","major":"计算机科学","degree":"本科","startDate":"2018-09","endDate":"2022-06"}]`),
			},
		},
	}

	err := svc.Create(ctx, resume)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Retrieve the resume
	retrieved, err := svc.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	// Verify JSON parsing
	if len(retrieved.Modules) != 1 {
		t.Fatalf("Expected 1 module, got %d", len(retrieved.Modules))
	}
	if retrieved.Modules[0].Type != "education" {
		t.Errorf("Expected module type 'education', got '%s'", retrieved.Modules[0].Type)
	}

	// Verify markdown was generated
	if retrieved.MarkdownContent == "" {
		t.Error("MarkdownContent should be generated")
	}
}

func TestListResumes(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	svc := NewResumeService()
	ctx := context.Background()

	// Create 3 resumes with different timestamps
	for i := 1; i <= 3; i++ {
		resume := &model.Resume{
			Title: "简历",
		}
		if err := svc.Create(ctx, resume); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		// Add a small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// List resumes
	items, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 resumes, got %d", len(items))
	}

	// Verify sorted by updated_at descending (newest first)
	if items[0].UpdatedAt.Before(items[2].UpdatedAt) {
		t.Error("List should be sorted by updated_at descending")
	}
}

func TestUpdateResume(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	svc := NewResumeService()
	ctx := context.Background()

	// Create a resume
	resume := &model.Resume{
		Title: "原始标题",
		BasicInfo: model.BasicInfo{
			Name: "原始姓名",
		},
	}

	if err := svc.Create(ctx, resume); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update the resume
	resume.Title = "新标题"
	resume.BasicInfo.Name = "新姓名"
	resume.Modules = []model.Module{
		{
			Type:    "education",
			Title:   "教育经历",
			Order:   1,
			Visible: true,
			Items:   json.RawMessage(`[{"school":"北京大学"}]`),
		},
	}

	if err := svc.Update(ctx, resume); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify the update
	retrieved, err := svc.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrieved.Title != "新标题" {
		t.Errorf("Expected title '新标题', got '%s'", retrieved.Title)
	}
	if retrieved.BasicInfo.Name != "新姓名" {
		t.Errorf("Expected name '新姓名', got '%s'", retrieved.BasicInfo.Name)
	}

	// Verify markdown was regenerated
	if retrieved.MarkdownContent == "" {
		t.Error("MarkdownContent should be regenerated after update")
	}
}

func TestDeleteResume(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	svc := NewResumeService()
	ctx := context.Background()

	// Create a resume
	resume := &model.Resume{
		Title: "待删除简历",
	}

	if err := svc.Create(ctx, resume); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete the resume
	if err := svc.Delete(ctx, resume.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify GetByID returns error
	_, err := svc.GetByID(ctx, resume.ID)
	if err != ErrResumeNotFound {
		t.Errorf("Expected ErrResumeNotFound, got %v", err)
	}

	// Verify List does not include deleted resume
	items, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("Expected 0 resumes after delete, got %d", len(items))
	}
}

func TestUpdateJSON(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	svc := NewResumeService()
	ctx := context.Background()

	// Create a resume
	resume := &model.Resume{
		Title: "JSON更新测试",
	}

	if err := svc.Create(ctx, resume); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update JSON data
	newJSON := `[
		{"type":"basicInfo","title":"基本信息","order":0,"visible":true,"items":{"name":"李四","phone":"13900139000"}},
		{"type":"education","title":"教育经历","order":1,"visible":true,"items":[{"school":"MIT","major":"CS","degree":"硕士"}]}
	]`

	if err := svc.UpdateJSON(ctx, resume.ID, newJSON); err != nil {
		t.Fatalf("UpdateJSON failed: %v", err)
	}

	// Verify the update
	retrieved, err := svc.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if len(retrieved.Modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(retrieved.Modules))
	}

	// Verify BasicInfo was extracted
	if retrieved.BasicInfo.Name != "李四" {
		t.Errorf("Expected name '李四', got '%s'", retrieved.BasicInfo.Name)
	}

	// Verify markdown was regenerated
	if retrieved.MarkdownContent == "" {
		t.Error("MarkdownContent should be regenerated after UpdateJSON")
	}
	if retrieved.MarkdownContent != "" && !containsString(retrieved.MarkdownContent, "MIT") {
		t.Error("MarkdownContent should contain MIT from updated JSON")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
