package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"Darvin-Resume/internal/database"
	"Darvin-Resume/internal/model"

	"github.com/google/uuid"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var ErrSnapshotNotFound = errors.New("snapshot not found")
var ErrInvalidSnapshot = errors.New("invalid snapshot")

const MaxSnapshotsPerResume = 50

type SnapshotService interface {
	CreateSnapshot(ctx context.Context, req *model.CreateSnapshotRequest) (*model.Snapshot, error)
	ListSnapshots(ctx context.Context, resumeId string) ([]*model.SnapshotListItem, error)
	GetSnapshot(ctx context.Context, snapshotId string) (*model.Snapshot, error)
	DiffSnapshots(ctx context.Context, id1 string, id2 string) (*model.DiffResult, error)
	RollbackToSnapshot(ctx context.Context, resumeId string, snapshotId string) (*model.Snapshot, error)
	DeleteSnapshot(ctx context.Context, snapshotId string) error
	CleanupOldSnapshots(ctx context.Context, resumeId string, maxCount int) error
}

type snapshotService struct{}

var dmp = diffmatchpatch.New()

func NewSnapshotService() SnapshotService {
	return &snapshotService{}
}

// CreateSnapshot 创建版本快照
func (s *snapshotService) CreateSnapshot(ctx context.Context, req *model.CreateSnapshotRequest) (*model.Snapshot, error) {
	if req.ResumeID == "" {
		return nil, errors.New("resume id is required")
	}

	// 获取当前简历的完整数据
	resume, err := s.getResume(ctx, req.ResumeID)
	if err != nil {
		return nil, err
	}

	snapshot := &model.Snapshot{
		ID:              uuid.New().String(),
		ResumeID:        req.ResumeID,
		Label:           req.Label,
		Note:            req.Note,
		TriggerType:     req.TriggerType,
		JSONData:        resume.JSONData,
		MarkdownContent: resume.MarkdownContent,
		TemplateID:      resume.TemplateID,
		CustomCSS:       resume.CustomCSS,
		CreatedAt:       time.Now(),
	}

	query := `
		INSERT INTO snapshots (id, resume_id, label, note, trigger_type, json_data, markdown_content, template_id, custom_css, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = database.DB.ExecContext(ctx, query,
		snapshot.ID,
		snapshot.ResumeID,
		snapshot.Label,
		snapshot.Note,
		snapshot.TriggerType,
		snapshot.JSONData,
		snapshot.MarkdownContent,
		snapshot.TemplateID,
		snapshot.CustomCSS,
		snapshot.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 清理超出上限的旧快照
	_ = s.CleanupOldSnapshots(ctx, req.ResumeID, MaxSnapshotsPerResume)

	return snapshot, nil
}

// ListSnapshots 获取简历的所有快照（不含完整数据，仅列表信息）
func (s *snapshotService) ListSnapshots(ctx context.Context, resumeId string) ([]*model.SnapshotListItem, error) {
	query := `
		SELECT id, resume_id, label, note, trigger_type, created_at
		FROM snapshots
		WHERE resume_id = ?
		ORDER BY created_at DESC
	`
	rows, err := database.DB.QueryContext(ctx, query, resumeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.SnapshotListItem
	for rows.Next() {
		item := &model.SnapshotListItem{}
		if err := rows.Scan(&item.ID, &item.ResumeID, &item.Label, &item.Note, &item.TriggerType, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// GetSnapshot 获取单个快照的完整数据
func (s *snapshotService) GetSnapshot(ctx context.Context, snapshotId string) (*model.Snapshot, error) {
	query := `
		SELECT id, resume_id, label, note, trigger_type, json_data, markdown_content, template_id, custom_css, created_at
		FROM snapshots
		WHERE id = ?
	`
	row := database.DB.QueryRowContext(ctx, query, snapshotId)

	snapshot := &model.Snapshot{}
	err := row.Scan(
		&snapshot.ID,
		&snapshot.ResumeID,
		&snapshot.Label,
		&snapshot.Note,
		&snapshot.TriggerType,
		&snapshot.JSONData,
		&snapshot.MarkdownContent,
		&snapshot.TemplateID,
		&snapshot.CustomCSS,
		&snapshot.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSnapshotNotFound
		}
		return nil, err
	}
	return snapshot, nil
}

// DiffSnapshots 对比两个快照的内容差异
func (s *snapshotService) DiffSnapshots(ctx context.Context, id1 string, id2 string) (*model.DiffResult, error) {
	snap1, err := s.GetSnapshot(ctx, id1)
	if err != nil {
		return nil, err
	}
	snap2, err := s.GetSnapshot(ctx, id2)
	if err != nil {
		return nil, err
	}

	// 使用 diffmatchpatch 计算差异
	diffs := dmp.DiffMain(snap1.MarkdownContent, snap2.MarkdownContent, false)
	diffs = dmp.DiffCleanupSemantic(diffs)

	// 统计
	var added, removed int
	for _, d := range diffs {
		switch d.Type {
		case diffmatchpatch.DiffInsert:
			added += strings.Count(d.Text, "\n") + 1
		case diffmatchpatch.DiffDelete:
			removed += strings.Count(d.Text, "\n") + 1
		}
	}

	// 生成 unified diff 格式
	diffStr := dmp.DiffToDelta(diffs)

	return &model.DiffResult{
		Snapshot1ID:    id1,
		Snapshot2ID:    id2,
		Snapshot1Label: snap1.Label,
		Snapshot2Label: snap2.Label,
		ContentDiff:    diffStr,
		Stats: model.DiffStats{
			AddedLines:   added,
			RemovedLines: removed,
			ChangedLines: 0,
		},
	}, nil
}

// RollbackToSnapshot 回滚到指定快照
// 回滚前自动创建当前状态的快照（trigger_type = 'rollback'）
func (s *snapshotService) RollbackToSnapshot(ctx context.Context, resumeId string, snapshotId string) (*model.Snapshot, error) {
	// 1. 获取目标快照
	targetSnap, err := s.GetSnapshot(ctx, snapshotId)
	if err != nil {
		return nil, err
	}

	// 2. 回滚前自动创建当前快照
	autoSnap, err := s.CreateSnapshot(ctx, &model.CreateSnapshotRequest{
		ResumeID:    resumeId,
		Label:       fmt.Sprintf("回滚至「%s」前自动备份", targetSnap.Label),
		Note:        "系统自动创建的回滚前备份",
		TriggerType: "rollback",
	})
	if err != nil {
		// 回滚前备份失败不影响回滚操作，继续执行
		autoSnap = nil
	}

	// 3. 更新简历数据
	query := `
		UPDATE resumes
		SET json_data = ?, markdown_content = ?, template_id = ?, custom_css = ?, updated_at = ?
		WHERE id = ? AND is_deleted = FALSE
	`
	now := time.Now()
	result, err := database.DB.ExecContext(ctx, query,
		targetSnap.JSONData,
		targetSnap.MarkdownContent,
		targetSnap.TemplateID,
		targetSnap.CustomCSS,
		now,
		resumeId,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrResumeNotFound
	}

	return autoSnap, nil
}

// DeleteSnapshot 删除快照
func (s *snapshotService) DeleteSnapshot(ctx context.Context, snapshotId string) error {
	query := `DELETE FROM snapshots WHERE id = ?`
	result, err := database.DB.ExecContext(ctx, query, snapshotId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrSnapshotNotFound
	}
	return nil
}

// CleanupOldSnapshots 清理超出上限的旧快照
// 保留最新的 maxCount 个快照，删除更旧的
func (s *snapshotService) CleanupOldSnapshots(ctx context.Context, resumeId string, maxCount int) error {
	// 检查当前快照数量
	countQuery := `SELECT COUNT(*) FROM snapshots WHERE resume_id = ?`
	var count int
	if err := database.DB.QueryRowContext(ctx, countQuery, resumeId).Scan(&count); err != nil {
		return err
	}

	if count <= maxCount {
		return nil
	}

	// 删除超出上限的旧快照
	deleteCount := count - maxCount
	deleteQuery := `
		DELETE FROM snapshots
		WHERE id IN (
			SELECT id FROM snapshots
			WHERE resume_id = ?
			ORDER BY created_at ASC
			LIMIT ?
		)
	`
	_, err := database.DB.ExecContext(ctx, deleteQuery, resumeId, deleteCount)
	return err
}

// snapshotData 用于存储简历的快照相关数据
type snapshotData struct {
	JSONData        string
	MarkdownContent string
	TemplateID      string
	CustomCSS       string
}

// getResume 获取简历数据（内部辅助方法）
func (s *snapshotService) getResume(ctx context.Context, resumeId string) (*snapshotData, error) {
	query := `
		SELECT json_data, markdown_content, template_id, custom_css
		FROM resumes
		WHERE id = ? AND is_deleted = FALSE
	`
	row := database.DB.QueryRowContext(ctx, query, resumeId)
	data := &snapshotData{}
	err := row.Scan(&data.JSONData, &data.MarkdownContent, &data.TemplateID, &data.CustomCSS)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrResumeNotFound
		}
		return nil, err
	}
	return data, nil
}
