package ai

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"Darvin-Resume/internal/database"
)

// ErrChatNotFound is returned when no chat history exists for a resume.
var ErrChatNotFound = errors.New("chat not found")

// GetChatHistory retrieves all chat messages for a given resume ID.
// Messages are ordered by creation time ascending.
func GetChatHistory(ctx context.Context, resumeId string) ([]ChatMessage, error) {
	query := `
		SELECT id, resume_id, role, content, COALESCE(quoted_text, ''), created_at
		FROM ai_messages
		WHERE resume_id = ?
		ORDER BY created_at ASC
	`
	rows, err := database.DB.QueryContext(ctx, query, resumeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		var createdAt time.Time
		if err := rows.Scan(&msg.ID, &msg.ResumeID, &msg.Role, &msg.Content, &msg.QuotedText, &createdAt); err != nil {
			return nil, err
		}
		msg.CreatedAt = createdAt.UnixMilli()
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// SaveChatMessage persists a chat message to the database.
func SaveChatMessage(ctx context.Context, msg ChatMessage) error {
	if msg.ID == "" {
		return errors.New("message ID is required")
	}
	if msg.ResumeID == "" {
		return errors.New("resume ID is required")
	}
	if msg.Role != "user" && msg.Role != "assistant" {
		return errors.New("role must be 'user' or 'assistant'")
	}

	query := `
		INSERT INTO ai_messages (id, resume_id, role, content, quoted_text, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			role = excluded.role,
			content = excluded.content,
			quoted_text = excluded.quoted_text
	`
	createdAt := time.Now()
	if msg.CreatedAt > 0 {
		createdAt = time.UnixMilli(msg.CreatedAt)
	}

	_, err := database.DB.ExecContext(ctx, query,
		msg.ID, msg.ResumeID, msg.Role, msg.Content, msg.QuotedText, createdAt)
	return err
}

// ClearChatHistory removes all chat messages for a given resume ID.
func ClearChatHistory(ctx context.Context, resumeId string) error {
	query := `DELETE FROM ai_messages WHERE resume_id = ?`
	_, err := database.DB.ExecContext(ctx, query, resumeId)
	return err
}

// GetChatHistoryOrEmpty safely returns chat history or an empty slice (never nil).
func GetChatHistoryOrEmpty(ctx context.Context, resumeId string) ([]ChatMessage, error) {
	messages, err := GetChatHistory(ctx, resumeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []ChatMessage{}, nil
		}
		return nil, err
	}
	if messages == nil {
		return []ChatMessage{}, nil
	}
	return messages, nil
}
