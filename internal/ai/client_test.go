package ai

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBuildMessages(t *testing.T) {
	tests := []struct {
		name               string
		systemPrompt       string
		userContent        string
		jobTarget          string
		includeFullContext bool
		fullResume         string
		wantRoleCount      int
	}{
		{
			name:          "basic message",
			systemPrompt: "You are a helpful assistant.",
			userContent:  "Hello",
			jobTarget:    "",
			wantRoleCount: 2,
		},
		{
			name:          "with job target",
			systemPrompt: "You are a resume assistant.",
			userContent:  "Polish this",
			jobTarget:    "Frontend Engineer",
			wantRoleCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := BuildMessages(tt.systemPrompt, tt.userContent, tt.jobTarget, tt.includeFullContext, tt.fullResume)
			if len(messages) != tt.wantRoleCount {
				t.Errorf("BuildMessages() returned %d messages, want %d", len(messages), tt.wantRoleCount)
			}
			// Check that job target is included in system
			if tt.jobTarget != "" && !strings.Contains(messages[0].Content, tt.jobTarget) {
				t.Errorf("BuildMessages() system prompt should contain job target %q", tt.jobTarget)
			}
			// Check user content
			if messages[len(messages)-1].Content != tt.userContent {
				t.Errorf("BuildMessages() user content = %q, want %q", messages[len(messages)-1].Content, tt.userContent)
			}
		})
	}
}

func TestAIClient_Chat(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		// Create a test server that returns a valid Claude response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check headers
			if r.Header.Get("x-api-key") != "test-key" {
				t.Errorf("missing or wrong x-api-key header")
			}
			if r.Header.Get("anthropic-version") != "2023-06-01" {
				t.Errorf("missing or wrong anthropic-version header")
			}

			// Parse request body to verify structure
			var req messagesRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode request: %v", err)
			}
			if req.Model != "claude-sonnet-4-20250514" {
				t.Errorf("model = %q, want claude-sonnet-4-20250514", req.Model)
			}
			if req.Stream != false {
				t.Errorf("stream = %v, want false", req.Stream)
			}
			if len(req.Messages) == 0 {
				t.Errorf("Messages should not be empty")
			}

			// Return a valid response
			resp := messagesResponse{
				ID:   "msg_test123",
				Type: "message",
				Role: "assistant",
				Content: []ContentBlock{
					{Type: "text", Text: "Hello, this is a test response."},
				},
				Model:       "claude-sonnet-4-20250514",
				StopReason:  "end_turn",
				Usage:       Usage{InputTokens: 10, OutputTokens: 20},
			}
			w.Header().Set("content-type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient(AIConfig{
			APIKey:  "test-key",
			BaseURL: server.URL,
		})

		result, err := client.Chat(context.Background(), "claude-sonnet-4-20250514",
			[]Message{{Role: "user", Content: "Hello"}}, 100, "system prompt")
		if err != nil {
			t.Fatalf("Chat() error = %v", err)
		}
		if result != "Hello, this is a test response." {
			t.Errorf("Chat() = %q, want %q", result, "Hello, this is a test response.")
		}
	})

	t.Run("network error", func(t *testing.T) {
		client := NewClient(AIConfig{
			APIKey:  "test-key",
			BaseURL: "http://localhost:99999", // Invalid port
		})

		_, err := client.Chat(context.Background(), "claude-sonnet-4-20250514",
			[]Message{{Role: "user", Content: "Hello"}}, 100, "")
		if err == nil {
			t.Fatal("Chat() expected error for unreachable server")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.Code != AIErrNetwork {
			t.Errorf("error code = %q, want %q", apiErr.Code, AIErrNetwork)
		}
	})

	t.Run("auth error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"type": "error", "error": {"type": "authentication_error", "message": "invalid api key"}}`))
		}))
		defer server.Close()

		client := NewClient(AIConfig{
			APIKey:  "bad-key",
			BaseURL: server.URL,
		})

		_, err := client.Chat(context.Background(), "claude-sonnet-4-20250514",
			[]Message{{Role: "user", Content: "Hello"}}, 100, "")
		if err == nil {
			t.Fatal("Chat() expected error for auth failure")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.Code != AIErrAuth {
			t.Errorf("error code = %q, want %q", apiErr.Code, AIErrAuth)
		}
	})

	t.Run("rate limit error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer server.Close()

		client := NewClient(AIConfig{
			APIKey:  "test-key",
			BaseURL: server.URL,
		})

		_, err := client.Chat(context.Background(), "claude-sonnet-4-20250514",
			[]Message{{Role: "user", Content: "Hello"}}, 100, "")
		if err == nil {
			t.Fatal("Chat() expected error for rate limit")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.Code != AIErrRateLimit {
			t.Errorf("error code = %q, want %q", apiErr.Code, AIErrRateLimit)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond) // Longer than timeout
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		// Create client with very short timeout
		client := &client{
			httpClient: &http.Client{Timeout: 50 * time.Millisecond},
			config:     AIConfig{BaseURL: server.URL, TimeoutSecs: 1},
		}

		_, err := client.Chat(context.Background(), "claude-sonnet-4-20250514",
			[]Message{{Role: "user", Content: "Hello"}}, 100, "")
		if err == nil {
			t.Fatal("Chat() expected error for timeout")
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.Code != AIErrTimeout && apiErr.Code != AIErrNetwork {
			t.Errorf("error code = %q, expected %q or %q", apiErr.Code, AIErrTimeout, AIErrNetwork)
		}
	})
}

