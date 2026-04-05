package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"open-resume/internal/converter"
	"open-resume/internal/database"
	"open-resume/internal/model"

	"github.com/google/uuid"
)

// ResumeService defines the interface for resume operations
type ResumeService interface {
	Create(ctx context.Context, resume *model.Resume) error
	GetByID(ctx context.Context, id string) (*model.Resume, error)
	List(ctx context.Context) ([]*model.ResumeListItem, error)
	Update(ctx context.Context, resume *model.Resume) error
	Delete(ctx context.Context, id string) error
	UpdateJSON(ctx context.Context, id string, jsonData string) error
}

// ErrResumeNotFound is returned when a resume is not found
var ErrResumeNotFound = errors.New("resume not found")

// ErrInvalidData is returned when data is invalid
var ErrInvalidData = errors.New("invalid data")

// resumeService implements ResumeService interface
type resumeService struct{}

// NewResumeService creates a new resume service instance
func NewResumeService() ResumeService {
	return &resumeService{}
}

// Create creates a new resume with generated UUID and initializes empty JSON
func (s *resumeService) Create(ctx context.Context, resume *model.Resume) error {
	// Generate UUID if not set
	if resume.ID == "" {
		resume.ID = uuid.New().String()
	}

	// Initialize timestamps
	now := time.Now()
	resume.CreatedAt = now
	resume.UpdatedAt = now

	// Ensure modules is not nil
	if resume.Modules == nil {
		resume.Modules = []model.Module{}
	}

	// Add BasicInfo as a module if it has data
	if hasBasicInfoData(resume.BasicInfo) {
		basicInfoItems, err := json.Marshal(resume.BasicInfo)
		if err != nil {
			return err
		}
		// Insert basicInfo at the beginning (order 0)
		resume.Modules = append([]model.Module{
			{
				Type:    "basicInfo",
				Title:   "基本信息",
				Order:   0,
				Visible: true,
				Items:   basicInfoItems,
			},
		}, resume.Modules...)
	}

	// Serialize modules to JSON
	modulesJSON, err := json.Marshal(resume.Modules)
	if err != nil {
		return err
	}

	// Convert to markdown
	markdown, err := converter.ConvertResumeToMarkdown(resume)
	if err != nil {
		return err
	}

	// Insert into database
	query := `
		INSERT INTO resumes (id, title, json_data, markdown_content, template_id, custom_css, module_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = database.DB.ExecContext(ctx, query,
		resume.ID,
		resume.Title,
		string(modulesJSON),
		markdown,
		resume.TemplateID,
		resume.CustomCSS,
		"[]",
		resume.CreatedAt,
		resume.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a resume by its ID (only non-deleted)
func (s *resumeService) GetByID(ctx context.Context, id string) (*model.Resume, error) {
	query := `
		SELECT id, title, json_data, markdown_content, template_id, custom_css, module_order, created_at, updated_at
		FROM resumes
		WHERE id = ? AND is_deleted = FALSE
	`
	row := database.DB.QueryRowContext(ctx, query, id)

	var resume model.Resume
	var jsonData string
	var moduleOrder string
	err := row.Scan(
		&resume.ID,
		&resume.Title,
		&jsonData,
		&resume.MarkdownContent,
		&resume.TemplateID,
		&resume.CustomCSS,
		&moduleOrder,
		&resume.CreatedAt,
		&resume.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrResumeNotFound
		}
		return nil, err
	}

	// Parse JSON data
	if err := json.Unmarshal([]byte(jsonData), &resume.Modules); err != nil {
		return nil, err
	}

	// Also parse basic info from modules if present
	for _, module := range resume.Modules {
		if module.Type == "basicInfo" {
			if err := json.Unmarshal(module.Items, &resume.BasicInfo); err != nil {
				return nil, err
			}
			break
		}
	}

	return &resume, nil
}

// List returns all non-deleted resumes sorted by updated_at descending
func (s *resumeService) List(ctx context.Context) ([]*model.ResumeListItem, error) {
	query := `
		SELECT id, title, updated_at
		FROM resumes
		WHERE is_deleted = FALSE
		ORDER BY updated_at DESC
	`
	rows, err := database.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.ResumeListItem
	for rows.Next() {
		item := &model.ResumeListItem{}
		if err := rows.Scan(&item.ID, &item.Title, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// Update updates an existing resume and triggers JSON->Markdown sync
func (s *resumeService) Update(ctx context.Context, resume *model.Resume) error {
	// Update timestamp
	resume.UpdatedAt = time.Now()

	// Ensure modules is not nil
	if resume.Modules == nil {
		resume.Modules = []model.Module{}
	}

	// Add BasicInfo as a module if it has data (and not already present)
	if hasBasicInfoData(resume.BasicInfo) {
		// Check if basicInfo module already exists
		hasBasicInfoModule := false
		for _, m := range resume.Modules {
			if m.Type == "basicInfo" {
				hasBasicInfoModule = true
				break
			}
		}
		if !hasBasicInfoModule {
			basicInfoItems, err := json.Marshal(resume.BasicInfo)
			if err != nil {
				return err
			}
			// Insert basicInfo at the beginning
			resume.Modules = append([]model.Module{
				{
					Type:    "basicInfo",
					Title:   "基本信息",
					Order:   0,
					Visible: true,
					Items:   basicInfoItems,
				},
			}, resume.Modules...)
		}
	}

	// Serialize modules to JSON
	modulesJSON, err := json.Marshal(resume.Modules)
	if err != nil {
		return err
	}

	// Convert to markdown
	markdown, err := converter.ConvertResumeToMarkdown(resume)
	if err != nil {
		return err
	}

	query := `
		UPDATE resumes
		SET title = ?, json_data = ?, markdown_content = ?, template_id = ?, custom_css = ?, updated_at = ?
		WHERE id = ? AND is_deleted = FALSE
	`
	result, err := database.DB.ExecContext(ctx, query,
		resume.Title,
		string(modulesJSON),
		markdown,
		resume.TemplateID,
		resume.CustomCSS,
		resume.UpdatedAt,
		resume.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrResumeNotFound
	}

	return nil
}

// Delete soft-deletes a resume by setting is_deleted=true and deleted_at
func (s *resumeService) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE resumes
		SET is_deleted = TRUE, deleted_at = ?
		WHERE id = ? AND is_deleted = FALSE
	`
	now := time.Now()
	result, err := database.DB.ExecContext(ctx, query, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrResumeNotFound
	}

	return nil
}

