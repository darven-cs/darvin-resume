package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"Darvin-Resume/internal/settings"
)

// Message represents a chat message in the Messages API format.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIMessage represents an assistant message from the API.
type AIMessage struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence"`
	Usage        Usage          `json:"usage"`
}

// ContentBlock represents a content block in the response.
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Usage contains token usage information.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// streamingEvent represents a single SSE event from Claude's streaming response.
type streamingEvent struct {
	Type         string        `json:"type"`
	Index        int           `json:"index"`
	ContentBlock *ContentBlock `json:"content_block"`
	Delta        *TextDelta    `json:"delta"`
}

// TextDelta represents a text delta in a streaming response.
type TextDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// messagesRequest is the request body for the Claude Messages API.
type messagesRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
}

// messagesResponse is the non-streaming response from the Claude Messages API.
type messagesResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence"`
	Usage        Usage          `json:"usage"`
}

// client is the HTTP client for Claude API calls.
type client struct {
	httpClient *http.Client
	config     AIConfig
}

// NewClient creates a new AI client with the given config.
func NewClient(config AIConfig) *client {
	timeout := time.Duration(config.TimeoutSecs) * time.Second
	return &client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		config: config,
	}
}

// Chat sends a non-streaming chat request and returns the full response.
func (c *client) Chat(ctx context.Context, model string, messages []Message, maxTokens int, systemPrompt string) (string, error) {
	if model == "" {
		model = c.config.DefaultModel
	}
	if maxTokens <= 0 {
		maxTokens = c.config.MaxTokens
	}

	reqBody := messagesRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Stream:    false,
		Messages:  messages,
		System:    systemPrompt,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", &APIError{Code: AIErrBuildReq, Message: err.Error()}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", &APIError{Code: AIErrBuildReq, Message: err.Error()}
	}

	c.setHeaders(req, len(body))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", &APIError{Code: AIErrCancelled, Message: "request cancelled"}
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return "", &APIError{Code: AIErrTimeout, Message: "request timed out"}
		}
		return "", &APIError{Code: AIErrNetwork, Message: err.Error()}
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return "", err
	}

	var result messagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", &APIError{Code: AIErrParseResp, Message: err.Error()}
	}

	return extractTextContent(result.Content), nil
}

// ChatStream sends a streaming chat request and returns a ReadCloser for SSE data.
// The caller is responsible for closing the returned ReadCloser.
func (c *client) ChatStream(ctx context.Context, model string, messages []Message, maxTokens int, systemPrompt string) (io.ReadCloser, error) {
	if model == "" {
		model = c.config.DefaultModel
	}
	if maxTokens <= 0 {
		maxTokens = c.config.MaxTokens
	}

	reqBody := messagesRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Stream:    true,
		Messages:  messages,
		System:    systemPrompt,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &APIError{Code: AIErrBuildReq, Message: err.Error()}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, &APIError{Code: AIErrBuildReq, Message: err.Error()}
	}

	c.setHeaders(req, len(body))

	// Use a pipe so we can wrap the response body for streaming
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, &APIError{Code: AIErrCancelled, Message: "request cancelled"}
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &APIError{Code: AIErrTimeout, Message: "request timed out"}
		}
		return nil, &APIError{Code: AIErrNetwork, Message: err.Error()}
	}

	if err := c.checkResponse(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}

	return resp.Body, nil
}

// StreamHandler is a function type for handling streaming chunks.
type StreamHandler func(chunk string) error

