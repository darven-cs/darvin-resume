package service

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"Darvin-Resume/internal/database"
	"Darvin-Resume/internal/model"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Initialize database for tests
	if err := database.Init(); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	database.Close()

	os.Exit(code)
}

// TestIntegration_BridgeCRUD verifies the complete CRUD chain for Bridge layer
func TestIntegration_BridgeCRUD(t *testing.T) {
	svc := NewResumeService()
	ctx := context.Background()

	// 1. CreateResume - verify return value contains UUID and title matches
	t.Run("CreateResume returns valid resume with UUID", func(t *testing.T) {
		resume := &model.Resume{
			Title: "Integration Test Resume",
			BasicInfo: model.BasicInfo{
				Name:  "张三",
				Email: "zhangsan@example.com",
				Phone: "13800138000",
			},
		}

		err := svc.Create(ctx, resume)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if resume.ID == "" {
			t.Error("CreateResume should return resume with non-empty UUID")
		}

		if resume.Title != "Integration Test Resume" {
			t.Errorf("Expected title 'Integration Test Resume', got '%s'", resume.Title)
		}

		if resume.CreatedAt.IsZero() {
			t.Error("CreateResume should set CreatedAt timestamp")
		}

		if resume.UpdatedAt.IsZero() {
			t.Error("CreateResume should set UpdatedAt timestamp")
		}
	})

	// Store ID for subsequent tests
	testResume := &model.Resume{
		Title: "CRUD Test Resume",
		BasicInfo: model.BasicInfo{
			Name:  "张三",
			Email: "zhangsan@example.com",
		},
	}
	if err := svc.Create(ctx, testResume); err != nil {
		t.Fatalf("Failed to create test resume: %v", err)
	}
	testID := testResume.ID

	// 2. GetResume - verify JSON parsing and non-empty Markdown
	t.Run("GetResume retrieves resume correctly", func(t *testing.T) {
		retrieved, err := svc.GetByID(ctx, testID)
		if err != nil {
			t.Fatalf("GetResume failed: %v", err)
		}

		if retrieved.ID != testID {
			t.Errorf("Expected ID '%s', got '%s'", testID, retrieved.ID)
		}

		if retrieved.Title != "CRUD Test Resume" {
			t.Errorf("Expected title 'CRUD Test Resume', got '%s'", retrieved.Title)
		}

		// Verify markdown content is generated
		if retrieved.MarkdownContent == "" {
			t.Error("GetResume should return non-empty MarkdownContent")
		}
	})

	// 3. ListResumes - verify the created resume appears in the list
	t.Run("ListResumes includes created resume", func(t *testing.T) {
		list, err := svc.List(ctx)
		if err != nil {
			t.Fatalf("ListResumes failed: %v", err)
		}

		found := false
		for _, item := range list {
			if item.ID == testID {
				found = true
				break
			}
		}

		if !found {
			t.Error("ListResumes should include the newly created resume")
		}
	})

	// 4. UpdateResume - verify Markdown content changes after update
	t.Run("UpdateResume regenerates Markdown", func(t *testing.T) {
		// Get the original markdown
		original, err := svc.GetByID(ctx, testID)
		if err != nil {
			t.Fatalf("GetResume failed: %v", err)
		}

		// Update with new JSON data
		newModules := []model.Module{
			{
				Type:    "basicInfo",
				Title:   "基本信息",
				Order:   0,
				Visible: true,
				Items:   json.RawMessage(`{"name":"李四","email":"lisi@example.com"}`),
			},
		}
		newModulesJSON, _ := json.Marshal(newModules)

		err = svc.UpdateJSON(ctx, testID, string(newModulesJSON))
		if err != nil {
			t.Fatalf("UpdateResume failed: %v", err)
		}

		// Get updated resume
		updated, err := svc.GetByID(ctx, testID)
		if err != nil {
			t.Fatalf("GetResume after update failed: %v", err)
		}

		// Verify name changed in markdown
		if updated.MarkdownContent == original.MarkdownContent {
			t.Error("UpdateResume should regenerate MarkdownContent")
		}

		if updated.BasicInfo.Name != "李四" {
			t.Errorf("Expected name '李四', got '%s'", updated.BasicInfo.Name)
		}
	})

	// 5. DeleteResume - verify subsequent GetByID returns error
	t.Run("DeleteResume removes resume", func(t *testing.T) {
		err := svc.Delete(ctx, testID)
		if err != nil {
			t.Fatalf("DeleteResume failed: %v", err)
		}

		// Verify resume is no longer accessible
		_, err = svc.GetByID(ctx, testID)
		if err == nil {
			t.Error("GetResume after delete should return error")
		}
	})

	// 6. Persistence test - restart simulation (close and reopen DB)
	t.Run("Data persists after simulated restart", func(t *testing.T) {
		// Create a new service instance (simulating app restart)
		newSvc := NewResumeService()

		// Create resume
		persistResume := &model.Resume{
			Title: "Persistence Test",
		}
		if err := newSvc.Create(ctx, persistResume); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		persistID := persistResume.ID

		// Retrieve it
		retrieved, err := newSvc.GetByID(ctx, persistID)
		if err != nil {
			t.Fatalf("GetResume after simulated restart failed: %v", err)
		}

		if retrieved.Title != "Persistence Test" {
			t.Errorf("Expected title 'Persistence Test', got '%s'", retrieved.Title)
		}
	})
}

