package ai

import (
	"strings"
)

// AIErrorCode represents a categorized error code for AI operations.
// These codes are sent to the frontend for user-facing error handling.
type AIErrorCode string

// Error code constants matching frontend AIErrorCode type.
const (
	CodeNetworkError  AIErrorCode = "NETWORK_ERROR"
	CodeTimeout       AIErrorCode = "TIMEOUT"
	CodeAuthFailed    AIErrorCode = "AUTH_FAILED"
	CodeRateLimit     AIErrorCode = "RATE_LIMIT"
	CodeTokenLimit    AIErrorCode = "TOKEN_LIMIT"
	CodeFormatError   AIErrorCode = "FORMAT_ERROR"
	CodeAborted       AIErrorCode = "ABORTED"
	CodeUnknown       AIErrorCode = "UNKNOWN"
)

// errorMapping maps backend error codes and error message substrings
// to structured AI error codes for the frontend.
var errorMapping = map[string]AIErrorCode{
	"context_length_exceeded": CodeTokenLimit,
	"rate_limit_exceeded":     CodeRateLimit,
	"authentication_error":    CodeAuthFailed,
	"invalid_request_error":   CodeFormatError,
}

// ClassifyError converts a backend error to a structured AI error code
// and user-friendly message based on error type, code, and HTTP status.
func ClassifyError(err error) (AIErrorCode, string) {
	if err == nil {
		return CodeUnknown, "未知错误"
	}

	errMsg := err.Error()

	// Check for APIError with specific code
	if apiErr, ok := err.(*APIError); ok {
		// Check substring mappings first
		if mapped, ok := errorMapping[apiErr.Code]; ok {
			return mapped, apiErr.Message
		}

		// Fallback code-to-type mapping
		switch apiErr.Code {
		case AIErrNetwork:
			return CodeNetworkError, apiErr.Message
		case AIErrAuth:
			return CodeAuthFailed, apiErr.Message
		case AIErrRateLimit:
			return CodeRateLimit, apiErr.Message
		case AIErrAPI:
			// Check API error message for specific issues
			if strings.Contains(errMsg, "context_length") || strings.Contains(errMsg, "token") {
				return CodeTokenLimit, "内容超出模型限制，请精简后重试"
			}
			if strings.Contains(errMsg, "rate_limit") || strings.Contains(errMsg, "429") {
				return CodeRateLimit, "请求频率超限，请稍后再试"
			}
			return CodeUnknown, apiErr.Message
		case AIErrTimeout:
			return CodeTimeout, apiErr.Message
		case AIErrCancelled:
			return CodeAborted, "操作已取消"
		case AIErrBuildReq:
			return CodeFormatError, "请求构建失败，请检查输入内容"
		case AIErrParseResp:
			return CodeFormatError, "响应解析失败，AI 返回格式异常"
		}
	}

	// Check error message substring patterns
	msgLower := strings.ToLower(errMsg)
	switch {
	case strings.Contains(msgLower, "connection") ||
		strings.Contains(msgLower, "network") ||
		strings.Contains(msgLower, "dial") ||
		strings.Contains(msgLower, "refused"):
		return CodeNetworkError, "网络连接失败，请检查网络后重试"
	case strings.Contains(msgLower, "timeout") ||
		strings.Contains(msgLower, "deadline"):
		return CodeTimeout, "请求超时，请增加超时时间或减少内容后重试"
	case strings.Contains(msgLower, "unauthorized") ||
		strings.Contains(msgLower, "auth") ||
		strings.Contains(msgLower, "401") ||
		strings.Contains(msgLower, "403"):
		return CodeAuthFailed, "API 密钥无效或已过期，请在设置中更新"
	case strings.Contains(msgLower, "rate limit") ||
		strings.Contains(msgLower, "429"):
		return CodeRateLimit, "请求过于频繁，请稍后再试"
	case strings.Contains(msgLower, "context_length") ||
		strings.Contains(msgLower, "token_limit") ||
		strings.Contains(msgLower, "maximum context"):
		return CodeTokenLimit, "简历内容超出了模型处理的 Token 限制，请精简后重试"
	case strings.Contains(msgLower, "cancel"):
		return CodeAborted, "操作已取消"
	}

	return CodeUnknown, errMsg
}

// ToAIError converts a backend error to a structured AIError for the frontend.
func ToAIError(err error) AIError {
	code, msg := ClassifyError(err)

	recoverable := code != CodeAuthFailed
	retryable := code == CodeNetworkError || code == CodeTimeout || code == CodeRateLimit

	return AIError{
		Code:         code,
		Message:      msg,
		Detail:       err.Error(),
		Recoverable:  recoverable,
		Retryable:    retryable,
	}
}

// AIError is a structured error with user-facing metadata.
// Matches the frontend AIError interface.
type AIError struct {
	Code        AIErrorCode
	Message     string
	Detail      string
	Recoverable bool
	Retryable   bool
}

// IsRetryable returns true if the error can be retried.
func (e AIError) IsRetryable() bool {
	return e.Retryable
}

// IsRecoverable returns true if the error allows data recovery.
func (e AIError) IsRecoverable() bool {
	return e.Recoverable
}