// StreamEvents reads SSE events from the response body and calls the handler for each chunk.
// It handles both content chunks and error events.
func StreamEvents(body io.Reader, handler StreamHandler) error {
	// Use a scanner to read line by line (SSE format)
	// Claude SSE events come as: data: {"type": "...", ...}
	// We need to handle both chunk deltas and message deltas

	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 1024)
	for {
		n, err := body.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
			// Process complete lines
			for {
				lineEnd := bytes.Index(buf, []byte("\n"))
				if lineEnd < 0 {
					break
				}
				line := string(buf[:lineEnd])
				buf = buf[lineEnd+1:]

				line = strings.TrimSpace(line)
				if !strings.HasPrefix(line, "data: ") {
					continue
				}

				jsonStr := strings.TrimPrefix(line, "data: ")
				if jsonStr == "" || jsonStr == "[DONE]" {
					continue
				}

				var event streamingEvent
				if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
					// Try to parse as message delta for backwards compat
					var msgDelta struct {
						Type    string `json:"type"`
						Index   int    `json:"index"`
						Content struct {
							Type string `json:"type"`
							Text string `json:"text"`
						} `json:"delta"`
					}
					if unmarshalErr := json.Unmarshal([]byte(jsonStr), &msgDelta); unmarshalErr != nil {
						continue
					}
					if msgDelta.Type == "content_block_delta" && msgDelta.Content.Type == "text_delta" {
						if err := handler(msgDelta.Content.Text); err != nil {
							return err
						}
					}
					continue
				}

				switch event.Type {
				case "content_block_delta":
					if event.Delta != nil && event.Delta.Type == "text_delta" && event.Delta.Text != "" {
						if err := handler(event.Delta.Text); err != nil {
							return err
						}
					}
				case "message_stop":
					// End of stream
				case "error":
					// Error event from API
					var errData struct {
						Type  string `json:"type"`
						Error struct {
							Type    string `json:"type"`
							Message string `json:"message"`
						} `json:"error"`
					}
					if err := json.Unmarshal([]byte(jsonStr), &errData); err == nil {
						return &APIError{Code: AIErrAPI, Message: errData.Error.Message}
					}
				}
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

// ValidateAPIKey sends a lightweight request to verify the API key is valid.
func (c *client) ValidateAPIKey(ctx context.Context, apiKey string, baseURL string) (bool, error) {
	// Create a minimal client with the provided credentials
	testClient := &client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		config: AIConfig{
			APIKey:  apiKey,
			BaseURL: baseURL,
		},
	}

	// Send a minimal streaming request with empty content
	messages := []Message{{Role: "user", Content: "ping"}}
	body, err := json.Marshal(messagesRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 1,
		Stream:    false,
		Messages:  messages,
	})
	if err != nil {
		return false, &APIError{Code: AIErrBuildReq, Message: err.Error()}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return false, &APIError{Code: AIErrBuildReq, Message: err.Error()}
	}

	testClient.setHeaders(req, len(body))

	resp, err := testClient.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return false, nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return false, nil
		}
		return false, &APIError{Code: AIErrNetwork, Message: err.Error()}
	}
	defer resp.Body.Close()

	// 200 = valid, 401/403 = auth error, 429 = rate limit, anything else = maybe valid
	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return false, nil
	case http.StatusTooManyRequests:
		return false, &APIError{Code: AIErrRateLimit, Message: "rate limit exceeded"}
	default:
		// Other status codes might indicate other issues
		return false, nil
	}
}

// BuildMessages constructs the messages array for the Claude API.
// jobTarget is prepended to the system prompt.
// includeFullContext adds the full resume content as a user message.
func BuildMessages(systemPromptTemplate, userContent, jobTarget string, includeFullContext bool, fullResumeContent string) []Message {
	messages := make([]Message, 0, 3)

	// Build system prompt
	systemPrompt := systemPromptTemplate
	if jobTarget != "" {
		systemPrompt = fmt.Sprintf("%s\n\n目标岗位：%s", systemPrompt, jobTarget)
	}

	// Add system message
	if systemPrompt != "" {
		messages = append(messages, Message{
			Role:    "user",
			Content: systemPrompt,
		})
	}

	// Add user message
	messages = append(messages, Message{
		Role:    "user",
		Content: userContent,
	})

	return messages
}

