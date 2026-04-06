package ai

import "fmt"

// Default configuration values
const (
	DefaultBaseURL      = "https://api.anthropic.com"
	DefaultModel        = "claude-sonnet-4-20250514"
	DefaultMaxTokens    = 4096
	DefaultTimeoutSecs  = 60
)

// Settings keys (stored in SQLite settings table with ai. prefix)
const (
	SettingKeyAPIKey       = "ai.apiKey"
	SettingKeyBaseURL       = "ai.baseURL"
	SettingKeyDefaultModel  = "ai.defaultModel"
	SettingKeyMaxTokens     = "ai.maxTokens"
	SettingKeyTimeoutSecs   = "ai.timeoutSeconds"
)

// AIConfig holds the AI service configuration.
// APIKey is encrypted with AES-256-GCM via secure_config.go. Uses device-key
// derived from platform identifiers (Windows MachineGuid / macOS IOPlatformUUID /
// Linux machine-id). Same plaintext API key produces different ciphertext each
// encryption due to random salt + nonce (Argon2id key derivation).
type AIConfig struct {
	APIKey       string `json:"apiKey"`
	BaseURL      string `json:"baseURL"`
	DefaultModel string `json:"defaultModel"`
	MaxTokens    int    `json:"maxTokens"`
	TimeoutSecs  int    `json:"timeoutSeconds"`
}

// DefaultAIConfig returns a config with sensible defaults.
func DefaultAIConfig() AIConfig {
	return AIConfig{
		APIKey:       "",
		BaseURL:      DefaultBaseURL,
		DefaultModel: DefaultModel,
		MaxTokens:    DefaultMaxTokens,
		TimeoutSecs:  DefaultTimeoutSecs,
	}
}

// Validate checks if the config has required fields.
func (c *AIConfig) Validate() error {
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.DefaultModel == "" {
		c.DefaultModel = DefaultModel
	}
	if c.MaxTokens <= 0 {
		c.MaxTokens = DefaultMaxTokens
	}
	if c.TimeoutSecs <= 0 {
		c.TimeoutSecs = DefaultTimeoutSecs
	}
	if c.TimeoutSecs > 300 {
		return fmt.Errorf("timeout must be <= 300 seconds")
	}
	return nil
}

// APIError represents an error from the AI API.
type APIError struct {
	Code    string
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("AIError(%s): %s", e.Code, e.Message)
}

// Error codes for AI API errors.
const (
	AIErrNetwork    = "network"
	AIErrAuth      = "auth"
	AIErrRateLimit = "rate_limit"
	AIErrAPI       = "api"
	AIErrTimeout   = "timeout"
	AIErrCancelled = "cancelled"
	AIErrBuildReq  = "build_request"
	AIErrParseResp = "parse_response"
)
