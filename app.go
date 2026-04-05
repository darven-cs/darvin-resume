package main

import (
	"context"
	"log"

	"open-resume/internal/ai"
	"open-resume/internal/database"
	"open-resume/internal/model"
	"open-resume/internal/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	svc    service.ResumeService
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
}

// shutdown is called when the app stops
func (a *App) shutdown(ctx context.Context) {
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

// AI Configuration Bridge Methods

// GetAIConfig returns the current AI configuration as a map.
func (a *App) GetAIConfig() (map[string]interface{}, error) {
	cfg, err := ai.LoadConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"apiKey":        cfg.APIKey,
		"baseURL":       cfg.BaseURL,
		"defaultModel":  cfg.DefaultModel,
		"maxTokens":     cfg.MaxTokens,
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
		TimeoutSecs: getInt(config, "timeoutSeconds", ai.DefaultTimeoutSecs),
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