// ChatMessage represents a stored chat message with role and content.
type ChatMessage struct {
	ID         string `json:"id"`
	ResumeID   string `json:"resumeId"`
	Role       string `json:"role"`
	Content    string `json:"content"`
	QuotedText string `json:"quotedText,omitempty"`
	CreatedAt  int64  `json:"createdAt"`
}

// BuildChatMessages constructs the messages array including conversation history.
// It builds a system prompt with context, prepends history messages (up to 10),
// and appends the current user prompt (with optional quoted text).
func BuildChatMessages(systemPromptTemplate, userContent, jobTarget string, historyMessages []ChatMessage, resumeContent string) []Message {
	messages := make([]Message, 0, 2+len(historyMessages))

	// Build context info for system prompt
	var contextInfo string
	if jobTarget != "" {
		contextInfo += fmt.Sprintf("\n目标岗位：%s", jobTarget)
	}
	if resumeContent != "" {
		contextInfo += fmt.Sprintf("\n当前简历内容：\n%s", resumeContent)
	}

	// Build system prompt
	systemPrompt := systemPromptTemplate
	if contextInfo != "" {
		systemPrompt += "\n\n当前简历编辑上下文：" + contextInfo
	}
	systemPrompt += "\n\n请以简历优化助手的身份回答。"

	messages = append(messages, Message{
		Role:    "user",
		Content: systemPrompt,
	})

	// Add history messages (up to 10 most recent)
	maxHistory := 10
	if len(historyMessages) > maxHistory {
		historyMessages = historyMessages[len(historyMessages)-maxHistory:]
	}
	for _, msg := range historyMessages {
		role := "user"
		if msg.Role == "assistant" {
			role = "assistant"
		}
		messages = append(messages, Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Add current user message with optional quoted text
	currentContent := userContent
	if resumeContent != "" {
		// Prepend resume content for context in current message
		currentContent = fmt.Sprintf("【简历内容】\n%s\n\n【你的问题】\n%s", resumeContent, userContent)
	}
	messages = append(messages, Message{
		Role:    "user",
		Content: currentContent,
	})

	return messages
}

// setHeaders sets the required headers for Claude API requests.
func (c *client) setHeaders(req *http.Request, bodyLen int) {
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("content-length", strconv.Itoa(bodyLen))
}

// checkResponse checks the HTTP response status and returns an appropriate error.
func (c *client) checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return &APIError{Code: AIErrAuth, Message: fmt.Sprintf("authentication failed: %s", resp.Status)}
	case http.StatusTooManyRequests:
		return &APIError{Code: AIErrRateLimit, Message: "rate limit exceeded"}
	default:
		// Try to read error message from body
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return &APIError{Code: AIErrAPI, Message: fmt.Sprintf("API error %d: %s", resp.StatusCode, string(body))}
	}
}

// extractTextContent extracts text from content blocks.
func extractTextContent(blocks []ContentBlock) string {
	var sb strings.Builder
	for _, block := range blocks {
		if block.Type == "text" && block.Text != "" {
			sb.WriteString(block.Text)
		}
	}
	return sb.String()
}

// LoadConfig loads AI configuration from settings storage.
// API Key is loaded from encrypted storage via LoadSecureAPIKey.
func LoadConfig(ctx context.Context) (AIConfig, error) {
	cfg := DefaultAIConfig()

	apiKey, err := LoadSecureAPIKey(ctx)
	if err != nil && !errors.Is(err, settings.ErrSettingNotFound) {
		return cfg, err
	}
	cfg.APIKey = apiKey

	baseURL, err := settings.Get(ctx, SettingKeyBaseURL)
	if err == nil {
		cfg.BaseURL = baseURL
	}

	defaultModel, err := settings.Get(ctx, SettingKeyDefaultModel)
	if err == nil {
		cfg.DefaultModel = defaultModel
	}

	maxTokensStr, err := settings.Get(ctx, SettingKeyMaxTokens)
	if err == nil {
		if val, parseErr := strconv.Atoi(maxTokensStr); parseErr == nil {
			cfg.MaxTokens = val
		}
	}

	timeoutStr, err := settings.Get(ctx, SettingKeyTimeoutSecs)
	if err == nil {
		if val, parseErr := strconv.Atoi(timeoutStr); parseErr == nil {
			cfg.TimeoutSecs = val
		}
	}

	return cfg, nil
}

