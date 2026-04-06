package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"Darvin-Resume/internal/ai"
	"Darvin-Resume/internal/backup"
	"Darvin-Resume/internal/database"
	"Darvin-Resume/internal/export"
	"Darvin-Resume/internal/model"
	"Darvin-Resume/internal/service"
	"Darvin-Resume/internal/settings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	svc         service.ResumeService
	backupSched *backup.Scheduler
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		svc: service.NewResumeService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize auto backup scheduler based on saved settings
	a.initAutoBackup(ctx)
}

// initAutoBackup reads auto backup settings and starts the scheduler if enabled.
func (a *App) initAutoBackup(ctx context.Context) {
	enabled, _ := settings.Get(ctx, backup.SettingKeyAutoBackupEnabled)
	interval, _ := settings.Get(ctx, backup.SettingKeyAutoBackupInterval)

	if enabled == "true" {
		dur := backup.ParseInterval(interval)
		a.backupSched = backup.NewScheduler(dur)
		a.backupSched.Start(ctx)
		log.Printf("[App] auto backup scheduler started: interval=%v", dur)
	}
}

// shutdown is called when the app stops
func (a *App) shutdown(ctx context.Context) {
	// Stop auto backup scheduler
	if a.backupSched != nil {
		a.backupSched.Stop()
		a.backupSched = nil
	}

	// Close database connection
	if err := database.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return "Hello " + name + ", It's show time!"
}

// CreateResume creates a new resume with the given title
func (a *App) CreateResume(title string) (*model.Resume, error) {
	resume := &model.Resume{
		Title: title,
	}
	if err := a.svc.Create(context.Background(), resume); err != nil {
		return nil, err
	}
	return resume, nil
}

// GetResume retrieves a resume by its ID
func (a *App) GetResume(id string) (*model.Resume, error) {
	return a.svc.GetByID(context.Background(), id)
}

// ListResumes returns all non-deleted resumes
func (a *App) ListResumes() ([]*model.ResumeListItem, error) {
	return a.svc.List(context.Background())
}

// UpdateResume updates an existing resume with new JSON data
func (a *App) UpdateResume(id string, jsonData string) error {
	return a.svc.UpdateJSON(context.Background(), id, jsonData)
}

// DeleteResume soft-deletes a resume by its ID
func (a *App) DeleteResume(id string) error {
	return a.svc.Delete(context.Background(), id)
}

// UpdateResumeModule updates a specific module within a resume
func (a *App) UpdateResumeModule(id string, moduleType string, moduleData string) error {
	return a.svc.UpdateJSON(context.Background(), id, moduleData)
}

// RenameResume 修改简历标题
func (a *App) RenameResume(id string, title string) error {
	return a.svc.RenameResume(context.Background(), id, title)
}

// DuplicateResume 复制简历（新记录 + "(副本)" 后缀）
func (a *App) DuplicateResume(id string) (*model.Resume, error) {
	return a.svc.DuplicateResume(context.Background(), id)
}

// RestoreResume 恢复软删除的简历
func (a *App) RestoreResume(id string) error {
	return a.svc.RestoreResume(context.Background(), id)
}

// PermanentDeleteResume 物理删除已软删除的简历
func (a *App) PermanentDeleteResume(id string) error {
	return a.svc.PermanentDeleteResume(context.Background(), id)
}

// ListDeletedResumes 查询已软删除的简历列表
func (a *App) ListDeletedResumes() ([]*model.ResumeListItem, error) {
	return a.svc.ListDeletedResumes(context.Background())
}

// AI Configuration Bridge Methods

// GetAIConfig returns the current AI configuration as a map.
func (a *App) GetAIConfig() (map[string]interface{}, error) {
	cfg, err := ai.LoadConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"apiKey":         cfg.APIKey,
		"baseURL":        cfg.BaseURL,
		"defaultModel":   cfg.DefaultModel,
		"maxTokens":      cfg.MaxTokens,
		"timeoutSeconds": cfg.TimeoutSecs,
	}, nil
}