// UpdateJSON updates only the JSON data and triggers Markdown regeneration
func (s *resumeService) UpdateJSON(ctx context.Context, id string, jsonData string) error {
	// Verify resume exists and is not deleted
	checkQuery := `SELECT id FROM resumes WHERE id = ? AND is_deleted = FALSE`
	var existingID string
	err := database.DB.QueryRowContext(ctx, checkQuery, id).Scan(&existingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrResumeNotFound
		}
		return err
	}

	// Parse the JSON to create a resume model for markdown conversion
	var modules []model.Module
	if err := json.Unmarshal([]byte(jsonData), &modules); err != nil {
		return err
	}

	// Extract basic info if present
	var basicInfo model.BasicInfo
	for _, module := range modules {
		if module.Type == "basicInfo" {
			if err := json.Unmarshal(module.Items, &basicInfo); err != nil {
				return err
			}
			break
		}
	}

	// Create temporary resume for markdown conversion
	tempResume := &model.Resume{
		ID:        id,
		BasicInfo: basicInfo,
		Modules:   modules,
	}

	// Convert to markdown
	markdown, err := converter.ConvertResumeToMarkdown(tempResume)
	if err != nil {
		return err
	}

	// Update only json_data and markdown_content
	updateQuery := `
		UPDATE resumes
		SET json_data = ?, markdown_content = ?, updated_at = ?
		WHERE id = ? AND is_deleted = FALSE
	`
	result, err := database.DB.ExecContext(ctx, updateQuery, jsonData, markdown, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrResumeNotFound
	}

	return nil
}

// hasBasicInfoData checks if BasicInfo has any non-empty fields
func hasBasicInfoData(info model.BasicInfo) bool {
	return info.Name != "" ||
		info.Phone != "" ||
		info.Email != "" ||
		info.Avatar != "" ||
		info.Website != "" ||
		info.GitHub != "" ||
		info.Address != "" ||
		info.Summary != ""
}