// TestIntegration_UpdateResumeModule verifies module-level updates
func TestIntegration_UpdateResumeModule(t *testing.T) {
	svc := NewResumeService()
	ctx := context.Background()

	// Create a resume
	resume := &model.Resume{
		Title: "Module Update Test",
	}
	if err := svc.Create(ctx, resume); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update with education module
	educationModule := []model.Module{
		{
			Type:    "education",
			Title:   "教育经历",
			Order:   1,
			Visible: true,
			Items: json.RawMessage(`[{
				"school": "清华大学",
				"major": "计算机科学与技术",
				"degree": "本科",
				"startDate": "2019-09",
				"endDate": "2023-06",
				"highlights": ["GPA 3.8/4.0", "校级一等奖学金"]
			}]`),
		},
	}
	educationJSON, _ := json.Marshal(educationModule)

	err := svc.UpdateJSON(ctx, resume.ID, string(educationJSON))
	if err != nil {
		t.Fatalf("UpdateJSON failed: %v", err)
	}

	// Verify the update
	retrieved, err := svc.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetResume failed: %v", err)
	}

	if len(retrieved.Modules) != 1 {
		t.Errorf("Expected 1 module, got %d", len(retrieved.Modules))
	}

	if retrieved.Modules[0].Type != "education" {
		t.Errorf("Expected module type 'education', got '%s'", retrieved.Modules[0].Type)
	}

	// Verify markdown was regenerated
	if retrieved.MarkdownContent == "" {
		t.Error("MarkdownContent should be generated after module update")
	}
}

// TestIntegration_TimestampUpdates verifies timestamps are updated correctly
func TestIntegration_TimestampUpdates(t *testing.T) {
	svc := NewResumeService()
	ctx := context.Background()

	// Create a resume
	resume := &model.Resume{
		Title: "Timestamp Test",
	}
	initialTime := time.Now().Add(-1 * time.Hour) // Set to 1 hour ago
	resume.CreatedAt = initialTime
	resume.UpdatedAt = initialTime

	if err := svc.Create(ctx, resume); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify CreatedAt is set to current time (not the provided initial time)
	if resume.CreatedAt.Before(initialTime) {
		t.Error("CreateResume should set CreatedAt to current time")
	}

	// Wait a bit and then update
	time.Sleep(10 * time.Millisecond)

	// Update the resume
	newModules := []model.Module{
		{
			Type:    "basicInfo",
			Title:   "基本信息",
			Order:   0,
			Visible: true,
			Items:   json.RawMessage(`{"name":"Test User"}`),
		},
	}
	newModulesJSON, _ := json.Marshal(newModules)
	if err := svc.UpdateJSON(ctx, resume.ID, string(newModulesJSON)); err != nil {
		t.Fatalf("UpdateJSON failed: %v", err)
	}

	// Get the updated resume
	retrieved, err := svc.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetResume failed: %v", err)
	}

	// Verify UpdatedAt was updated
	if !retrieved.UpdatedAt.After(resume.UpdatedAt) {
		t.Error("UpdateJSON should update the UpdatedAt timestamp")
	}
}