// SaveConfig persists AI configuration to settings storage.
// API Key is encrypted with AES-256-GCM via SaveSecureAPIKey.
func SaveConfig(ctx context.Context, cfg AIConfig) error {
	if err := SaveSecureAPIKey(ctx, cfg.APIKey); err != nil {
		return err
	}
	if err := settings.Set(ctx, SettingKeyBaseURL, cfg.BaseURL); err != nil {
		return err
	}
	if err := settings.Set(ctx, SettingKeyDefaultModel, cfg.DefaultModel); err != nil {
		return err
	}
	if err := settings.Set(ctx, SettingKeyMaxTokens, strconv.Itoa(cfg.MaxTokens)); err != nil {
		return err
	}
	if err := settings.Set(ctx, SettingKeyTimeoutSecs, strconv.Itoa(cfg.TimeoutSecs)); err != nil {
		return err
	}
	return nil
}

// SystemPrompt is the default system prompt for resume optimization.
const SystemPrompt = `你是一位专业的简历优化助手。请根据用户需求优化简历内容。
要求：只返回优化后的文本内容，不要添加解释。`

// ChatWithRetry sends a chat request with automatic retry on format errors.
// It validates the response format and retries at most maxRetries times.
// If all attempts fail, returns the last response with a FORMAT_ERROR.
// This implements AIAI-12: format validation with auto-retry.
func (c *client) ChatWithRetry(ctx context.Context, model string, messages []Message, maxTokens int, systemPrompt string, maxRetries int) (string, error) {
	if maxRetries < 1 {
		maxRetries = 1
	}

	var lastResponse string

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := c.Chat(ctx, model, messages, maxTokens, systemPrompt)
		if err != nil {
			// Real error (not format) - return immediately
			return resp, err
		}

		lastResponse = resp

		// Validate response format
		if isValidResponse(resp) {
			return resp, nil
		}

		// Format invalid - retry unless we've exhausted retries
		if attempt < maxRetries {
			continue
		}
	}

	// All attempts exhausted - return raw response with format error
	// This implements the "fallback to raw output" requirement from AIAI-12
	return lastResponse, &APIError{
		Code:    AIErrParseResp,
		Message: "AI 返回格式异常，已保留原始输出，可手动编辑",
	}
}

// isValidResponse checks if the AI response looks like valid resume content.
func isValidResponse(content string) bool {
	if content == "" {
		return false
	}
	trimmed := strings.TrimSpace(content)
	if len(trimmed) < 10 {
		return false
	}
	// Check it's not a pure error message or junk
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "error:") ||
		strings.HasPrefix(lower, "sorry,") ||
		strings.HasPrefix(lower, "i can't") ||
		strings.HasPrefix(lower, "i cannot") {
		return false
	}
	return true
}

// activeOperations tracks which streaming operations are still active.
var activeOperations = &sync.Map{}

// RegisterOperation marks an operation as active for cancellation tracking.
func RegisterOperation(operationId string) {
	activeOperations.Store(operationId, true)
}

// UnregisterOperation removes an operation from the active set.
func UnregisterOperation(operationId string) {
	activeOperations.Delete(operationId)
}

// IsOperationActive checks if an operation is still registered as active.
func IsOperationActive(operationId string) bool {
	_, ok := activeOperations.Load(operationId)
	return ok
}

// CancelOperation removes an operation from the active set.
// The streaming loop should check IsOperationActive periodically.
// This implements AIAI-13: user-initiated abort with content preservation.
func CancelOperation(operationId string) bool {
	_, existed := activeOperations.LoadAndDelete(operationId)
	return existed
}