// SaveAIConfig saves AI configuration to persistent storage.
func (a *App) SaveAIConfig(config map[string]interface{}) error {
	cfg := ai.AIConfig{
		APIKey:       getString(config, "apiKey"),
		BaseURL:      getString(config, "baseURL", ai.DefaultBaseURL),
		DefaultModel: getString(config, "defaultModel", ai.DefaultModel),
		MaxTokens:    getInt(config, "maxTokens", ai.DefaultMaxTokens),
		TimeoutSecs:  getInt(config, "timeoutSeconds", ai.DefaultTimeoutSecs),
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	return ai.SaveConfig(context.Background(), cfg)
}

// ValidateAPIKey validates the given API key by sending a test request.
func (a *App) ValidateAPIKey(apiKey string, baseURL string) (bool, error) {
	testCfg := ai.AIConfig{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
	client := ai.NewClient(testCfg)
	return client.ValidateAPIKey(context.Background(), apiKey, baseURL)
}

// AISendMessage sends a streaming chat message to the AI service.
// operationId is used by the frontend to track and cancel the operation.
// Streamed chunks are emitted via Wails EventsEmit("ai:stream:{operationId}", chunk).
func (a *App) AISendMessage(operationId string, prompt string, jobTarget string, includeFullContext bool) (string, error) {
	// Load config
	cfg, err := ai.LoadConfig(context.Background())
	if err != nil {
		a.emitStreamError(operationId, err)
		return "", err
	}

	if cfg.APIKey == "" {
		err := &ai.APIError{Code: ai.AIErrAuth, Message: "API key not configured"}
		a.emitStreamError(operationId, err)
		return "", err
	}

	client := ai.NewClient(cfg)
	messages := ai.BuildMessages(ai.SystemPrompt, prompt, jobTarget, includeFullContext, "")

	// Start streaming
	body, err := client.ChatStream(context.Background(), cfg.DefaultModel, messages, cfg.MaxTokens, "")
	if err != nil {
		a.emitStreamError(operationId, err)
		return "", err
	}
	defer body.Close()

	// Stream chunks to frontend
	fullContent := ""
	err = ai.StreamEvents(body, func(chunk string) error {
		fullContent += chunk
		runtime.EventsEmit(a.ctx, "ai:stream:"+operationId, map[string]interface{}{
			"type":    "content",
			"content": chunk,
		})
		return nil
	})

	if err != nil {
		a.emitStreamError(operationId, err)
		return fullContent, err
	}

	// Emit done event
	runtime.EventsEmit(a.ctx, "ai:stream:"+operationId, map[string]interface{}{
		"type": "done",
	})

	return fullContent, nil
}

// AISendMessageSync sends a non-streaming chat message and returns the full response.
// Use this as a fallback when streaming is not supported.
func (a *App) AISendMessageSync(operationId string, prompt string, jobTarget string, includeFullContext bool) (string, error) {
	cfg, err := ai.LoadConfig(context.Background())
	if err != nil {
		return "", err
	}

	if cfg.APIKey == "" {
		return "", &ai.APIError{Code: ai.AIErrAuth, Message: "API key not configured"}
	}

	client := ai.NewClient(cfg)
	messages := ai.BuildMessages(ai.SystemPrompt, prompt, jobTarget, includeFullContext, "")

	content, err := client.Chat(context.Background(), cfg.DefaultModel, messages, cfg.MaxTokens, "")
	if err != nil {
		return "", err
	}

	// Emit done event
	runtime.EventsEmit(a.ctx, "ai:stream:"+operationId, map[string]interface{}{
		"type": "done",
	})

	return content, nil
}

// emitStreamError emits an error event to the frontend.
func (a *App) emitStreamError(operationId string, err error) {
	var msg string
	if apiErr, ok := err.(*ai.APIError); ok {
		msg = apiErr.Message
	} else {
		msg = err.Error()
	}
	runtime.EventsEmit(a.ctx, "ai:stream:"+operationId, map[string]interface{}{
		"type":  "error",
		"error": msg,
	})
}

// getString safely extracts a string from a map with optional default.
func getString(m map[string]interface{}, key string, defaults ...string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// getInt safely extracts an int from a map with optional default.
func getInt(m map[string]interface{}, key string, defaults ...int) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	if v, ok := m[key].(int); ok {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// ============================================================
// Backup Bridge Methods (EXPT-10, EXPT-11)
// ============================================================

// CreateManualBackup creates a manual backup file.
// password is optional; if empty, the backup is not encrypted.
func (a *App) CreateManualBackup(password string) (string, error) {
	path, err := backup.CreateBackup(password)
	if err != nil {
		return "", err
	}
	return path, nil
}

// RestoreFromBackup restores data from a backup file.
// password is required if the backup is encrypted.
func (a *App) RestoreFromBackup(filePath string, password string) error {
	return backup.RestoreBackup(filePath, password)
}

// ListBackups returns JSON array of all local backups.
func (a *App) ListBackups() (string, error) {
	backups, err := backup.ListBackups()
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(backups)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ShowSaveBackupDialog opens a save dialog for the user to choose a backup destination.
// Returns the selected path, or empty string if cancelled.
func (a *App) ShowSaveBackupDialog() (string, error) {
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "导出备份",
		DefaultFilename: fmt.Sprintf("darvin-resume-backup-%s.darvin-backup",
			time.Now().Format("20060102-150405")),
		Filters: []runtime.FileFilter{
			{DisplayName: "Darvin Resume 备份", Pattern: "*.darvin-backup"},
		},
	})
}

// ShowOpenBackupDialog opens a file dialog for the user to select a backup file.
// Returns the selected path, or empty string if cancelled.
func (a *App) ShowOpenBackupDialog() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择备份文件",
		Filters: []runtime.FileFilter{
			{DisplayName: "Darvin Resume 备份", Pattern: "*.darvin-backup"},
		},
	})
}