func TestValidateAPIKey(t *testing.T) {
	t.Run("valid key returns true", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := messagesResponse{
				ID:   "msg_valid",
				Type: "message",
				Role: "assistant",
				Content: []ContentBlock{
					{Type: "text", Text: "pong"},
				},
				StopReason: "end_turn",
				Usage:      Usage{InputTokens: 2, OutputTokens: 2},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient(AIConfig{BaseURL: server.URL})
		valid, err := client.ValidateAPIKey(context.Background(), "valid-key", server.URL)
		if err != nil {
			t.Fatalf("ValidateAPIKey() error = %v", err)
		}
		if !valid {
			t.Error("ValidateAPIKey() = false, want true for valid key")
		}
	})

	t.Run("invalid key returns false", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		client := NewClient(AIConfig{BaseURL: server.URL})
		valid, err := client.ValidateAPIKey(context.Background(), "bad-key", server.URL)
		if err != nil {
			t.Fatalf("ValidateAPIKey() error = %v", err)
		}
		if valid {
			t.Error("ValidateAPIKey() = true, want false for invalid key")
		}
	})
}

func TestStreamEvents(t *testing.T) {
	t.Run("parses SSE content chunks", func(t *testing.T) {
		// Simulate SSE stream with Claude-style events
		sseData := `data: {"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta", "text": "Hello"}}
data: {"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta", "text": " world"}}
data: {"type": "message_stop"}

`

		body := io.NopCloser(strings.NewReader(sseData))
		var chunks []string
		err := StreamEvents(body, func(chunk string) error {
			chunks = append(chunks, chunk)
			return nil
		})
		if err != nil {
			t.Fatalf("StreamEvents() error = %v", err)
		}
		if len(chunks) != 2 {
			t.Errorf("StreamEvents() collected %d chunks, want 2", len(chunks))
		}
		if chunks[0] != "Hello" {
			t.Errorf("chunk[0] = %q, want %q", chunks[0], "Hello")
		}
		if chunks[1] != " world" {
			t.Errorf("chunk[1] = %q, want %q", chunks[1], " world")
		}
	})

	t.Run("parses backwards-compatible message_delta format", func(t *testing.T) {
		sseData := `data: {"type": "content_block_delta", "index":0,"content_block":{"type":"text"},"delta":{"type":"text_delta","text":"Test"}}

`

		body := io.NopCloser(strings.NewReader(sseData))
		var chunks []string
		err := StreamEvents(body, func(chunk string) error {
			chunks = append(chunks, chunk)
			return nil
		})
		if err != nil {
			t.Fatalf("StreamEvents() error = %v", err)
		}
		if len(chunks) != 1 {
			t.Errorf("StreamEvents() collected %d chunks, want 1", len(chunks))
		}
	})

	t.Run("skips non-data lines", func(t *testing.T) {
		sseData := `event: message
data: {"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta", "text": "Skip test"}}

`

		body := io.NopCloser(strings.NewReader(sseData))
		var chunks []string
		err := StreamEvents(body, func(chunk string) error {
			chunks = append(chunks, chunk)
			return nil
		})
		if err != nil {
			t.Fatalf("StreamEvents() error = %v", err)
		}
		if len(chunks) != 1 {
			t.Errorf("StreamEvents() collected %d chunks, want 1 (event lines should be skipped)", len(chunks))
		}
	})
}

func TestAIConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  AIConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  AIConfig{BaseURL: "https://api.anthropic.com", DefaultModel: "claude-sonnet-4-20250514", MaxTokens: 4096, TimeoutSecs: 60},
			wantErr: false,
		},
		{
			name:    "empty baseurl gets default",
			config:  AIConfig{DefaultModel: "claude-sonnet-4-20250514", MaxTokens: 4096, TimeoutSecs: 60},
			wantErr: false,
		},
		{
			name:    "zero max tokens gets default",
			config:  AIConfig{BaseURL: "https://api.anthropic.com", DefaultModel: "claude-sonnet-4-20250514", TimeoutSecs: 60},
			wantErr: false,
		},
		{
			name:    "timeout too large",
			config:  AIConfig{BaseURL: "https://api.anthropic.com", DefaultModel: "claude-sonnet-4-20250514", MaxTokens: 4096, TimeoutSecs: 400},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AIConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAPIError(t *testing.T) {
	err := &APIError{Code: AIErrNetwork, Message: "connection refused"}
	if err.Error() != "AIError(network): connection refused" {
		t.Errorf("APIError.Error() = %q, want %q", err.Error(), "AIError(network): connection refused")
	}
}