// ExportBackupToPath creates a backup and exports it to the specified path.
// password is optional.
func (a *App) ExportBackupToPath(savePath string, password string) (string, error) {
	tmpPath, err := backup.CreateBackup(password)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(savePath, data, 0644); err != nil {
		return "", err
	}
	return savePath, nil
}

// GetBackupDir returns the local backup directory path.
func (a *App) GetBackupDir() (string, error) {
	return backup.GetBackupDir()
}

// GetAutoBackupSettings returns the current auto backup settings as JSON.
func (a *App) GetAutoBackupSettings() (string, error) {
	enabled, _ := settings.Get(a.ctx, backup.SettingKeyAutoBackupEnabled)
	interval, _ := settings.Get(a.ctx, backup.SettingKeyAutoBackupInterval)

	if enabled == "" {
		enabled = "false"
	}
	if interval == "" {
		interval = "daily"
	}

	result := map[string]string{
		"enabled":  enabled,
		"interval": interval,
	}
	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetAutoBackupSettings enables or disables automatic backups.
func (a *App) SetAutoBackupSettings(enabled bool, interval string) error {
	enabledStr := "false"
	if enabled {
		enabledStr = "true"
	}

	if err := settings.Set(a.ctx, backup.SettingKeyAutoBackupEnabled, enabledStr); err != nil {
		return err
	}
	if err := settings.Set(a.ctx, backup.SettingKeyAutoBackupInterval, interval); err != nil {
		return err
	}

	// Update scheduler: stop existing, start new if enabled
	if a.backupSched != nil {
		a.backupSched.Stop()
		a.backupSched = nil
	}

	if enabled {
		dur := backup.ParseInterval(interval)
		a.backupSched = backup.NewScheduler(dur)
		a.backupSched.Start(a.ctx)
	}

	return nil
}

// ============================================================
// Snapshot Bridge Methods (EXPT-05 ~ EXPT-09)
// ============================================================

// CreateSnapshot 创建版本快照
func (a *App) CreateSnapshot(resumeId string, label string, note string, triggerType string) (*model.Snapshot, error) {
	svc := service.NewSnapshotService()
	req := &model.CreateSnapshotRequest{
		ResumeID:    resumeId,
		Label:       label,
		Note:        note,
		TriggerType: triggerType,
	}
	return svc.CreateSnapshot(context.Background(), req)
}

// ListSnapshots 获取简历的所有快照列表
func (a *App) ListSnapshots(resumeId string) ([]*model.SnapshotListItem, error) {
	svc := service.NewSnapshotService()
	return svc.ListSnapshots(context.Background(), resumeId)
}

// GetSnapshot 获取单个快照的完整数据
func (a *App) GetSnapshot(snapshotId string) (*model.Snapshot, error) {
	svc := service.NewSnapshotService()
	return svc.GetSnapshot(context.Background(), snapshotId)
}

// DiffSnapshots 对比两个快照
func (a *App) DiffSnapshots(id1 string, id2 string) (*model.DiffResult, error) {
	svc := service.NewSnapshotService()
	return svc.DiffSnapshots(context.Background(), id1, id2)
}

// RollbackToSnapshot 回滚到指定快照
func (a *App) RollbackToSnapshot(resumeId string, snapshotId string) (*model.Snapshot, error) {
	svc := service.NewSnapshotService()
	return svc.RollbackToSnapshot(context.Background(), resumeId, snapshotId)
}

// DeleteSnapshot 删除快照
func (a *App) DeleteSnapshot(snapshotId string) error {
	svc := service.NewSnapshotService()
	return svc.DeleteSnapshot(context.Background(), snapshotId)
}

// GetSnapshotMarkdown 获取快照的 Markdown 内容（用于编辑器加载）
func (a *App) GetSnapshotMarkdown(snapshotId string) (string, string, string, string, error) {
	svc := service.NewSnapshotService()
	snap, err := svc.GetSnapshot(context.Background(), snapshotId)
	if err != nil {
		return "", "", "", "", err
	}
	return snap.MarkdownContent, snap.TemplateID, snap.CustomCSS, snap.JSONData, nil
}

// ============================================================
// Chat History Bridge Methods
// ============================================================

// GetChatHistory retrieves all chat messages for a given resume.
func (a *App) GetChatHistory(resumeId string) ([]ai.ChatMessage, error) {
	return ai.GetChatHistoryOrEmpty(context.Background(), resumeId)
}

// SaveChatMessage persists a chat message to the database.
func (a *App) SaveChatMessage(msg ai.ChatMessage) error {
	return ai.SaveChatMessage(context.Background(), msg)
}

// ClearChatHistory removes all chat messages for a given resume.
func (a *App) ClearChatHistory(resumeId string) error {
	return ai.ClearChatHistory(context.Background(), resumeId)
}

// AISendChatMessage sends a chat message with conversation history context.
// It streams the response via Wails EventsOn("ai:stream:{operationId}").
// historyMessages contains up to 10 recent messages for context.
func (a *App) AISendChatMessage(operationId string, prompt string, jobTarget string, historyMessages []ai.ChatMessage, resumeContent string) (string, error) {
	// Load config
	cfg, err := ai.LoadConfig(context.Background())
	if err != nil {
		a.emitStreamError(operationId, err)
		return "", err
	}

	if cfg.APIKey == "" {
		err := &ai.APIError{Code: ai.AIErrAuth, Message: "API key not configured"}
		a.emitStreamError(operationId, err)
		return "", err
	}

	client := ai.NewClient(cfg)
	messages := ai.BuildChatMessages(ai.SystemPrompt, prompt, jobTarget, historyMessages, resumeContent)

	// Start streaming
	body, err := client.ChatStream(context.Background(), cfg.DefaultModel, messages, cfg.MaxTokens, "")
	if err != nil {
		a.emitStreamError(operationId, err)
		return "", err
	}
	defer body.Close()

	// Stream chunks to frontend
	fullContent := ""
	err = ai.StreamEvents(body, func(chunk string) error {
		fullContent += chunk
		runtime.EventsEmit(a.ctx, "ai:stream:"+operationId, map[string]interface{}{
			"type":    "content",
			"content": chunk,
		})
		return nil
	})

	if err != nil {
		a.emitStreamError(operationId, err)
		return fullContent, err
	}

	// Emit done event
	runtime.EventsEmit(a.ctx, "ai:stream:"+operationId, map[string]interface{}{
		"type": "done",
	})

	return fullContent, nil
}

// AICancelOperation cancels a running AI streaming operation.
// It removes the operation from the active set, allowing the streaming loop to exit.
// Implements AIAI-13: user-initiated abort with content preservation.
func (a *App) AICancelOperation(operationId string) error {
	ai.CancelOperation(operationId)
	return nil
}

// UpdateResumeTemplate 更新简历模板 ID
func (a *App) UpdateResumeTemplate(id string, templateId string) error {
	return a.svc.UpdateTemplateID(context.Background(), id, templateId)
}

// UpdateResumeCustomCSS 更新简历自定义 CSS
func (a *App) UpdateResumeCustomCSS(id string, customCss string) error {
	return a.svc.UpdateCustomCSS(context.Background(), id, customCss)
}

// ExportPDFFromHTML 使用 Chromedp 无头浏览器导出 PDF
// htmlContent: 完整的 HTML 内容（含 style 标签）
// outputPath: 输出文件路径
func (a *App) ExportPDFFromHTML(htmlContent string, outputPath string) (string, error) {
	ctx := context.Background()
	opts := &export.PDFOptions{
		PaperWidth:  8.27,
		PaperHeight: 11.69,
		Scale:       1.0, // DPI 96 / 72
		PrintBg:     true,
	}
	err := export.ExportPDFFromHTML(ctx, htmlContent, outputPath, opts)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

// ShowSaveDialog 显示系统文件保存对话框并返回用户选择的路径
func (a *App) ShowSaveDialog(dialogConfig map[string]interface{}) (map[string]interface{}, error) {
	title := getString(dialogConfig, "title", "保存文件")
	defaultFilename := getString(dialogConfig, "defaultPath", "")

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"filePath": path,
		"canceled": path == "",
	}, nil
}
